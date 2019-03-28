package generator

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/data"
	"github.com/bazo-blockchain/lazo/generator/emit"
)

type Generator struct {
	symbolTable *symbol.SymbolTable
	ilBuilder   *emit.ILBuilder
	errors      []error
}

func New(symbolTable *symbol.SymbolTable) *Generator {
	g := &Generator{
		symbolTable: symbolTable,
		ilBuilder:   emit.NewILBuilder(symbolTable),
	}
	return g
}

func (g *Generator) Run() (*data.Metadata, []error) {
	contractNode := g.symbolTable.GetNodeBySymbol(g.symbolTable.GlobalScope.Contract)

	v := emit.NewCodeGenerationVisitor(g.symbolTable, g.ilBuilder)
	contractNode.Accept(v)
	g.errors = v.Errors
	metadata := g.ilBuilder.Complete()

	return metadata, g.errors
}
