package phase

import "github.com/bazo-blockchain/lazo/checker/symbol"

type TypeResolution struct {
	symTable *symbol.Table
}

func RunTypeResolution(symTable *symbol.Table) {
	resolution :=TypeResolution{
		symTable: symTable,
	}
	resolution.resolveTypesInSymbols()
}

func (tr *TypeResolution) resolveTypesInSymbols() {}