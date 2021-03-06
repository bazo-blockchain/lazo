package designatorresolution

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type designatorResolution struct {
	symTable *symbol.SymbolTable
	errors   []error
}

// Run resolves designators to its declaration
// Returns errors that occurred during resolution
func Run(symTable *symbol.SymbolTable) []error {
	resolution := designatorResolution{
		symTable: symTable,
	}
	resolution.resolveDesignators()
	return resolution.errors
}

func (dr *designatorResolution) resolveDesignators() {
	contractSymbol := dr.symTable.GlobalScope.Contract
	v := newDesignatorResolutionVisitor(dr.symTable, contractSymbol)
	contractNode := dr.symTable.GetNodeBySymbol(contractSymbol).(*node.ContractNode)

	contractNode.Accept(v)
	dr.errors = v.Errors
}
