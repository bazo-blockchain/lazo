package phase

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/checker/visitor"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type designatorResolution struct {
	symTable *symbol.SymbolTable
}


func RunDesignatorResolution(symTable *symbol.SymbolTable) {
	resolution :=designatorResolution{
		symTable: symTable,
	}
	resolution.resolveDesignators()
}

func (dr *designatorResolution) resolveDesignators() {
	contractSymbol := dr.symTable.GlobalScope.Contract
	v := visitor.NewDesignatorResolutionVisitor(dr.symTable, contractSymbol)
	contractNode := dr.symTable.GetNodeBySymbol(contractSymbol).(*node.ContractNode)
	contractNode.Accept(v)
}