package symbolconstruction

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type symbolConstruction struct {
	programNode *node.ProgramNode
	symbolTable *symbol.SymbolTable
	globalScope *symbol.GlobalScope
	errors      []error
}

// Run prepares global scope, creates symbols and checks identifiers
// Returns errors that occurred during construction
func Run(programNode *node.ProgramNode) (*symbol.SymbolTable, []error) {
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
	sc.registerBuiltinField()
}

func (sc *symbolConstruction) registerBuiltInTypes() {
	sc.globalScope.BoolType = sc.registerBuiltInType("bool")
	sc.globalScope.CharType = sc.registerBuiltInType("char")
	sc.globalScope.IntType = sc.registerBuiltInType("int")
	sc.globalScope.StringType = sc.registerBuiltInType("String")
}

func (sc *symbolConstruction) registerBuiltInType(name string) *symbol.BasicTypeSymbol {
	baseType := symbol.NewBasicTypeSymbol(sc.globalScope, name)
	sc.globalScope.Types[name] = baseType
	sc.globalScope.BuiltInTypes = append(sc.globalScope.BuiltInTypes, baseType)
	return baseType
}

func (sc *symbolConstruction) registerBuiltInConstants() {
	sc.globalScope.NullConstant = sc.registerBuiltInConstant(sc.globalScope.NullType, "null")
	sc.globalScope.FalseConstant = sc.registerBuiltInConstant(sc.globalScope.BoolType, "false")
	sc.globalScope.TrueConstant = sc.registerBuiltInConstant(sc.globalScope.BoolType, "true")
}

func (sc *symbolConstruction) registerBuiltInConstant(typeSymbol *symbol.BasicTypeSymbol, name string) *symbol.ConstantSymbol {
	constant := symbol.NewConstantSymbol(sc.globalScope, name, typeSymbol)
	sc.globalScope.Constants = append(sc.globalScope.Constants, constant)
	return constant
}

func (sc *symbolConstruction) registerBuiltinField() {
	arrayLength := &symbol.FieldSymbol{
		AbstractSymbol: symbol.NewAbstractSymbol(sc.globalScope, "length"),
		Type:           sc.symbolTable.FindTypeByIdentifier("int"),
	}
	sc.globalScope.ArrayLength = arrayLength
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

	for _, fieldNode := range contractNode.Fields {
		sc.registerField(contractSymbol, fieldNode)
	}

	for _, structNode := range contractNode.Structs {
		sc.registerStruct(contractSymbol, structNode)
	}

	if contractNode.Constructor != nil {
		sc.registerConstructor(contractSymbol, contractNode.Constructor)
	}

	for _, functionNode := range contractNode.Functions {
		sc.registerFunction(contractSymbol, functionNode)
	}
}

func (sc *symbolConstruction) registerField(contractSymbol *symbol.ContractSymbol, node *node.FieldNode) {
	fieldSymbol := symbol.NewFieldSymbol(contractSymbol, node.Identifier)
	contractSymbol.Fields = append(contractSymbol.Fields, fieldSymbol)
	sc.symbolTable.MapSymbolToNode(fieldSymbol, node)
}

func (sc *symbolConstruction) registerStruct(contractSymbol *symbol.ContractSymbol, node *node.StructNode) {
	structType := symbol.NewStructTypeSymbol(contractSymbol, node.Name)
	sc.symbolTable.MapSymbolToNode(structType, node)

	if _, ok := sc.globalScope.Structs[node.Name]; ok {
		sc.reportError(structType,
			fmt.Sprintf("Struct '%s' is already declared", structType.Identifier()))
		return
	}

	sc.globalScope.Structs[node.Name] = structType
	sc.globalScope.Types[node.Name] = structType

	for _, fieldNode := range node.Fields {
		fieldSymbol := symbol.NewFieldSymbol(structType, fieldNode.Identifier)
		structType.Fields = append(structType.Fields, fieldSymbol)
		sc.symbolTable.MapSymbolToNode(fieldSymbol, fieldNode)
	}
}

func (sc *symbolConstruction) registerConstructor(contractSymbol *symbol.ContractSymbol, node *node.ConstructorNode) {
	constructor := symbol.NewFunctionSymbol(contractSymbol, "constructor")
	contractSymbol.Constructor = constructor
	sc.symbolTable.MapSymbolToNode(constructor, node)

	for _, parameter := range node.Parameters {
		sc.registerParameter(constructor, parameter)
	}

	v := newLocalVariableVisitor(sc.symbolTable, constructor)
	v.VisitStatementBlock(node.Body)
}

func (sc *symbolConstruction) registerFunction(contractSymbol *symbol.ContractSymbol, node *node.FunctionNode) {
	functionSymbol := symbol.NewFunctionSymbol(contractSymbol, node.Name)
	contractSymbol.Functions = append(contractSymbol.Functions, functionSymbol)
	sc.symbolTable.MapSymbolToNode(functionSymbol, node)

	for _, parameter := range node.Parameters {
		sc.registerParameter(functionSymbol, parameter)
	}

	v := newLocalVariableVisitor(sc.symbolTable, functionSymbol)
	v.VisitStatementBlock(node.Body)
}

func (sc *symbolConstruction) registerParameter(functionSymbol *symbol.FunctionSymbol, node *node.ParameterNode) {
	parameterSymbol := symbol.NewParameterSymbol(functionSymbol, node.Identifier)
	functionSymbol.Parameters = append(functionSymbol.Parameters, parameterSymbol)
	sc.symbolTable.MapSymbolToNode(parameterSymbol, node)
}

func (sc *symbolConstruction) checkValidIdentifiers() {
	contract := sc.globalScope.Contract
	sc.checkValidIdentifier(contract)
	for _, field := range contract.Fields {
		sc.checkValidIdentifier(field)
	}

	for _, structType := range sc.globalScope.Structs {
		sc.checkValidIdentifier(structType)
		for _, field := range structType.Fields {
			sc.checkValidIdentifier(field)
		}
	}

	if contract.Constructor != nil {
		for _, decl := range contract.Constructor.AllDeclarations() {
			sc.checkValidIdentifier(decl)
		}
	}

	for _, function := range contract.Functions {
		sc.checkValidIdentifier(function)
		for _, decl := range function.AllDeclarations() {
			sc.checkValidIdentifier(decl)
		}
	}
}

var reservedKeywords = []string{"char", "int", "bool", "string", "this", "null", "void"}

func (sc *symbolConstruction) checkValidIdentifier(sym symbol.Symbol) {
	for _, keyword := range reservedKeywords {
		if sym.Identifier() == keyword {
			sc.reportError(sym, fmt.Sprintf("Reserved keyword '%s' cannot be used as an identifier", keyword))
			return
		}
	}
	for _, structType := range sc.globalScope.Structs {
		if sym != structType && sym.Identifier() == structType.Identifier() {
			sc.reportError(sym, fmt.Sprintf("Struct name %s cannot be used as an identifier",
				structType.Identifier()))
			return
		}
	}
}

func (sc *symbolConstruction) checkUniqueIdentifiers() {
	sc.checkUniqueIdentifier(sc.globalScope)
	sc.checkUniqueIdentifier(sc.globalScope.Contract)

	for _, structType := range sc.globalScope.Structs {
		sc.checkUniqueIdentifier(structType)
	}

	if sc.globalScope.Contract.Constructor != nil {
		sc.checkUniqueIdentifier(sc.globalScope.Contract.Constructor)
	}

	for _, function := range sc.globalScope.Contract.Functions {
		sc.checkUniqueIdentifier(function)
	}
}

func (sc *symbolConstruction) checkUniqueIdentifier(sym symbol.Symbol) {
	allDecl := sym.AllDeclarations()
	for r, decl := range allDecl {
		for c, otherDecl := range allDecl {
			if c > r && decl.Identifier() == otherDecl.Identifier() {
				sc.reportError(otherDecl,
					fmt.Sprintf("Identifier '%s' is already declared", otherDecl.Identifier()))
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
	sc.errors = append(sc.errors, fmt.Errorf("[%s] %s", pos, msg))
}
