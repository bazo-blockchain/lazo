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
		syntaxTree: syntaxTree,
	}
	return p
}

func (c *Checker) Run() (*symbol.SymbolTable, []error) {
	phase.RunSymbolConstruction(c.symbolTable, c.syntaxTree)
	if c.hasErrors() {
		return nil, c.errors
	}
	//phase.RunTypeResolution(c.symbolTable)
	//if c.hasErrors() {
	//	return nil, c.errors
	//}

	return c.symbolTable, c.errors
}

func (c *Checker) hasErrors() bool {
	return len(c.errors) > 0
}
