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

func (tc *TypeResolution) resolveTypesInContractSymbol() {
	contractSymbol := tc.symTable.GlobalScope.Contract
	for _, field := range contractSymbol.Fields {
		tc.resolveTypeInFieldSymbol(field)
	}

	for _, function := range contractSymbol.Functions {
		tc.resolveTypeInFunctionSymbol(function)
	}
}

func (tc *TypeResolution) resolveTypeInFieldSymbol(symbol *symbol.FieldSymbol) {
	fieldNode := tc.symTable.GetNodeBySymbol(symbol).(*node.VariableNode)
	symbol.Type = tc.resolveType(fieldNode.Type)
}

func (tc *TypeResolution) resolveTypeInFunctionSymbol(sym *symbol.FunctionSymbol) {
	functionNode := tc.symTable.GetNodeBySymbol(sym).(*node.FunctionNode)

	if functionNode.HasReturnTypes() {
		sym.ReturnTypes = make([]*symbol.TypeSymbol, len(functionNode.ReturnTypes))
		for i, rtype := range functionNode.ReturnTypes {
			sym.ReturnTypes[i] = tc.resolveType(rtype)
		}
	}

	for _, param := range sym.Parameters {
		paramNode := tc.symTable.GetNodeBySymbol(param).(*node.VariableNode)
		param.Type = tc.resolveType(paramNode.Type)
	}

	for _, locVar := range sym.LocalVariables {
		locVarNode := tc.symTable.GetNodeBySymbol(locVar).(*node.VariableNode)
		locVar.Type = tc.resolveType(locVarNode.Type)
	}
}

func (tc *TypeResolution) resolveType(node *node.TypeNode) *symbol.TypeSymbol {
	result := tc.symTable.FindTypeByNode(node)
	if result == nil {
		tc.reportError(node, fmt.Sprintf("Invalid type '%s'", node.Identifier))
	}
	return result
}

func (tc *TypeResolution) reportError(node node.Node, msg string) {
	tc.errors = append(tc.errors, errors.New(
		fmt.Sprintf("[%s] %s", node.Pos(), msg)))
}
