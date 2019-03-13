package visitor

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type LocalVariableVisitor struct {
	node.AbstractVisitor
	symbolTable     *symbol.SymbolTable
	function        *symbol.FunctionSymbol
}

func NewLocalVariableVisitor(symbolTable *symbol.SymbolTable, function *symbol.FunctionSymbol) *LocalVariableVisitor {
	v := &LocalVariableVisitor{
		symbolTable: symbolTable,
		function:    function,
	}
	v.ConcreteVisitor = v
	return v
}

func (v *LocalVariableVisitor) VisitVariableNode(node *node.VariableNode) {
	fmt.Println("variable " + node.Identifier)
}
