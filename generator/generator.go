package generator

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/emit"
	"github.com/bazo-blockchain/lazo/generator/il"
)

type Generator struct {
	symbolTable *symbol.SymbolTable
	ilBuilder   *emit.ILBuilder
	metaData    *il.MetaData
	errors      []error
}

func New(symbolTable *symbol.SymbolTable) *Generator {
	g := &Generator{
		symbolTable: symbolTable,
		ilBuilder:   emit.NewILBuilder(symbolTable),
	}
	g.metaData = g.ilBuilder.MetaData
	return g
}

func (g *Generator) Run() (*il.MetaData, []error) {
	for _, function := range g.symbolTable.GlobalScope.Contract.Functions {
		g.generateIL(function)
	}
	// TODO Check how Metadata.Functions.Code is set
	g.ilBuilder.Complete()
	return g.metaData, g.errors
}

func (g *Generator) generateIL(function *symbol.FunctionSymbol) {
	fmt.Println("Generating IL Code")
}
