package checker

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/phase"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type Checker struct {
	syntaxTree *node.ProgramNode
	currentNode node.Node
	symTable *symbol.Table
	errors []error
}

func New(syntaxTree *node.ProgramNode) *Checker {
	p := &Checker{
		syntaxTree: syntaxTree,
		currentNode: syntaxTree,
	}
	return p
}

func (check *Checker) Run() (*symbol.Table, []error) {
	fmt.Print(check.syntaxTree)
	phase.RunSymbolConstruction(check.symTable, check.syntaxTree)
	if check.hasErrors() {
		return nil, check.errors
	}
	phase.RunTypeResolution(check.symTable)
	if check.hasErrors() {
		return nil, check.errors
	}

	return nil, check.errors
}

func (check *Checker) hasErrors() bool {
	return len(check.errors) > 0
}