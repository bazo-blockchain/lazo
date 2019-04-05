package designatorresolution

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type designatorResolution struct {
	symTable *symbol.SymbolTable
	errors   []error
}

// RunDesignatorResolution resolves designators to its declaration
// Returns errors that occurred during resolution
func RunDesignatorResolution(symTable *symbol.SymbolTable) []error {
	resolution := designatorResolution{
		symTable: symTable,
	}
	resolution.resolveDesignators()
	return resolution.errors
}

func (dr *designatorResolution) resolveDesignators() {
	contractSymbol := dr.symTable.GlobalScope.Contract
	v := NewDesignatorResolutionVisitor(dr.symTable, contractSymbol)
	contractNode := dr.symTable.GetNodeBySymbol(contractSymbol).(*node.ContractNode)

	contractNode.Accept(v)
	dr.errors = v.Errors
}
