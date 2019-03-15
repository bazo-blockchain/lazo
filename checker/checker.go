package checker

import (
	"github.com/bazo-blockchain/lazo/checker/phase"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type Checker struct {
	syntaxTree  *node.ProgramNode
	symbolTable *symbol.SymbolTable
	errors      []error
}

func New(syntaxTree *node.ProgramNode) *Checker {
	p := &Checker{
		syntaxTree:  syntaxTree,
	}
	return p
}

func (c *Checker) Run() (*symbol.SymbolTable, []error) {
	c.symbolTable, c.errors = phase.RunSymbolConstruction(c.syntaxTree)
	if !c.hasErrors() {
		phase.RunTypeResolution(c.symbolTable)
	}

	if !c.hasErrors() {
		phase.RunDesignatorResolution(c.symbolTable)
	}
	if !c.hasErrors() {
		phase.RunTypeChecker(c.symbolTable)
	}
	return c.symbolTable, c.errors
}

func (c *Checker) hasErrors() bool {
	return len(c.errors) > 0
}
