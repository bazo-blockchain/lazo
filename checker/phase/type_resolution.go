package phase

import "github.com/bazo-blockchain/lazo/checker/symbol"

type TypeResolution struct {
	symTable *symbol.SymbolTable
}

func RunTypeResolution(symTable *symbol.SymbolTable) {
	resolution :=TypeResolution{
		symTable: symTable,
	}
	resolution.resolveTypesInSymbols()
}

func (tr *TypeResolution) resolveTypesInSymbols() {}