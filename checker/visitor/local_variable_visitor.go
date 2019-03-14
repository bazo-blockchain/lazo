package visitor

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type LocalVariableVisitor struct {
	node.AbstractVisitor
	symbolTable *symbol.SymbolTable
	function    *symbol.FunctionSymbol
	blockScopes [][]*symbol.LocalVariableSymbol
}

func NewLocalVariableVisitor(symbolTable *symbol.SymbolTable, function *symbol.FunctionSymbol) *LocalVariableVisitor {
	v := &LocalVariableVisitor{
		symbolTable: symbolTable,
		function:    function,
	}
	v.ConcreteVisitor = v
	return v
}

func (v *LocalVariableVisitor) VisitStatementBlock(stmts []node.StatementNode) {
	v.blockScopes = append(v.blockScopes, []*symbol.LocalVariableSymbol{}) // add new block scope
	v.AbstractVisitor.VisitStatementBlock(stmts)
	v.blockScopes = v.blockScopes[:len(v.blockScopes)-1] // remove last block scope
}

func (v *LocalVariableVisitor) VisitVariableNode(node *node.VariableNode) {
	sym := symbol.NewLocalVariableSymbol(v.function, node.Identifier)
	v.function.LocalVariables = append(v.function.LocalVariables, sym)
	v.symbolTable.MapSymbolToNode(sym, node)

	// append the local variable to the actual block scope
	v.blockScopes[len(v.blockScopes)-1] = append(v.blockScopes[len(v.blockScopes)-1], sym)
}

func (v *LocalVariableVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.recordVisiblity(node)
}

func (v *LocalVariableVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	v.recordVisiblity(node)
	v.AbstractVisitor.VisitIfStatementNode(node)
}

func (v *LocalVariableVisitor) recordVisiblity(stmt node.StatementNode) {
	for _, scope := range v.blockScopes {
		for _, localVariable := range scope {
			localVariable.VisibleIn = append(localVariable.VisibleIn, stmt)
		}
	}
}
