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
	symbolTable  *symbol.SymbolTable
	MetaData     *il.MetaData
	functionRefs map[*symbol.FunctionSymbol]int
	Errors       []error
}

func NewILBuilder(symbolTable *symbol.SymbolTable) *ILBuilder {
	builder := &ILBuilder{
		symbolTable:  symbolTable,
		MetaData:     &il.MetaData{},
		functionRefs: map[*symbol.FunctionSymbol]int{},
	}
	builder.GenerateMetaData()
	return builder
}

func (b *ILBuilder) GenerateMetaData() {
	contract := b.symbolTable.GlobalScope.Contract
	b.registerContract(contract)
	b.fixContract(contract)
}

func (b *ILBuilder) Complete() {
	for _, function := range b.MetaData.Contract.Functions {
		b.fixOperands(function.Code)
	}
}

func (b *ILBuilder) fixOperands(code []*il.Instruction) {
	for _, instruction := range code {
		if typeSymbol, ok := instruction.Operand.(*symbol.TypeSymbol); ok {
			instruction.Operand = b.getTypeRef(typeSymbol)
		} else if functionSymbol, ok := instruction.Operand.(*symbol.FunctionSymbol); ok {
			instruction.Operand = b.getFunctionRef(functionSymbol)
		}
	}
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
	functionData := &il.FunctionData{
		Identifier: function.GetIdentifier(),
	}
	b.MetaData.Contract.Functions = append(b.MetaData.Contract.Functions, functionData)
	b.functionRefs[function] = len(b.MetaData.Contract.Functions) - 1
}

func (b *ILBuilder) fixContract(contract *symbol.ContractSymbol) {
	contractData := b.MetaData.Contract

	for _, field := range contract.Fields {
		contractData.Fields = append(contractData.Fields, b.getTypeRef(field.Type))
	}

	for _, function := range contract.Functions {
		b.fixFunction(function)
	}

}

func (b *ILBuilder) fixFunction(function *symbol.FunctionSymbol) {
	functionData := b.getFunctionData(function)

	for _, rtype := range function.ReturnTypes {
		functionData.ReturnTypes = append(functionData.ReturnTypes, b.getTypeRef(rtype))
	}

	for _, param := range function.Parameters {
		functionData.ParamTypes = append(functionData.ParamTypes, b.getTypeRef(param.Type))
	}

	for _, local := range function.LocalVariables {
		functionData.LocalTypes = append(functionData.LocalTypes, b.getTypeRef(local.Type))
	}

}

func (b *ILBuilder) getFunctionRef(symbol *symbol.FunctionSymbol) int {
	return b.functionRefs[symbol]
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

func (b *ILBuilder) getFunctionData(symbol *symbol.FunctionSymbol) *il.FunctionData {
	return b.MetaData.Contract.Functions[b.getFunctionRef(symbol)]
}
