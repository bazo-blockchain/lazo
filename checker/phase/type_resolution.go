package phase

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
)

type TypeResolution struct {
	symTable *symbol.SymbolTable
	errors   []error
}

func RunTypeResolution(symTable *symbol.SymbolTable) []error {
	resolution := TypeResolution{
		symTable: symTable,
	}
	resolution.resolveTypesInContractSymbol()
	return resolution.errors
}

func (tr *TypeResolution) resolveTypesInContractSymbol() {
	contractSymbol := tr.symTable.GlobalScope.Contract
	for _, field := range contractSymbol.Fields {
		tr.resolveTypeInFieldSymbol(field)
	}

	for _, function := range contractSymbol.Functions {
		tr.resolveTypeInFunctionSymbol(function)
	}
}

func (tr *TypeResolution) resolveTypeInFieldSymbol(symbol *symbol.FieldSymbol) {
	fieldNode := tr.symTable.GetNodeBySymbol(symbol).(*node.VariableNode)
	symbol.Type = tr.resolveType(fieldNode.Type)
}

func (tr *TypeResolution) resolveTypeInFunctionSymbol(sym *symbol.FunctionSymbol) {
	functionNode := tr.symTable.GetNodeBySymbol(sym).(*node.FunctionNode)

	tr.resolveReturnTypes(sym, functionNode)

	for _, param := range sym.Parameters {
		paramNode := tr.symTable.GetNodeBySymbol(param).(*node.VariableNode)
		param.Type = tr.resolveType(paramNode.Type)
	}

	for _, locVar := range sym.LocalVariables {
		locVarNode := tr.symTable.GetNodeBySymbol(locVar).(*node.VariableNode)
		locVar.Type = tr.resolveType(locVarNode.Type)
	}
}

func (tr *TypeResolution) resolveReturnTypes(sym *symbol.FunctionSymbol, functionNode *node.FunctionNode) {
	total := len(functionNode.ReturnTypes)
	if  total > 3 {
		tr.reportError(functionNode, "More than 3 return types are not allowed")
	}

	for _, rtype := range functionNode.ReturnTypes {
		if rtype.Identifier == "void" {
			if total > 1 {
				tr.reportError(rtype, "'void' is invalid with multiple return types")
			}
		} else {
			sym.ReturnTypes = append(sym.ReturnTypes, tr.resolveType(rtype))
		}
	}
}

func (tr *TypeResolution) resolveType(node *node.TypeNode) *symbol.TypeSymbol {
	result := tr.symTable.FindTypeByNode(node)
	if result == nil {
		tr.reportError(node, fmt.Sprintf("Invalid type '%s'", node.Identifier))
	}
	return result
}

func (tr *TypeResolution) reportError(node node.Node, msg string) {
	tr.errors = append(tr.errors, errors.New(
		fmt.Sprintf("[%s] %s", node.Pos(), msg)))
}
