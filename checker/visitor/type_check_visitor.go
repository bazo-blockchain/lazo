package visitor

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type TypeCheckVisitor struct {
	node.AbstractVisitor
	symbolTable *symbol.SymbolTable
	functionSymbol *symbol.FunctionSymbol
}

func NewTypeCheckVisitor(symbolTable *symbol.SymbolTable, functionSymbol *symbol.FunctionSymbol) *TypeCheckVisitor {
	v := &TypeCheckVisitor{
		symbolTable: symbolTable,
		functionSymbol: functionSymbol,
	}
	v.ConcreteVisitor = v
	return v
}

func (v *TypeCheckVisitor) VisitTypeNode(node *node.TypeNode) {
	// To be done as soon as own types are introduced
}

func (v *TypeCheckVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
}

func (v *TypeCheckVisitor) VisitVariableNode(node *node.VariableNode) {
	v.AbstractVisitor.VisitVariableNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)
	expType := v.symbolTable.FindTypeByExpressionNode(node.Expression)
	if targetType != expType {
		fmt.Print("Error Type mismatch")
	}
}