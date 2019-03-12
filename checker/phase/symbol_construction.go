package phase

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type SymbolConstruction struct {
	programNode *node.ProgramNode
	symTable *symbol.Table
	globalScope *symbol.CompilationUnit
}

func RunSymbolConstruction(symTable *symbol.Table, programNode *node.ProgramNode) {
	construction := SymbolConstruction{
		symTable: symTable,
		programNode: programNode,
		globalScope: symTable.Compilation,
	}
	construction.registerDeclarations(construction.programNode)
	construction.checkValidIdentifiers()
	construction.checkUniqueIdentifiers()
}

func (sc *SymbolConstruction) registerDeclarations(programNode *node.ProgramNode) {
	sc.registerBuiltins()
	sc.registerContract(programNode.Contract)
}

func (sc *SymbolConstruction) registerBuiltins() {
	sc.registerBuiltinTypes()
	sc.registerBuiltinConstants()
}

func (sc *SymbolConstruction) registerContract(contractNode *node.ContractNode) {
	sym := symbol.ContractSymbol{}.NewSymbol(sc.globalScope, contractNode.Name)

	contractSymbol, _ := sym.(*symbol.ContractSymbol)
	sc.globalScope.Contract = contractSymbol

	typeSymbol, _ := sym.(*symbol.TypeSymbol)
	sc.globalScope.Types = append(sc.globalScope.Types, typeSymbol)

	sc.symTable.LinkDeclaration(contractNode, contractSymbol)
	for _, variableNode := range contractNode.Variables {
		sc.registerField(contractSymbol, variableNode)
	}

	for _, functionNode := range contractNode.Functions {
		sc.registerFunction(contractSymbol, functionNode)
	}
}

func (sc *SymbolConstruction) registerField(contractSymbol *symbol.ContractSymbol, node *node.VariableNode) {
	sym := symbol.FieldSymbol{}.NewSymbol(contractSymbol, node.Identifier)
	fieldSymbol, _ := sym.(*symbol.FieldSymbol)
	contractSymbol.Fields = append(contractSymbol.Fields, fieldSymbol)
	sc.symTable.LinkDeclaration(node, fieldSymbol)
}

func (sc *SymbolConstruction) registerFunction(contractSymbol *symbol.ContractSymbol, node *node.FunctionNode) {
	sym := symbol.FunctionSymbol{}.NewSymbol(contractSymbol, node.Name)
	functionSymbol, _ := sym.(*symbol.FunctionSymbol)
	contractSymbol.Functions = append(contractSymbol.Functions, functionSymbol)$
	sc.symTable.LinkDeclaration(node, functionSymbol)
	for _, parameter := range node.Parameters {
		sc.registerParameter(functionSymbol, parameter)
	}
	for _, statement := range node.Body {
		// TODO Pass visitor
		statement.Accept(nil)
	}
}

func (sc *SymbolConstruction) registerParameter(functionSymbol *symbol.FunctionSymbol, node *node.VariableNode) {
	sym := symbol.ParameterSymbol{}.NewSymbol(functionSymbol, node.Identifier)
	parameterSymbol, _ := sym.(*symbol.ParameterSymbol)
	functionSymbol.Parameters = append(functionSymbol.Parameters, parameterSymbol)
	sc.symTable.LinkDeclaration(node, parameterSymbol)
}

func (sc *SymbolConstruction) registerBuiltinTypes() {
	// TODO Implement
}

func (sc *SymbolConstruction) registerBuiltinConstants() {
	// TODO Implement
}

func (sc *SymbolConstruction) checkValidIdentifiers() {
	// TODO Implement
}

func (sc *SymbolConstruction) checkUniqueIdentifiers() {
	// TODO Implement
}
