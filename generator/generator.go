package generator

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/emit"
	"github.com/bazo-blockchain/lazo/generator/il"
)

type Generator struct {
	symbolTable *symbol.SymbolTable
	ilBuilder   *emit.ILBuilder
	metaData    *il.Metadata
	errors      []error
}

func New(symbolTable *symbol.SymbolTable) *Generator {
	g := &Generator{
		symbolTable: symbolTable,
		ilBuilder:   emit.NewILBuilder(symbolTable),
	}
	g.metaData = g.ilBuilder.Metadata
	return g
}

func (g *Generator) Run() (*il.Metadata, []error) {
	for _, function := range g.symbolTable.GlobalScope.Contract.Functions {
		g.generateIL(function)
	}
	// TODO Check how Metadata.Functions.Code is set
	g.ilBuilder.Complete()
	return g.metaData, g.errors
}

func (g *Generator) generateIL(function *symbol.FunctionSymbol) {
	funcData := g.ilBuilder.GetFunctionData(function)
	_ = g.symbolTable.GetNodeBySymbol(function)

	assembler := emit.NewILAssembler(funcData)
	// visit function body
	assembler.Complete()
}
