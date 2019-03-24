package emit

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/il"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"math/big"
)

type ILCodeGenerationVisitor struct {
	node.AbstractVisitor
	symbolTable *symbol.SymbolTable
	ilBuilder   *ILBuilder
	function    *symbol.FunctionSymbol
	assembler   *ILAssembler
	bytePos     uint16
	Errors      []error
}

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

func (v *ILCodeGenerationVisitor) VisitContractNode(node *node.ContractNode) {
	contractSymbol := v.symbolTable.GlobalScope.Contract
	contractData := v.ilBuilder.Metadata.Contract

	//for _, variable := range node.Variables {
	//	variable.Accept(v.ConcreteVisitor)
	//}

	v.assembler = NewILAssembler(&v.bytePos)
	v.assembler.Call(contractSymbol.Functions[0])
	contractData.Instructions = v.assembler.Complete(true)

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

func (v *ILCodeGenerationVisitor) VisitVariableNode(node *node.VariableNode) {
	v.AbstractVisitor.VisitVariableNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)

	if node.Expression == nil {
		v.pushDefault(targetType)
	}

	// TODO: Differentiate between local variable and field variable
	index := v.function.GetVarIndex(node.Identifier)
	v.assembler.Store(byte(index))
}

func (v *ILCodeGenerationVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	node.Right.Accept(v)

	decl := v.symbolTable.GetDeclByDesignator(node.Left)
	index, isContractField := v.getVarIndex(decl)

	if !isContractField {
		v.assembler.Store(index)
	}
	// TODO: Implement contract field SSTORE
}

func (v *ILCodeGenerationVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(node)
	v.assembler.Emit(il.RET)
}

func (v *ILCodeGenerationVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	elseLabel := v.assembler.CreateLabel()
	endLabel := v.assembler.CreateLabel()

	// Condition
	node.Condition.Accept(v)
	v.assembler.NegBool()
	v.assembler.JmpIfTrue(elseLabel)

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
	token.Addition:       il.ADD,
	token.Subtraction:    il.SUB,
	token.Multiplication: il.MULT,
	token.Division:       il.DIV,
	token.Modulo:         il.MOD,
	token.Exponent:       il.MULT,
	token.Greater:        il.GT,
	token.GreaterEqual:   il.GTE,
	token.LessEqual:      il.LTE,
	token.Less:           il.LT,
	token.Equal:          il.EQ,
	token.Unequal:        il.NEQ,
}

func (v *ILCodeGenerationVisitor) VisitBinaryExpressionNode(expNode *node.BinaryExpressionNode) {
	if op, ok := binaryOpCodes[expNode.Operator]; ok {
		switch expNode.Operator {
		case token.Exponent:
			exponent := expNode.Right.(*node.IntegerLiteralNode)
			left := expNode.Left.(*node.IntegerLiteralNode)
			expNode.Right = expNode.Left
			v.AbstractVisitor.VisitBinaryExpressionNode(expNode)
			for i := big.NewInt(0); lessThan(i, sub(exponent.Value, big.NewInt(2))); i = add(i, big.NewInt(1)) {
				v.assembler.Emit(op)
				v.assembler.PushInt(left.Value)
			}
			v.assembler.Emit(op)
		default:
			v.AbstractVisitor.VisitBinaryExpressionNode(expNode)
			v.assembler.Emit(op)
		}
		return
	}

	if expNode.Operator == token.And {
		falseLabel := v.assembler.CreateLabel()
		endLabel := v.assembler.CreateLabel()

		expNode.Left.Accept(v)
		v.assembler.NegBool()
		v.assembler.JmpIfTrue(falseLabel)
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
		// ConvertToBool fixes Bug in JMPIF on VM
		// VM Stores [0 1] on stack for value 1 but JMP IF only reads the first Byte
		v.assembler.ConvertToBool()
		v.assembler.JmpIfTrue(trueLabel)
		expNode.Right.Accept(v)
		v.assembler.Jmp(endLabel)

		v.assembler.SetLabel(trueLabel)
		v.assembler.PushBool(true)

		v.assembler.SetLabel(endLabel)
		return
	}

	// TODO complete binary exp logic
	panic("binary operator not supported")
}

var unaryOpCodes = map[token.Symbol]il.OpCode{
	token.Subtraction: il.NEG,
	token.Addition:    il.NOP,
}

func (v *ILCodeGenerationVisitor) VisitUnaryExpressionNode(node *node.UnaryExpression) {
	if op, ok := unaryOpCodes[node.Operator]; ok {
		v.AbstractVisitor.VisitUnaryExpressionNode(node)
		v.assembler.Emit(op)
		return
	}

	if node.Operator == token.Not {
		v.AbstractVisitor.VisitUnaryExpressionNode(node)
		v.assembler.NegBool()
		return
	}

	panic("unary operator not supported")
}

func (v *ILCodeGenerationVisitor) VisitDesignatorNode(node *node.DesignatorNode) {
	decl := v.symbolTable.GetDeclByDesignator(node)
	index, isContractField := v.getVarIndex(decl)

	if !isContractField {
		v.assembler.Load(index)
	}
	// TODO: Implement contract field SLOAD
}

func (v *ILCodeGenerationVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	v.assembler.PushInt(node.Value)
}

func (v *ILCodeGenerationVisitor) VisitBoolLiteralNode(node *node.BoolLiteralNode) {
	v.assembler.PushBool(node.Value)
}

func (v *ILCodeGenerationVisitor) VisitStringLiteralNode(node *node.StringLiteralNode) {
	// TODO Implement
	v.AbstractVisitor.VisitStringLiteralNode(node)
}

func (v *ILCodeGenerationVisitor) VisitCharacterLiteralNode(node *node.CharacterLiteralNode) {
	// TODO Implement
	v.AbstractVisitor.VisitCharacterLiteralNode(node)
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
	default:
		panic(fmt.Sprintf("%s not supported", typeSymbol.Identifier))
	}
}

// Returns: variable index and isContractField
func (v *ILCodeGenerationVisitor) getVarIndex(decl symbol.Symbol) (byte, bool) {
	switch decl.(type) {
	case *symbol.LocalVariableSymbol:
		index := v.function.GetVarIndex(decl.GetIdentifier())
		if index == -1 {
			panic(fmt.Sprintf("Variable not found %s", decl.GetIdentifier()))
		}
		return byte(index), false
	default:
		panic("Not implemented")
	}
}

func lessThan(x *big.Int, y *big.Int) bool {
	value := x.Cmp(y) == -1
	return value
}

func sub(x *big.Int, y *big.Int) *big.Int {
	value := big.NewInt(0).Sub(x, y)
	valstr := value.String()
	fmt.Println(valstr)
	return value
}

func add(x *big.Int, y *big.Int) *big.Int {
	value := big.NewInt(0).Add(x, y)
	return value
}
