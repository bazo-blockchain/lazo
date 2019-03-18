package emit

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/il"
)

type ILBuilder struct {
	symbolTable *symbol.SymbolTable
	typeRefs map[string]int
	MetaData *il.MetaData
}

func NewILBuilder(symbolTable *symbol.SymbolTable) *ILBuilder {
	return &ILBuilder{
		symbolTable: symbolTable,
		MetaData: &il.MetaData{},
	}
}

func (b * ILBuilder) GenerateMetaData() {
	// Check if is Contract Symbol, if so register (1) all fields and (2) all methods
	return
}
