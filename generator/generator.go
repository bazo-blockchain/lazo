package generator

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/emit"
)

type Generator struct {
	symbolTable *symbol.SymbolTable
	ilBuilder *emit.ILBuilder
	errors []error
}

func New(symbolTable *symbol.SymbolTable) *Generator {
	p := &Generator{
		symbolTable: symbolTable,
		ilBuilder: emit.NewILBuilder(symbolTable),
	}
	return p
}

func (g *Generator) Run() []error {
	for _, function := range g.symbolTable.GlobalScope.Contract.Functions {
		g.generateIL(function)
	}
	// TODO Complete IL Builder

	return g.errors
}

func (g *Generator) generateIL(function *symbol.FunctionSymbol) {
	// TODO IL Assembler erstellen
	fmt.Println("Generating IL Code")
}
