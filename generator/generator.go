package generator

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/emit"
	"github.com/bazo-blockchain/lazo/generator/il"
	"github.com/bazo-blockchain/lazo/parser/node"
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
	// call first function
	// Halt
	for _, function := range g.symbolTable.GlobalScope.Contract.Functions {
		g.generateIL(function)
	}
	// TODO Check how Metadata.Functions.Code is set
	g.ilBuilder.Complete()
	return g.metaData, g.errors
}

func (g *Generator) generateIL(function *symbol.FunctionSymbol) {
	funcData := g.ilBuilder.GetFunctionData(function)
	funcNode := g.symbolTable.GetNodeBySymbol(function).(*node.FunctionNode)

	assembler := emit.NewILAssembler(0)
	v := emit.NewCodeGenerationVisitor(g.symbolTable, function, assembler)
	v.VisitStatementBlock(funcNode.Body)
	funcData.Instructions = assembler.Complete()
}
