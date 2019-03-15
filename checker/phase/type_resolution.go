package phase

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
)

type TypeResolution struct {
	symTable *symbol.SymbolTable
	errors []error
}

func RunTypeResolution(symTable *symbol.SymbolTable) []error {
	resolution :=TypeResolution{
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

	if functionNode.HasReturnTypes() {
		sym.ReturnTypes = make([]*symbol.TypeSymbol, len(functionNode.ReturnTypes))
		for i, rtype := range functionNode.ReturnTypes {
			sym.ReturnTypes[i] = tr.resolveType(rtype)
		}
	}

	for _, param := range sym.Parameters {
		paramNode := tr.symTable.GetNodeBySymbol(param).(*node.VariableNode)
		param.Type = tr.resolveType(paramNode.Type)
	}

	for _, locVar := range sym.LocalVariables {
		locVarNode := tr.symTable.GetNodeBySymbol(locVar).(*node.VariableNode)
		locVar.Type = tr.resolveType(locVarNode.Type)
	}
}

func (tr *TypeResolution) resolveType(node *node.TypeNode) *symbol.TypeSymbol {
	result := tr.symTable.FindTypeByNode(node)
	if result == nil {

		fmt.Printf("Error: Could not find type for node %s", node.Identifier)
	}
	return result
}

func (tr *TypeResolution) reportTypeResolutionError(sym symbol.Symbol, msg string) {
	tr.errors = append(tr.errors, errors.New(fmt.Sprintf("[%s] %s", tr.symTable.GetNodeBySymbol(sym).Pos(), msg)))
}