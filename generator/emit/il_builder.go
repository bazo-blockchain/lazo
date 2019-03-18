package emit

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/il"
)

/**
 *	IL Builder constructs Metadata
 */
type ILBuilder struct {
	symbolTable *symbol.SymbolTable
	MetaData *il.MetaData
	Errors []error
}

func NewILBuilder(symbolTable *symbol.SymbolTable) *ILBuilder {
	builder := &ILBuilder{
		symbolTable: symbolTable,
		MetaData: &il.MetaData{},
	}
	builder.GenerateMetaData()
	return builder
}

func (b *ILBuilder) GenerateMetaData() {
	contract := b.symbolTable.GlobalScope.Contract
	b.registerContract(contract)
	b.fixContract(contract)
	return
}

func (b *ILBuilder) registerContract(contract *symbol.ContractSymbol) {
	b.MetaData.Contract = &il.ContractData{
		Identifier: contract.GetIdentifier(),
	}
	for _, function := range contract.Functions {
		b.registerFunction(function)
	}
}

func (b *ILBuilder) registerFunction(function *symbol.FunctionSymbol) {
	_ = &il.FunctionData{
		Identifier: function.GetIdentifier(),
	}
}

func (b *ILBuilder) fixContract(contract *symbol.ContractSymbol) {
	contractData := b.MetaData.Contract

	for _, field := range contract.Fields {
		contractData.Fields = append(contractData.Fields, b.getTypeRef(field.Type))
	}

	for _ = range contract.Functions {

	}

}

func (b *ILBuilder) getTypeRef(sym *symbol.TypeSymbol) il.TypeData {
	scope := b.symbolTable.GlobalScope
	if sym.GetIdentifier() == scope.BoolType.GetIdentifier() {
		return il.BoolType
	} else if sym.GetIdentifier() == scope.CharType.GetIdentifier() {
		return il.CharType
	} else if sym.GetIdentifier() == scope.StringType.GetIdentifier() {
		return il.StringType
	} else if sym.GetIdentifier() == scope.IntType.GetIdentifier() {
		return il.IntType
	} else {
		panic(fmt.Sprintf("Error: Unsupported Type %s", sym.GetIdentifier()))
	}
}
