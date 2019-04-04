package checker

import (
	"github.com/bazo-blockchain/lazo/checker/phase"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

// Checker contains the syntax tree, symbol table and errors. It performs semantic analysis on abstract syntax tree and
// generates a symbol table which holds meta information about each node (e.g. scope).
type Checker struct {
	syntaxTree  *node.ProgramNode
	symbolTable *symbol.SymbolTable
	errors      []error
}

// New creates a new Checker
func New(syntaxTree *node.ProgramNode) *Checker {
	p := &Checker{
		syntaxTree: syntaxTree,
	}
	return p
}

// Run performs all the checker phases
// Executed phases are symbol construction, type resolution, designator resolution and type checking.
// If errors occur during one of those phases, the process is stopped at the end of the failing phase.
// Returns the symbol table and errors
func (c *Checker) Run() (*symbol.SymbolTable, []error) {
	c.symbolTable, c.errors = phase.RunSymbolConstruction(c.syntaxTree)
	if !c.hasErrors() {
		c.errors = phase.RunTypeResolution(c.symbolTable)
	}
	if !c.hasErrors() {
		c.errors = phase.RunDesignatorResolution(c.symbolTable)
	}
	if !c.hasErrors() {
		c.errors = phase.RunTypeChecker(c.symbolTable)
	}
	return c.symbolTable, c.errors
}

func (c *Checker) hasErrors() bool {
	return len(c.errors) > 0
}
