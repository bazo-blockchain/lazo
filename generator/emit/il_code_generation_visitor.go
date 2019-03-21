package emit

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/il"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
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

func (v *ILCodeGenerationVisitor) VisitBinaryExpressionNode(node *node.BinaryExpressionNode) {
	if op, ok := binaryOpCodes[node.Operator]; ok {
		v.AbstractVisitor.VisitBinaryExpressionNode(node)
		v.assembler.Emit(op)
	} else {
		// TODO complete binary exp logic
		panic("binary operator not supported")
	}
}

func (v *ILCodeGenerationVisitor) VisitUnaryExpressionNode(node *node.UnaryExpression) {
	if op, ok := unaryOpCodes[node.Operator]; ok {
		v.AbstractVisitor.VisitUnaryExpressionNode(node)
		v.assembler.Emit(op)
	} else {
		panic("unary operator not supported")
	}
}

func (v *ILCodeGenerationVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(node)
	v.assembler.Emit(il.RET)
}

func (v *ILCodeGenerationVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	v.assembler.EmitOperand(il.PUSH, node.Value)
}

// Helper Functions
// ----------------

var binaryOpCodes = map[token.Symbol]il.OpCode{
	token.Addition:       il.ADD,
	token.Subtraction:    il.SUB,
	token.Multiplication: il.MULT,
	token.Division:       il.DIV,
}

var unaryOpCodes = map[token.Symbol]il.OpCode{
	token.Subtraction: il.NEG,
}
