package symbolconstruction

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

// LocalVariableVisitor contains the symbol table, the function and the block scopes. It traverses the abstract syntax
// tree to record the variables visibility
type LocalVariableVisitor struct {
	node.AbstractVisitor
	symbolTable *symbol.SymbolTable
	function    *symbol.FunctionSymbol
	blockScopes [][]*symbol.LocalVariableSymbol
}

// NewLocalVariableVisitor creates a new LocalVariableVisitor
func NewLocalVariableVisitor(symbolTable *symbol.SymbolTable, function *symbol.FunctionSymbol) *LocalVariableVisitor {
	v := &LocalVariableVisitor{
		symbolTable: symbolTable,
		function:    function,
	}
	v.ConcreteVisitor = v
	return v
}

// VisitStatementBlock adds a new block scope, visits the statements and removes the last blockscope as otherwise
// the variable will be visible to all statements.
func (v *LocalVariableVisitor) VisitStatementBlock(stmts []node.StatementNode) {
	v.blockScopes = append(v.blockScopes, []*symbol.LocalVariableSymbol{}) // add new block scope
	v.AbstractVisitor.VisitStatementBlock(stmts)
	v.blockScopes = v.blockScopes[:len(v.blockScopes)-1] // remove last block scope
}

// VisitVariableNode records the visibility of the local variable and adds the local variable to the block scopes
func (v *LocalVariableVisitor) VisitVariableNode(node *node.VariableNode) {
	v.recordVisiblity(node)

	sym := symbol.NewLocalVariableSymbol(v.function, node.Identifier)
	v.function.LocalVariables = append(v.function.LocalVariables, sym)
	v.symbolTable.MapSymbolToNode(sym, node)

	// append the local variable to the actual block scope
	v.blockScopes[len(v.blockScopes)-1] = append(v.blockScopes[len(v.blockScopes)-1], sym)
}

// VisitAssignmentStatementNode records the visibility
func (v *LocalVariableVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.recordVisiblity(node)
}

// VisitIfStatementNode records the visibility
func (v *LocalVariableVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	v.recordVisiblity(node)
	v.AbstractVisitor.VisitIfStatementNode(node)
}

// VisitReturnStatementNode records the visibility
func (v *LocalVariableVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.recordVisiblity(node)
}

func (v *LocalVariableVisitor) recordVisiblity(stmt node.StatementNode) {
	for _, scope := range v.blockScopes {
		for _, localVariable := range scope {
			localVariable.VisibleIn = append(localVariable.VisibleIn, stmt)
		}
	}
}
