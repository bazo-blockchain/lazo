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
		v.assembler.Call(contractSymbol.Functions[i])
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

	// constructor code comes here
	v.assembler.Call(contractSymbol.Functions[0])
	contractData.Instructions = v.assembler.Complete(true)
}

// VisitFieldNode generates the IL Code for a contract field node and default initializes it if required
func (v *ILCodeGenerationVisitor) VisitFieldNode(node *node.FieldNode) {
	v.AbstractVisitor.VisitFieldNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)

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

	// parameter symbol should not be default initialized
	if node.Expression == nil && v.function != nil && !v.function.IsLocalVar(node.Identifier) {
		return
	}

	if node.Expression == nil {
		v.pushDefault(targetType)
	}

	index := v.function.GetVarIndex(node.Identifier)
	isContractField := !v.function.IsLocalVar(node.Identifier)
	if isContractField {
		v.assembler.StoreState(byte(index))
	} else {
		v.assembler.StoreLocal(byte(index))
	}
}

// VisitAssignmentStatementNode generates the IL Code for an assignment
func (v *ILCodeGenerationVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	node.Right.Accept(v)

	decl := v.symbolTable.GetDeclByDesignator(node.Left)
	index, isContractField := v.getVarIndex(decl)

	if isContractField {
		v.assembler.StoreState(index)
	} else {
		v.assembler.StoreLocal(index)
	}
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
	token.Addition:       il.Add,
	token.Subtraction:    il.Sub,
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
	token.Subtraction: il.Neg,
}

// VisitUnaryExpressionNode generates the IL Code for all unary expressions
func (v *ILCodeGenerationVisitor) VisitUnaryExpressionNode(expNode *node.UnaryExpressionNode) {
	if op, ok := unaryOpCodes[expNode.Operator]; ok {
		v.AbstractVisitor.VisitUnaryExpressionNode(expNode)
		v.assembler.Emit(op)
		return
	}

	if expNode.Operator == token.Addition {
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

// VisitDesignatorNode generates the IL Code for a designator
func (v *ILCodeGenerationVisitor) VisitDesignatorNode(node *node.DesignatorNode) {
	decl := v.symbolTable.GetDeclByDesignator(node)
	index, isContractField := v.getVarIndex(decl)

	if isContractField {
		v.assembler.LoadState(index)
	} else {
		v.assembler.LoadLocal(index)
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

func (v *ILCodeGenerationVisitor) pushDefault(typeSymbol *symbol.TypeSymbol) {
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
		typeNode := v.symbolTable.GetNodeBySymbol(typeSymbol)
		v.reportError(typeNode, fmt.Sprintf("%s not supported", typeSymbol.ID))
	}
}

// Returns: variable index and isContractField
func (v *ILCodeGenerationVisitor) getVarIndex(decl symbol.Symbol) (byte, bool) {
	switch decl.(type) {
	case *symbol.LocalVariableSymbol, *symbol.ParameterSymbol:
		index := v.function.GetVarIndex(decl.Identifier())
		if index == -1 {
			panic(fmt.Sprintf("Variable not found %s", decl.Identifier()))
		}
		return byte(index), false
	case *symbol.FieldSymbol:
		contract := decl.Scope().(*symbol.ContractSymbol)
		index := contract.GetFieldIndex(decl.Identifier())
		if index == -1 {
			panic(fmt.Sprintf("Variable not found %s", decl.Identifier()))
		}
		return byte(index), true
	default:
		panic("Not implemented")
	}
}

func (v *ILCodeGenerationVisitor) reportError(node node.Node, msg string) {
	v.Errors = append(v.Errors, fmt.Errorf("[%s] %s", node.Pos(), msg))
}
