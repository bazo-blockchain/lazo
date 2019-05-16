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
	v.blockScopes = v.blockScopes[:v.currentBlockIndex()] // remove last block scope
}

// VisitVariableNode records the visibility of the local variable and adds the local variable to the block scopes
func (v *localVariableVisitor) VisitVariableNode(node *node.VariableNode) {
	v.recordVisiblity(node)

	sym := symbol.NewLocalVariableSymbol(v.function, node.Identifier)
	v.function.LocalVariables = append(v.function.LocalVariables, sym)
	v.symbolTable.MapSymbolToNode(sym, node)

	v.addToScope(sym)
}

func (v *localVariableVisitor) VisitMultiVariableNode(node *node.MultiVariableNode) {
	v.recordVisiblity(node)

	for _, id := range node.Identifiers {
		sym := symbol.NewLocalVariableSymbol(v.function, id)
		v.function.LocalVariables = append(v.function.LocalVariables, sym)
		v.symbolTable.MapSymbolToNode(sym, node)

		v.addToScope(sym)
	}
}

// VisitAssignmentStatementNode records the visibility
func (v *localVariableVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.recordVisiblity(node)
}

// VisitMultiAssignmentStatementNode records the visibility
func (v *localVariableVisitor) VisitMultiAssignmentStatementNode(node *node.MultiAssignmentStatementNode) {
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

// VisitDeleteStatementNode records the visibility
func (v *localVariableVisitor) VisitDeleteStatementNode(node *node.DeleteStatementNode) {
	v.recordVisiblity(node)
}

func (v *localVariableVisitor) recordVisiblity(stmt node.StatementNode) {
	for _, scope := range v.blockScopes {
		for _, localVariable := range scope {
			localVariable.VisibleIn = append(localVariable.VisibleIn, stmt)
		}
	}
}

func (v *localVariableVisitor) currentBlockIndex() int {
	return len(v.blockScopes) - 1
}

func (v *localVariableVisitor) addToScope(sym *symbol.LocalVariableSymbol) {
	// append the local variable to the actual block scope
	index := v.currentBlockIndex()
	v.blockScopes[index] = append(v.blockScopes[index], sym)
}
