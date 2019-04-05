package generator

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/data"
	"github.com/bazo-blockchain/lazo/generator/emit"
)

// Generator contains the symbol table, il builder and errors. It is the last step of the front-end compiler and
// generates the IL Instruction and transforms it to byte code which can then be interpreted by the VM.
type Generator struct {
	symbolTable *symbol.SymbolTable
	ilBuilder   *emit.ILBuilder
	errors      []error
}

// New returns a new Generator
func New(symbolTable *symbol.SymbolTable) *Generator {
	g := &Generator{
		symbolTable: symbolTable,
		ilBuilder:   emit.NewILBuilder(symbolTable),
	}
	return g
}

// Run performs the IL code generation process
func (g *Generator) Run() (*data.Metadata, []error) {
	contractNode := g.symbolTable.GetNodeBySymbol(g.symbolTable.GlobalScope.Contract)

	v := emit.NewCodeGenerationVisitor(g.symbolTable, g.ilBuilder)
	contractNode.Accept(v)
	g.errors = v.Errors
	metadata := g.ilBuilder.Complete()

	return metadata, g.errors
}
