package phase

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type TypeResolution struct {
	symTable *symbol.SymbolTable
}

func RunTypeResolution(symTable *symbol.SymbolTable) {
	resolution :=TypeResolution{
		symTable: symTable,
	}
	resolution.resolveTypesInContractSymbol()
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

func (tr *TypeResolution) resolveTypeInFunctionSymbol(symbol *symbol.FunctionSymbol) {

}

func (tr *TypeResolution) resolveType(node *node.TypeNode) *symbol.TypeSymbol {
	result := tr.symTable.FindTypeByNode(node)
	if result == nil {
		fmt.Printf("Error: Could not find type for node %s", node.Identifier)
	}
	return result
}