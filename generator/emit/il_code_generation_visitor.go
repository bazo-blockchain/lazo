package emit

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/data"
	"github.com/bazo-blockchain/lazo/generator/il"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"math/big"
)

// ILCodeGenerationVisitor generates the IL Code
type ILCodeGenerationVisitor struct {
	node.AbstractVisitor
	symbolTable *symbol.SymbolTable
	ilBuilder   *ILBuilder
	function    *symbol.FunctionSymbol
	assembler   *ILAssembler
	bytePos     uint16
	Errors      []error
}

// NewCodeGenerationVisitor creates a new CodeGenerationVisitor
func NewCodeGenerationVisitor(
	symbolTable *symbol.SymbolTable, ilBuilder *ILBuilder) *ILCodeGenerationVisitor {
	v := &ILCodeGenerationVisitor{
		symbolTable: symbolTable,
		ilBuilder:   ilBuilder,
		bytePos:     0,
	}
	v.ConcreteVisitor = v
	return v
}

// VisitContractNode creates a new IL Assembler, generates the ABI, Constructor IL Code and Function IL Code
func (v *ILCodeGenerationVisitor) VisitContractNode(node *node.ContractNode) {
	contractSymbol := v.symbolTable.GlobalScope.Contract
	contractData := v.ilBuilder.Metadata.Contract

	v.assembler = NewILAssembler(&v.bytePos)
	v.generateABI(contractSymbol, contractData)
	v.generateConstructorIL(node, contractSymbol, contractData)
	v.generateFunctionIL(node, contractSymbol, contractData)
}

func (v *ILCodeGenerationVisitor) generateABI(contractSymbol *symbol.ContractSymbol,
	contractData *data.ContractData) {
	v.assembler.Emit(il.CallData)

	for i, functionData := range contractData.Functions {
		v.assembler.Emit(il.Dup)
		v.assembler.PushFuncHash(functionData.Hash)
		v.assembler.Emit(il.NotEq)

		checkNextFuncLabel := v.assembler.CreateLabel()
		v.assembler.JmpTrue(checkNextFuncLabel)
		v.assembler.Emit(il.Pop) // Remove function hash from top of call stack
		v.assembler.CallFunc(contractSymbol.Functions[i])
		v.assembler.Emit(il.Halt)

		v.assembler.SetLabel(checkNextFuncLabel)
	}
}

func (v *ILCodeGenerationVisitor) generateConstructorIL(node *node.ContractNode,
	contractSymbol *symbol.ContractSymbol, contractData *data.ContractData) {
	constructorLabel := v.assembler.CreateLabel()

	v.assembler.PushInt(big.NewInt(0))
	v.assembler.Emit(il.Eq)
	v.assembler.JmpTrue(constructorLabel)
	v.assembler.Emit(il.Halt)

	v.assembler.SetLabel(constructorLabel)
	for _, variable := range node.Fields {
		variable.Accept(v.ConcreteVisitor)
	}

	if node.Constructor != nil {
		v.function = contractSymbol.Constructor
		v.assembler.CallFunc(contractSymbol.Constructor)
		v.ilBuilder.SetFunctionPos(v.function, v.bytePos)
		node.Constructor.Accept(v)
		v.function = nil
	}
	contractData.Instructions = v.assembler.Complete(true)
}

// VisitFieldNode generates the IL Code for a contract field node and default initializes it if required
func (v *ILCodeGenerationVisitor) VisitFieldNode(node *node.FieldNode) {
	v.AbstractVisitor.VisitFieldNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)

	if arrayType, ok := targetType.(*symbol.ArrayTypeSymbol); ok {
		if _, ok := arrayType.ElementType.(*symbol.BasicTypeSymbol); !ok {
			v.reportError(node, "Generator currently does not support array nesting")
			return
		}
	}

	if node.Expression == nil {
		v.pushDefault(targetType)
	}

	index := v.symbolTable.GlobalScope.Contract.GetFieldIndex(node.Identifier)
	v.assembler.StoreState(byte(index))
}

func (v *ILCodeGenerationVisitor) generateFunctionIL(node *node.ContractNode, contractSymbol *symbol.ContractSymbol,
	contractData *data.ContractData) {
	for i, function := range node.Functions {
		v.function = contractSymbol.Functions[i]
		funcData := contractData.Functions[i]

		v.ilBuilder.SetFunctionPos(v.function, v.bytePos)
		v.assembler = NewILAssembler(&v.bytePos)
		function.Accept(v.ConcreteVisitor)

		funcData.Instructions = v.assembler.Complete(false)
		v.function = nil
	}
}

// Statements
// -----------

// VisitVariableNode generates the IL Code for a variable node and default initializes it if required
func (v *ILCodeGenerationVisitor) VisitVariableNode(node *node.VariableNode) {
	v.AbstractVisitor.VisitVariableNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)

	if node.Expression == nil {
		v.pushDefault(targetType)
	}

	index := v.function.GetVarIndex(node.Identifier)
	v.assembler.StoreLocal(byte(index))
}

// VisitMultiVariableNode generates the IL Code for multi-variable initialization
func (v *ILCodeGenerationVisitor) VisitMultiVariableNode(node *node.MultiVariableNode) {
	v.AbstractVisitor.VisitMultiVariableNode(node)

	for i := len(node.Identifiers) - 1; i >= 0; i-- {
		index := v.function.GetVarIndex(node.Identifiers[i])
		v.assembler.StoreLocal(byte(index))
	}
}

// VisitAssignmentStatementNode generates the IL Code for an assignment
func (v *ILCodeGenerationVisitor) VisitAssignmentStatementNode(assignNode *node.AssignmentStatementNode) {
	switch assignNode.Left.(type) {
	case *node.BasicDesignatorNode:
		assignNode.Right.Accept(v)
		decl := v.symbolTable.GetDeclByDesignator(assignNode.Left)
		v.storeVariable(decl)

	case *node.ElementAccessNode:
		elementAccess, _ := assignNode.Left.(*node.ElementAccessNode)
		assignNode.Right.Accept(v)
		elementAccess.Expression.Accept(v)
		elementAccess.Designator.Accept(v)
		v.assembler.Emit(il.ArrInsert)
		fieldSymbol := v.symbolTable.GetDeclByDesignator(elementAccess.Designator)
		v.storeVariable(fieldSymbol)

	case *node.MemberAccessNode: // this.field or struct.field
		memberAccessNode := assignNode.Left.(*node.MemberAccessNode)
		memberAccessNode.Designator.Accept(v)
		assignNode.Right.Accept(v)
		fieldSymbol := v.symbolTable.GetDeclByDesignator(memberAccessNode)
		v.storeVariable(fieldSymbol)

		// Struct is a value type, therefore update struct explicitly.
		if _, ok := fieldSymbol.Scope().(*symbol.StructTypeSymbol); ok {
			v.updateStruct(memberAccessNode.Designator)
		}
	default:
		v.reportError(assignNode, fmt.Sprintf("Invalid assignment %v", assignNode.Left))
	}
}

// VisitMultiAssignmentStatementNode generates the IL Code for a multi-assignment
func (v *ILCodeGenerationVisitor) VisitMultiAssignmentStatementNode(assignNode *node.MultiAssignmentStatementNode) {
	assignNode.FuncCall.Accept(v)

	for i := len(assignNode.Designators) - 1; i >= 0; i-- {
		switch assignNode.Designators[i].(type) {
		case *node.BasicDesignatorNode:
			decl := v.symbolTable.GetDeclByDesignator(assignNode.Designators[i])
			v.storeVariable(decl)

		case *node.ElementAccessNode:
			elementAccess, _ := assignNode.Designators[i].(*node.ElementAccessNode)
			elementAccess.Expression.Accept(v)
			elementAccess.Designator.Accept(v)
			v.assembler.Emit(il.ArrInsert)
			fieldSymbol := v.symbolTable.GetDeclByDesignator(elementAccess.Designator)
			v.storeVariable(fieldSymbol)

		case *node.MemberAccessNode:
			memberAccessNode := assignNode.Designators[i].(*node.MemberAccessNode)
			memberAccessNode.Designator.Accept(v)

			fieldSymbol := v.symbolTable.GetDeclByDesignator(memberAccessNode)
			_, isStructField := fieldSymbol.Scope().(*symbol.StructTypeSymbol)
			if isStructField {
				v.assembler.Emit(il.Swap) // Swap field value and struct to match the StoreFld opcode
			}
			v.storeVariable(fieldSymbol)

			// Struct is a value type in VM. Therefore, struct variable should be updated explicitly.
			if isStructField {
				v.updateStruct(memberAccessNode.Designator)
			}
		default:
			v.reportError(assignNode, fmt.Sprintf("Invalid assignment %v", assignNode.Designators[i]))
		}
	}
}

// Struct is a value type in VM.
// Therefore, every struct field assignment should update its parent explicitly (e.g. struct.field = x)
func (v *ILCodeGenerationVisitor) updateStruct(targetStruct node.DesignatorNode) {
	// e.g. this.targetStruct.field or grandParent.targetStruct.field
	if targetStructMemberAccessNode, ok := targetStruct.(*node.MemberAccessNode); ok {
		if targetStructMemberAccessNode.Designator.String() != symbol.This {
			targetStructMemberAccessNode.Designator.Accept(v) // load grand parent struct on stack
			v.assembler.Emit(il.Swap)                         // Swap grand parent struct and element value
		}

		targetStructFieldSymbol := v.symbolTable.GetDeclByDesignator(targetStructMemberAccessNode)
		v.storeVariable(targetStructFieldSymbol)

		// Struct is a value type, therefore update struct explicitly.
		if _, ok := targetStructFieldSymbol.Scope().(*symbol.StructTypeSymbol); ok {
			v.updateStruct(targetStructMemberAccessNode.Designator)
		}
		return
	}

	// targetStruct.field = x
	v.storeVariable(v.symbolTable.GetDeclByDesignator(targetStruct))
}

// VisitMemberAccessNode generates the IL Code for a member access node
func (v *ILCodeGenerationVisitor) VisitMemberAccessNode(node *node.MemberAccessNode) {
	if node.Designator.String() == symbol.This {
		index := v.symbolTable.GlobalScope.Contract.GetFieldIndex(node.Identifier)
		v.assembler.LoadState(byte(index))
		return
	}

	node.Designator.Accept(v)

	// TODO as soon as VM has ARRLEN Opcode
	//if node.Identifier == "length" && v.isArray(node.Designator) {
	//	v.assembler.Emit(il.ArrLen)
	//}

	decl := v.symbolTable.GetDeclByDesignator(node)
	v.loadVariable(decl)
}

// VisitElementAccessNode generates the il code for an array element access
func (v *ILCodeGenerationVisitor) VisitElementAccessNode(node *node.ElementAccessNode) {
	node.Expression.Accept(v)
	node.Designator.Accept(v)
	v.assembler.Emit(il.ArrAt) // Load Array Element
}

// VisitReturnStatementNode generates the IL Code for returning within a function
func (v *ILCodeGenerationVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(node)
	v.assembler.Emit(il.Ret)
}

// VisitIfStatementNode generates the IL Code for an If or an If-Else Statement
func (v *ILCodeGenerationVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	elseLabel := v.assembler.CreateLabel()
	endLabel := v.assembler.CreateLabel()

	// Condition
	node.Condition.Accept(v)
	v.assembler.JmpFalse(elseLabel)

	// Then
	v.VisitStatementBlock(node.Then)
	v.assembler.Jmp(endLabel)

	// Else
	v.assembler.SetLabel(elseLabel)
	v.VisitStatementBlock(node.Else)

	v.assembler.SetLabel(endLabel)
}

// Expressions
// -----------

var binaryOpCodes = map[token.Symbol]il.OpCode{
	token.Plus:           il.Add,
	token.Minus:          il.Sub,
	token.Multiplication: il.Mul,
	token.Division:       il.Div,
	token.Modulo:         il.Mod,
	token.Greater:        il.Gt,
	token.GreaterEqual:   il.GtEq,
	token.LessEqual:      il.LtEq,
	token.Less:           il.Lt,
	token.Equal:          il.Eq,
	token.Unequal:        il.NotEq,
}

// VisitBinaryExpressionNode generates the IL Code for all Binary Expressions
func (v *ILCodeGenerationVisitor) VisitBinaryExpressionNode(expNode *node.BinaryExpressionNode) {
	if op, ok := binaryOpCodes[expNode.Operator]; ok {
		v.AbstractVisitor.VisitBinaryExpressionNode(expNode)
		v.assembler.Emit(op)
		return
	}

	if expNode.Operator == token.Exponent {
		// Visit right node first because of right associativity
		expNode.Right.Accept(v) // exponent
		expNode.Left.Accept(v)  // basis
		v.assembler.Emit(il.Exp)
		return
	}

	if expNode.Operator == token.And {
		falseLabel := v.assembler.CreateLabel()
		endLabel := v.assembler.CreateLabel()

		expNode.Left.Accept(v)
		v.assembler.JmpFalse(falseLabel)
		expNode.Right.Accept(v)
		v.assembler.Jmp(endLabel)

		v.assembler.SetLabel(falseLabel)
		v.assembler.PushBool(false)

		v.assembler.SetLabel(endLabel)
		return
	}

	if expNode.Operator == token.Or {
		trueLabel := v.assembler.CreateLabel()
		endLabel := v.assembler.CreateLabel()

		expNode.Left.Accept(v)
		v.assembler.JmpTrue(trueLabel)
		expNode.Right.Accept(v)
		v.assembler.Jmp(endLabel)

		v.assembler.SetLabel(trueLabel)
		v.assembler.PushBool(true)

		v.assembler.SetLabel(endLabel)
		return
	}
	v.reportError(expNode, fmt.Sprintf("binary operator %s not supported", token.SymbolLexeme[expNode.Operator]))
}

var unaryOpCodes = map[token.Symbol]il.OpCode{
	token.Minus: il.Neg,
}

// VisitUnaryExpressionNode generates the IL Code for all unary expressions
func (v *ILCodeGenerationVisitor) VisitUnaryExpressionNode(expNode *node.UnaryExpressionNode) {
	if op, ok := unaryOpCodes[expNode.Operator]; ok {
		v.AbstractVisitor.VisitUnaryExpressionNode(expNode)
		v.assembler.Emit(op)
		return
	}

	if expNode.Operator == token.Plus {
		v.AbstractVisitor.VisitUnaryExpressionNode(expNode)
		return
	}

	if expNode.Operator == token.Not {
		v.AbstractVisitor.VisitUnaryExpressionNode(expNode)
		v.assembler.Emit(il.Neg)
		return
	}
	v.reportError(expNode, fmt.Sprintf("unary operator %s not supported", token.SymbolLexeme[expNode.Operator]))
}

// VisitFuncCallNode generates the IL Code for the function call
func (v *ILCodeGenerationVisitor) VisitFuncCallNode(node *node.FuncCallNode) {
	for _, arg := range node.Args {
		arg.Accept(v.ConcreteVisitor)
	}

	funcSym := v.symbolTable.GetDeclByDesignator(node.Designator).(*symbol.FunctionSymbol)
	v.assembler.CallFunc(funcSym)
}

// VisitStructCreationNode generates the IL code for creating a new struct.
func (v *ILCodeGenerationVisitor) VisitStructCreationNode(node *node.StructCreationNode) {
	structType := v.symbolTable.GetTypeByExpression(node).(*symbol.StructTypeSymbol)
	v.assembler.NewStruct(uint16(len(structType.Fields)))

	for i, expr := range node.FieldValues {
		expr.Accept(v.ConcreteVisitor)
		v.assembler.StoreField(uint16(i))
	}

	// Set default value when field is not initialized
	for i := len(node.FieldValues); i < len(structType.Fields); i++ {
		v.pushDefault(structType.Fields[i].Type)
		v.assembler.StoreField(uint16(i))
	}
}

// VisitStructNamedCreationNode generates the IL code for creating a new struct with named field initialization.
func (v *ILCodeGenerationVisitor) VisitStructNamedCreationNode(node *node.StructNamedCreationNode) {
	structType := v.symbolTable.GetTypeByExpression(node).(*symbol.StructTypeSymbol)
	v.assembler.NewStruct(uint16(len(structType.Fields)))

	initializedFields := make([]bool, len(structType.Fields))
	for _, namedField := range node.FieldValues {
		namedField.Accept(v.ConcreteVisitor)
		fieldIndex := structType.GetFieldIndex(namedField.Name)
		v.assembler.StoreField(uint16(fieldIndex))
		initializedFields[fieldIndex] = true
	}

	// Set default value when field is not initialized
	for i := 0; i < len(structType.Fields); i++ {
		if !initializedFields[i] {
			v.pushDefault(structType.Fields[i].Type)
			v.assembler.StoreField(uint16(i))
		}
	}
}

// VisitArrayLengthCreationNode generates the IL Code for the array length creation
func (v *ILCodeGenerationVisitor) VisitArrayLengthCreationNode(node *node.ArrayLengthCreationNode) {
	if len(node.Lengths) > 1 {
		v.reportError(node, "Generator currently does not support array nesting")
		return
	}

	node.Lengths[0].Accept(v)
	v.assembler.Emit(il.NewArr) // Pass Array length as parameter
}

// VisitArrayValueCreationNode generates the IL Code for the array value creation
func (v *ILCodeGenerationVisitor) VisitArrayValueCreationNode(n *node.ArrayValueCreationNode) {
	if _, isNested := n.Elements.Values[0].(*node.ArrayInitializationNode); isNested {
		v.reportError(n, "Generator currently does not support array nesting")
		return
	}

	length := big.NewInt(int64(len(n.Elements.Values)))
	v.assembler.PushInt(length)
	v.assembler.Emit(il.NewArr)
	for i, value := range n.Elements.Values {
		value.Accept(v)
		v.assembler.Emit(il.Swap)                 // array is be popped from stack before value
		v.assembler.PushInt(big.NewInt(int64(i))) // array is popped from stack before index
		v.assembler.Emit(il.Swap)
		v.assembler.Emit(il.ArrInsert)
	}

}

// VisitBasicDesignatorNode generates the IL Code for a designator
func (v *ILCodeGenerationVisitor) VisitBasicDesignatorNode(node *node.BasicDesignatorNode) {
	decl := v.symbolTable.GetDeclByDesignator(node)
	if node.String() != symbol.This {
		v.loadVariable(decl)
	}
}

// VisitIntegerLiteralNode pushes an integer to the stack
func (v *ILCodeGenerationVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	v.assembler.PushInt(node.Value)
}

// VisitBoolLiteralNode pushes a boolean to the stack
func (v *ILCodeGenerationVisitor) VisitBoolLiteralNode(node *node.BoolLiteralNode) {
	v.assembler.PushBool(node.Value)
}

// VisitStringLiteralNode pushes a string to the stack
func (v *ILCodeGenerationVisitor) VisitStringLiteralNode(node *node.StringLiteralNode) {
	v.assembler.PushString(node.Value)
}

// VisitCharacterLiteralNode pushes a character to the stack
func (v *ILCodeGenerationVisitor) VisitCharacterLiteralNode(node *node.CharacterLiteralNode) {
	v.assembler.PushCharacter(node.Value)
}

// Helper Functions
// ----------------

func (v *ILCodeGenerationVisitor) loadVariable(decl symbol.Symbol) {
	index := v.getVarIndex(decl)

	switch decl.Scope().(type) {
	case *symbol.ContractSymbol:
		v.assembler.LoadState(byte(index))
	case *symbol.FunctionSymbol:
		v.assembler.LoadLocal(byte(index))
	case *symbol.StructTypeSymbol:
		v.assembler.LoadField(uint16(index))
	}
}

func (v *ILCodeGenerationVisitor) storeVariable(decl symbol.Symbol) {
	index := v.getVarIndex(decl)

	switch decl.Scope().(type) {
	case *symbol.ContractSymbol:
		v.assembler.StoreState(byte(index))
	case *symbol.FunctionSymbol:
		v.assembler.StoreLocal(byte(index))
	case *symbol.StructTypeSymbol:
		v.assembler.StoreField(uint16(index))
	}
}

func (v *ILCodeGenerationVisitor) pushDefault(typeSymbol symbol.TypeSymbol) {
	switch typeSymbol.(type) {
	case *symbol.StructTypeSymbol:
		v.pushDefaultStruct(typeSymbol.(*symbol.StructTypeSymbol))
		return
	case *symbol.ArrayTypeSymbol:
		v.assembler.PushNil()
		return
	}

	gs := v.symbolTable.GlobalScope

	switch typeSymbol {
	case gs.IntType:
		v.assembler.PushInt(big.NewInt(0))
	case gs.BoolType:
		v.assembler.PushBool(false)
	case gs.StringType:
		v.assembler.PushString("")
	case gs.CharType:
		v.assembler.PushCharacter('0')
	default:
		typeNode := v.symbolTable.GetNodeBySymbol(typeSymbol.(symbol.Symbol))
		v.reportError(typeNode, fmt.Sprintf("%s not supported", typeSymbol.Identifier()))
	}
}

func (v *ILCodeGenerationVisitor) pushDefaultStruct(structType *symbol.StructTypeSymbol) {
	v.assembler.NewStruct(uint16(len(structType.Fields)))

	for i, field := range structType.Fields {
		if _, ok := field.Type.(*symbol.StructTypeSymbol); ok {
			v.assembler.PushNil()
		} else {
			v.pushDefault(field.Type)
		}
		v.assembler.StoreField(uint16(i))
	}
}

// Returns: variable index and isContractField
func (v *ILCodeGenerationVisitor) getVarIndex(decl symbol.Symbol) int {
	switch decl.(type) {
	case *symbol.LocalVariableSymbol, *symbol.ParameterSymbol:
		index := v.function.GetVarIndex(decl.Identifier())
		if index == -1 {
			panic(fmt.Sprintf("Variable not found %s", decl.Identifier()))
		}
		return index
	case *symbol.FieldSymbol:
		scope := decl.Scope()
		index := -1
		if contract, ok := scope.(*symbol.ContractSymbol); ok {
			index = contract.GetFieldIndex(decl.Identifier())
		} else if structType, ok := scope.(*symbol.StructTypeSymbol); ok {
			index = structType.GetFieldIndex(decl.Identifier())
		} else {
			panic(fmt.Sprintf("Unsupported field scope %s", scope.Identifier()))
		}

		if index == -1 {
			panic(fmt.Sprintf("Variable not found %s", decl.Identifier()))
		}
		return index
	default:
		panic(fmt.Sprintf("Unsupported variable type %t", decl))
	}
}

func (v *ILCodeGenerationVisitor) isArray(node node.DesignatorNode) bool {
	sym := v.symbolTable.GetTypeByExpression(node)
	_, ok := sym.(*symbol.ArrayTypeSymbol)
	return ok
}

func (v *ILCodeGenerationVisitor) reportError(node node.Node, msg string) {
	v.Errors = append(v.Errors, fmt.Errorf("[%s] %s", node.Pos(), msg))
}
