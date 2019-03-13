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
	// sc.registerContract(sc.programNode.Contract)
}

func (sc *symbolConstruction) registerBuiltins() {
	sc.registerBuiltinTypes()
	sc.registerBuiltinConstants()
}

//func (sc *symbolConstruction) registerContract(contractNode *node.ContractNode) {
//	sym := symbol.ContractSymbol{}.NewSymbol(sc.globalScope, contractNode.Name)
//
//	contractSymbol, _ := sym.(*symbol.ContractSymbol)
//	sc.globalScope.Contract = contractSymbol
//
//	typeSymbol, _ := sym.(*symbol.TypeSymbol)
//	sc.globalScope.Types = append(sc.globalScope.Types, typeSymbol)
//
//	sc.symTable.LinkDeclaration(contractNode, contractSymbol)
//	for _, variableNode := range contractNode.Variables {
//		sc.registerField(contractSymbol, variableNode)
//	}
//
//	for _, functionNode := range contractNode.Functions {
//		sc.registerFunction(contractSymbol, functionNode)
//	}
//}

//func (sc *symbolConstruction) registerField(contractSymbol *symbol.ContractSymbol, node *node.VariableNode) {
//	sym := symbol.FieldSymbol{}.NewSymbol(contractSymbol, node.Identifier)
//	fieldSymbol, _ := sym.(*symbol.FieldSymbol)
//	contractSymbol.Fields = append(contractSymbol.Fields, fieldSymbol)
//	sc.symTable.LinkDeclaration(node, fieldSymbol)
//}
//
//func (sc *symbolConstruction) registerFunction(contractSymbol *symbol.ContractSymbol, node *node.FunctionNode) {
//	sym := symbol.FunctionSymbol{}.NewSymbol(contractSymbol, node.Name)
//	functionSymbol, _ := sym.(*symbol.FunctionSymbol)
//	contractSymbol.Functions = append(contractSymbol.Functions, functionSymbol)
//	sc.symTable.LinkDeclaration(node, functionSymbol)
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
//	sc.symTable.LinkDeclaration(node, parameterSymbol)
//}

func (sc *symbolConstruction) registerBuiltinTypes() {
	sc.globalScope.BoolType = sc.registerBuiltinType("bool")
	sc.globalScope.CharType = sc.registerBuiltinType("char")
	sc.globalScope.IntType = sc.registerBuiltinType("int")
	sc.globalScope.StringType = sc.registerBuiltinType("string")
}

func (sc *symbolConstruction) registerBuiltinFunctions(returnType *symbol.TypeSymbol, identifier string, paramType *symbol.TypeSymbol) {
	// TODO Implement
}

func (sc *symbolConstruction) registerBuiltinType(name string) *symbol.TypeSymbol {
	baseType := symbol.NewTypeSymbol(sc.globalScope, name)
	sc.globalScope.Types = append(sc.globalScope.Types, baseType)
	return baseType
}

func (sc *symbolConstruction) registerBuiltinConstants() {
	//sc.globalScope.NullConstant = sc.registerBuiltinConstant()
	//sc.globalScope.FalseConstant = sc.registerBuiltinConstant()
	//sc.globalScope.TrueConstant = sc.registerBuiltinConstant()
}

//func (sc *symbolConstruction) registerBuiltinConstant(typeSymbol *symbol.TypeSymbol, name string) *symbol.ConstantSymbol {
//	constant := symbol.ConstantSymbol{}.NewConstantSymbol(sc.globalScope, name, typeSymbol)
//	sc.globalScope.Constants = append(sc.globalScope.Constants, constant)
//	return constant
//}

func (sc *symbolConstruction) checkValidIdentifiers() {
	// TODO Implement
}

func (sc *symbolConstruction) checkUniqueIdentifiers() {
	// TODO Implement
}

func (sc *symbolConstruction) checkUniqueIdentifiersByScope() {
	// TODO Implement
}
