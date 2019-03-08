package phase

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type SymbolConstruction struct {
	programNode *node.ProgramNode
	symTable *symbol.Table
}

func RunSymbolConstruction(symTable *symbol.Table, programNode *node.ProgramNode) {
	construction := SymbolConstruction{
		symTable: symTable,
		programNode: programNode,
	}
	construction.registerDeclarations(construction.programNode)
	construction.checkValidIdentifiers()
	construction.checkUniqueIdentifiers()
}

func (sc *SymbolConstruction) registerDeclarations(programNode *node.ProgramNode) {}

func (sc *SymbolConstruction) checkValidIdentifiers() {}

func (sc *SymbolConstruction) checkUniqueIdentifiers() {}
