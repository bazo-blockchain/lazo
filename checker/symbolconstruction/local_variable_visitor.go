package symbolconstruction

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type localVariableVisitor struct {
	node.AbstractVisitor
	symbolTable *symbol.SymbolTable
	function    *symbol.FunctionSymbol
	blockScopes [][]*symbol.LocalVariableSymbol
}

func newLocalVariableVisitor(symbolTable *symbol.SymbolTable, function *symbol.FunctionSymbol) *localVariableVisitor {
	v := &localVariableVisitor{
		symbolTable: symbolTable,
		function:    function,
	}
	v.ConcreteVisitor = v
	return v
}

// VisitStatementBlock adds a new block scope, visits the statements and removes the last blockscope as otherwise
// the variable will be visible to all statements.
func (v *localVariableVisitor) VisitStatementBlock(stmts []node.StatementNode) {
	v.blockScopes = append(v.blockScopes, []*symbol.LocalVariableSymbol{}) // add new block scope
	v.AbstractVisitor.VisitStatementBlock(stmts)
	v.blockScopes = v.blockScopes[:len(v.blockScopes)-1] // remove last block scope
}

// VisitVariableNode records the visibility of the local variable and adds the local variable to the block scopes
func (v *localVariableVisitor) VisitVariableNode(node *node.VariableNode) {
	v.recordVisiblity(node)

	sym := symbol.NewLocalVariableSymbol(v.function, node.Identifier)
	v.function.LocalVariables = append(v.function.LocalVariables, sym)
	v.symbolTable.MapSymbolToNode(sym, node)

	// append the local variable to the actual block scope
	v.blockScopes[len(v.blockScopes)-1] = append(v.blockScopes[len(v.blockScopes)-1], sym)
}

// VisitAssignmentStatementNode records the visibility
func (v *localVariableVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.recordVisiblity(node)
}

// VisitIfStatementNode records the visibility
func (v *localVariableVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	v.recordVisiblity(node)
	v.AbstractVisitor.VisitIfStatementNode(node)
}

// VisitReturnStatementNode records the visibility
func (v *localVariableVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.recordVisiblity(node)
}

func (v *localVariableVisitor) recordVisiblity(stmt node.StatementNode) {
	for _, scope := range v.blockScopes {
		for _, localVariable := range scope {
			localVariable.VisibleIn = append(localVariable.VisibleIn, stmt)
		}
	}
}
