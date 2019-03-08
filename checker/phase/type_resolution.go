package phase

import "github.com/bazo-blockchain/lazo/checker/symbol"

type TypeResolution struct {
	symTable *symbol.Table
}

func RunTypeResolution(symTable *symbol.Table) {
	TypeResolution{
		symTable: symTable,
	}.resolveTypesInSymbols()
}

func (tr *TypeResolution) resolveTypesInSymbols() {}