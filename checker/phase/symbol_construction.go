package phase

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type symbolConstruction struct {
	programNode *node.ProgramNode
	symTable    *symbol.SymbolTable
	globalScope *symbol.CompilationUnit
}

func RunSymbolConstruction(symTable *symbol.SymbolTable, programNode *node.ProgramNode) {
	construction := symbolConstruction{
		symTable:    symTable,
		programNode: programNode,
		globalScope: symTable.GlobalScope,
	}
	construction.registerDeclarations()
	construction.checkValidIdentifiers()
	construction.checkUniqueIdentifiers()
}

func (sc *symbolConstruction) registerDeclarations() {
	sc.registerBuiltins()
	sc.registerContract()
}

func (sc *symbolConstruction) registerBuiltins() {
	sc.registerBuiltInTypes()
	sc.registerBuiltInConstants()
}

func (sc *symbolConstruction) registerBuiltInTypes() {
	sc.globalScope.BoolType = sc.registerBuiltInType("bool")
	sc.globalScope.CharType = sc.registerBuiltInType("char")
	sc.globalScope.IntType = sc.registerBuiltInType("int")
	sc.globalScope.StringType = sc.registerBuiltInType("string")
}

func (sc *symbolConstruction) registerBuiltInType(name string) *symbol.TypeSymbol {
	baseType := symbol.NewTypeSymbol(sc.globalScope, name)
	sc.globalScope.Types = append(sc.globalScope.Types, baseType)
	sc.globalScope.BuiltInTypes = append(sc.globalScope.BuiltInTypes, baseType)
	return baseType
}

func (sc *symbolConstruction) registerBuiltInConstants() {
	sc.globalScope.NullConstant = sc.registerBuiltInConstant(sc.globalScope.NullType, "null")
	sc.globalScope.FalseConstant = sc.registerBuiltInConstant(sc.globalScope.BoolType, "false")
	sc.globalScope.TrueConstant = sc.registerBuiltInConstant(sc.globalScope.BoolType, "true")
}

func (sc *symbolConstruction) registerBuiltInConstant(typeSymbol *symbol.TypeSymbol, name string) *symbol.ConstantSymbol {
	constant := symbol.NewConstantSymbol(sc.globalScope, name, typeSymbol)
	sc.globalScope.Constants = append(sc.globalScope.Constants, constant)
	return constant
}

func (sc *symbolConstruction) registerContract() {
	contractNode := sc.programNode.Contract

	contractSymbol := symbol.NewContractSymbol(sc.globalScope, contractNode.Name)
	sc.globalScope.Contract = contractSymbol

	sc.symTable.MapSymbolToNode(contractSymbol, contractNode)
	//for _, variableNode := range contractNode.Variables {
	//	sc.registerField(contractSymbol, variableNode)
	//}
	//
	//for _, functionNode := range contractNode.Functions {
	//	sc.registerFunction(contractSymbol, functionNode)
	//}
}

//func (sc *symbolConstruction) registerField(contractSymbol *symbol.ContractSymbol, node *node.VariableNode) {
//	sym := symbol.FieldSymbol{}.NewSymbol(contractSymbol, node.Identifier)
//	fieldSymbol, _ := sym.(*symbol.FieldSymbol)
//	contractSymbol.Fields = append(contractSymbol.Fields, fieldSymbol)
//	sc.symTable.MapSymbolToNode(node, fieldSymbol)
//}
//
//func (sc *symbolConstruction) registerFunction(contractSymbol *symbol.ContractSymbol, node *node.FunctionNode) {
//	sym := symbol.FunctionSymbol{}.NewSymbol(contractSymbol, node.Name)
//	functionSymbol, _ := sym.(*symbol.FunctionSymbol)
//	contractSymbol.Functions = append(contractSymbol.Functions, functionSymbol)
//	sc.symTable.MapSymbolToNode(node, functionSymbol)
//	for _, parameter := range node.Parameters {
//		sc.registerParameter(functionSymbol, parameter)
//	}
//	for _, statement := range node.Body {
//		// TODO Pass visitor
//		statement.Accept(nil)
//	}
//}
//
//func (sc *symbolConstruction) registerParameter(functionSymbol *symbol.FunctionSymbol, node *node.VariableNode) {
//	sym := symbol.ParameterSymbol{}.NewSymbol(functionSymbol, node.Identifier)
//	parameterSymbol, _ := sym.(*symbol.ParameterSymbol)
//	functionSymbol.Parameters = append(functionSymbol.Parameters, parameterSymbol)
//	sc.symTable.MapSymbolToNode(node, parameterSymbol)
//}

func (sc *symbolConstruction) registerBuiltinFunctions(returnType *symbol.TypeSymbol, identifier string, paramType *symbol.TypeSymbol) {
	// TODO Implement
}

func (sc *symbolConstruction) checkValidIdentifiers() {
	// TODO Implement
}

func (sc *symbolConstruction) checkUniqueIdentifiers() {
	// TODO Implement
}

func (sc *symbolConstruction) checkUniqueIdentifiersByScope() {
	// TODO Implement
}
