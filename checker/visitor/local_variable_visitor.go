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

func (v *LocalVariableVisitor) VisitFunctionNode(node *node.FunctionNode) {
	// create stack
	for _, statement := range node.Body {
		statement.Accept(v.ConcreteVisitor)
	}
	// stack pop
}

func (v *LocalVariableVisitor) VisitVariableNode(node *node.VariableNode) {
	sym := symbol.NewLocalVariableSymbol(v.function, node.Identifier)
	v.function.LocalVariables = append(v.function.LocalVariables, sym)
	v.symbolTable.MapSymbolToNode(sym, node)

	fmt.Println("variable " + node.Identifier)
	// add symbol to the stack
}

func (v *LocalVariableVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	// Record
	v.AbstractVisitor.VisitIfStatementNode(node)
}


