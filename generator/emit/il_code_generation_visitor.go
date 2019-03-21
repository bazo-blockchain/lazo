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
				v.assembler.EmitOperand(il.PUSH, left.Value)
			}
			v.assembler.Emit(op)
		default:
			v.AbstractVisitor.VisitBinaryExpressionNode(expNode)
			v.assembler.Emit(op)
		}
		return
	}

	if expNode.Operator == token.And {
		returnFalseLabel := v.assembler.CreateLabel()
		skipLabel := v.assembler.CreateLabel()
		expNode.Left.Accept(v)
		v.assembler.Emit(il.NEG)
		v.assembler.EmitOperand(il.JMPIF, returnFalseLabel)
		expNode.Right.Accept(v)
		v.assembler.Emit(il.NEG)
		v.assembler.EmitOperand(il.JMPIF, returnFalseLabel)
		// Load constant boolean true
		v.assembler.EmitOperand(il.PUSH, 1)
		v.assembler.EmitOperand(il.JMP, skipLabel)
		v.assembler.SetLabel(returnFalseLabel)
		// Load constant boolean false
		v.assembler.EmitOperand(il.PUSH, 0)
		v.assembler.SetLabel(skipLabel)
		return
	}

	if expNode.Operator == token.Or {
		returnTrueLabel := v.assembler.CreateLabel()
		skipLabel := v.assembler.CreateLabel()
		expNode.Left.Accept(v)
		v.assembler.EmitOperand(il.JMPIF, returnTrueLabel)
		expNode.Right.Accept(v)
		v.assembler.EmitOperand(il.JMPIF, returnTrueLabel)
		// Load constant boolean false
		v.assembler.EmitOperand(il.PUSH, 1)
		v.assembler.EmitOperand(il.JMP, skipLabel)
		v.assembler.SetLabel(returnTrueLabel)
		// Load constant boolean true
		v.assembler.EmitOperand(il.PUSH, 1)
		v.assembler.SetLabel(skipLabel)
		return
	}

	// TODO complete binary exp logic
	panic("binary operator not supported")
}

func (v *ILCodeGenerationVisitor) VisitUnaryExpressionNode(node *node.UnaryExpression) {
	if op, ok := unaryOpCodes[node.Operator]; ok {
		v.AbstractVisitor.VisitUnaryExpressionNode(node)
		v.assembler.Emit(op)
	} else {
		panic("unary operator not supported")
	}
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
	v.assembler.Emit(il.HALT)
}

func (v *ILCodeGenerationVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	v.assembler.PushInt(node.Value)
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
