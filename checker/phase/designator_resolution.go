package phase

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/checker/visitor"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
)

type designatorResolution struct {
	symTable *symbol.SymbolTable
	errors []error
}


func RunDesignatorResolution(symTable *symbol.SymbolTable) []error {
	resolution :=designatorResolution{
		symTable: symTable,
	}
	resolution.resolveDesignators()
	return resolution.errors
}

func (dr *designatorResolution) resolveDesignators() {
	contractSymbol := dr.symTable.GlobalScope.Contract
	v := visitor.NewDesignatorResolutionVisitor(dr.symTable, contractSymbol)
	contractNode := dr.symTable.GetNodeBySymbol(contractSymbol).(*node.ContractNode)
	contractNode.Accept(v)
}

func (tr *TypeResolution) reportDesignatorResolutionError(sym symbol.Symbol, msg string) {
	tr.errors = append(tr.errors, errors.New(fmt.Sprintf("[%s] %s", tr.symTable.GetNodeBySymbol(sym).Pos(), msg)))
}