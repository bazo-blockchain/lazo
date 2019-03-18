package phase

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/checker/visitor"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
)

type symbolConstruction struct {
	programNode *node.ProgramNode
	symbolTable *symbol.SymbolTable
	globalScope *symbol.GlobalScope
	errors      []error
}

func RunSymbolConstruction(programNode *node.ProgramNode) (*symbol.SymbolTable, []error) {
	symTable := symbol.NewSymbolTable()
	construction := symbolConstruction{
		symbolTable: symTable,
		programNode: programNode,
		globalScope: symTable.GlobalScope,
	}

	if programNode.Contract == nil {
		construction.reportError(nil, "Program has no contract")
		return symTable, construction.errors
	}

	construction.registerBuiltins()
	construction.registerDeclarations()
	construction.checkValidIdentifiers()
	construction.checkUniqueIdentifiers()

	return symTable, construction.errors
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

func (sc *symbolConstruction) registerDeclarations() {
	sc.registerContract()
	// interfaces
}

func (sc *symbolConstruction) registerContract() {
	contractNode := sc.programNode.Contract

	contractSymbol := symbol.NewContractSymbol(sc.globalScope, contractNode.Name)
	sc.globalScope.Contract = contractSymbol
	sc.symbolTable.MapSymbolToNode(contractSymbol, contractNode)

	for _, variableNode := range contractNode.Variables {
		sc.registerField(contractSymbol, variableNode)
	}

	for _, functionNode := range contractNode.Functions {
		sc.registerFunction(contractSymbol, functionNode)
	}
}

func (sc *symbolConstruction) registerField(contractSymbol *symbol.ContractSymbol, node *node.VariableNode) {
	fieldSymbol := symbol.NewFieldSymbol(contractSymbol, node.Identifier)
	contractSymbol.Fields = append(contractSymbol.Fields, fieldSymbol)
	sc.symbolTable.MapSymbolToNode(fieldSymbol, node)
}

func (sc *symbolConstruction) registerFunction(contractSymbol *symbol.ContractSymbol, node *node.FunctionNode) {
	functionSymbol := symbol.NewFunctionSymbol(contractSymbol, node.Name)
	contractSymbol.Functions = append(contractSymbol.Functions, functionSymbol)
	sc.symbolTable.MapSymbolToNode(functionSymbol, node)

	for _, parameter := range node.Parameters {
		sc.registerParameter(functionSymbol, parameter)
	}

	v := visitor.NewLocalVariableVisitor(sc.symbolTable, functionSymbol)
	v.VisitStatementBlock(node.Body)
}

func (sc *symbolConstruction) registerParameter(functionSymbol *symbol.FunctionSymbol, node *node.VariableNode) {
	parameterSymbol := symbol.NewParameterSymbol(functionSymbol, node.Identifier)
	functionSymbol.Parameters = append(functionSymbol.Parameters, parameterSymbol)
	sc.symbolTable.MapSymbolToNode(parameterSymbol, node)
}

func (sc *symbolConstruction) checkValidIdentifiers() {
	sc.checkValidIdentifier(sc.globalScope.Contract)
	for _, field := range sc.globalScope.Contract.Fields {
		sc.checkValidIdentifier(field)
	}
	for _, function := range sc.globalScope.Contract.Functions {
		sc.checkValidIdentifier(function)
		for _, decl := range function.AllDeclarations() {
			sc.checkValidIdentifier(decl)
		}
	}
}

var reservedKeywords = []string{"char", "int", "bool", "string", "this", "null", "void"}

func (sc *symbolConstruction) checkValidIdentifier(sym symbol.Symbol) {
	for _, keyword := range reservedKeywords {
		if sym.GetIdentifier() == keyword {
			sc.reportError(sym, fmt.Sprintf("Reserved keyword '%s' cannot be used as an identifier", keyword))
		}
	}
}

func (sc *symbolConstruction) checkUniqueIdentifiers() {
	sc.checkUniqueIdentifier(sc.globalScope)
	sc.checkUniqueIdentifier(sc.globalScope.Contract)
	for _, function := range sc.globalScope.Contract.Functions {
		sc.checkUniqueIdentifier(function)
	}
}

func (sc *symbolConstruction) checkUniqueIdentifier(sym symbol.Symbol) {
	allDecl := sym.AllDeclarations()
	for r, decl := range allDecl {
		for c, otherDecl := range allDecl {
			if c > r && decl.GetIdentifier() == otherDecl.GetIdentifier() {
				sc.reportError(otherDecl,
					fmt.Sprintf("Identifier '%s' is already declared", otherDecl.GetIdentifier()))
				break
			}
		}
	}
}

func (sc *symbolConstruction) reportError(sym symbol.Symbol, msg string) {
	var pos string
	if sym != nil {
		pos = sc.symbolTable.GetNodeBySymbol(sym).Pos().String()
	}
	sc.errors = append(sc.errors, errors.New(fmt.Sprintf("[%s] %s", pos, msg)))
}
