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
	function    *symbol.FunctionSymbol
	assembler   *ILAssembler
	Errors      []error
}

func NewCodeGenerationVisitor(
	symbolTable *symbol.SymbolTable, function *symbol.FunctionSymbol, assembler *ILAssembler) *ILCodeGenerationVisitor {
	v := &ILCodeGenerationVisitor{
		symbolTable: symbolTable,
		function:    function,
		assembler:   assembler,
	}
	v.ConcreteVisitor = v
	return v
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

func (v *ILCodeGenerationVisitor) VisitAssignmentNode(node *node.AssignmentStatementNode) {
	// TODO Implement
	v.AbstractVisitor.VisitAssignmentStatementNode(node)
}

func (v *ILCodeGenerationVisitor) VisitDesignatorNode(node *node.DesignatorNode) {
	// TODO Implement
	v.AbstractVisitor.VisitDesignatorNode(node)
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

func (v *ILCodeGenerationVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(node)
	v.assembler.Emit(il.RET)
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

var unaryOpCodes = map[token.Symbol]il.OpCode{
	token.Subtraction: il.NEG,
	token.Addition:    il.NOP,
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
