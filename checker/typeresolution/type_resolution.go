package typeresolution

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type typeResolution struct {
	symTable *symbol.SymbolTable
	errors   []error
}

// Run performs type resolution
// Returns errors that occurred during type resolution
func Run(symTable *symbol.SymbolTable) []error {
	tr := typeResolution{
		symTable: symTable,
	}
	tr.resolveTypesInContractSymbol()
	tr.resolveTypesInStruct()
	return tr.errors
}

func (tr *typeResolution) resolveTypesInContractSymbol() {
	contractSymbol := tr.symTable.GlobalScope.Contract
	for _, field := range contractSymbol.Fields {
		tr.resolveTypeInFieldSymbol(field)
	}

	if contractSymbol.Constructor != nil {
		tr.resolveTypeInFunctionSymbol(contractSymbol.Constructor)
	}

	for _, function := range contractSymbol.Functions {
		tr.resolveTypeInFunctionSymbol(function)
	}
}

func (tr *typeResolution) resolveTypesInStruct() {
	for _, structType := range tr.symTable.GlobalScope.Structs {
		for _, fieldSym := range structType.Fields {
			fieldNode := tr.symTable.GetNodeBySymbol(fieldSym).(*node.StructFieldNode)
			fieldSym.Type = tr.resolveType(fieldNode.Type)
		}
	}
}

func (tr *typeResolution) resolveTypeInFieldSymbol(symbol *symbol.FieldSymbol) {
	fieldNode := tr.symTable.GetNodeBySymbol(symbol).(*node.FieldNode)
	symbol.Type = tr.resolveType(fieldNode.Type)
}

func (tr *typeResolution) resolveTypeInFunctionSymbol(sym *symbol.FunctionSymbol) {
	if functionNode, ok := tr.symTable.GetNodeBySymbol(sym).(*node.FunctionNode); ok {
		tr.resolveReturnTypes(sym, functionNode)
	}

	for _, param := range sym.Parameters {
		paramNode := tr.symTable.GetNodeBySymbol(param).(*node.ParameterNode)
		param.Type = tr.resolveType(paramNode.Type)
	}

	for _, locSym := range sym.LocalVariables {
		locVarNode := tr.symTable.GetNodeBySymbol(locSym)

		if varNode, ok := locVarNode.(*node.VariableNode); ok {
			locSym.Type = tr.resolveType(varNode.Type)
		} else if multiVarNode, ok := locVarNode.(*node.MultiVariableNode); ok {
			locSym.Type = tr.resolveType(multiVarNode.GetType(locSym.ID))
		} else {
			tr.reportError(locVarNode, fmt.Sprintf("Unsupported local variable node type"))
		}
	}
}

func (tr *typeResolution) resolveReturnTypes(sym *symbol.FunctionSymbol, functionNode *node.FunctionNode) {
	total := len(functionNode.ReturnTypes)
	if total > 3 {
		tr.reportError(functionNode, "More than 3 return types are not allowed")
	}

	for _, rtype := range functionNode.ReturnTypes {
		if rtype.String() == "void" {
			if total > 1 {
				tr.reportError(rtype, "'void' is invalid with multiple return types")
			}
		} else {
			sym.ReturnTypes = append(sym.ReturnTypes, tr.resolveType(rtype))
		}
	}
}

func (tr *typeResolution) resolveType(node node.TypeNode) symbol.TypeSymbol {
	result := tr.symTable.FindTypeByNode(node)
	if result == nil {
		tr.reportError(node, fmt.Sprintf("Invalid type '%s'", node.String()))
	}
	return result
}

func (tr *typeResolution) reportError(node node.Node, msg string) {
	tr.errors = append(tr.errors, fmt.Errorf("[%s] %s", node.Pos(), msg))
}
