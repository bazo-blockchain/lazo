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
	bytePos     int
	errors      []error
}

func New(symbolTable *symbol.SymbolTable) *Generator {
	g := &Generator{
		symbolTable: symbolTable,
		ilBuilder:   emit.NewILBuilder(symbolTable),
		bytePos:     0,
	}
	g.metaData = g.ilBuilder.Metadata
	return g
}

func (g *Generator) Run() (*il.Metadata, []error) {
	contract := g.symbolTable.GlobalScope.Contract
	g.generateContractIL(contract)

	for _, function := range contract.Functions {
		g.generateFunctionIL(function)
	}
	g.ilBuilder.Complete()
	return g.metaData, g.errors
}

func (g *Generator) generateContractIL(contract *symbol.ContractSymbol) {
	// FIXME: Call constructor. Temporarily call the first function to test
	contractData := g.metaData.Contract

	assembler := emit.NewILAssembler(g.bytePos)
	assembler.Call(contract.Functions[0])
	contractData.Instructions, g.bytePos = assembler.Complete(true)
}

func (g *Generator) generateFunctionIL(function *symbol.FunctionSymbol) {
	funcData := g.ilBuilder.GetFunctionData(function)
	funcNode := g.symbolTable.GetNodeBySymbol(function).(*node.FunctionNode)
	g.ilBuilder.SetFunctionPos(function, g.bytePos)

	assembler := emit.NewILAssembler(g.bytePos)
	v := emit.NewCodeGenerationVisitor(g.symbolTable, function, assembler)
	v.VisitStatementBlock(funcNode.Body)
	funcData.Instructions, g.bytePos = assembler.Complete(false)
}
