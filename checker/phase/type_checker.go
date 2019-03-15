package phase

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/checker/visitor"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type TypeChecker struct {
	symTable *symbol.SymbolTable
	errors []error
}


func RunTypeChecker(symTable *symbol.SymbolTable) []error {
	check :=TypeChecker{
		symTable: symTable,
	}
	check.checkTypes()
	return check.errors
}

func (tc *TypeChecker) checkTypes() {
	contractSymbol := tc.symTable.GlobalScope.Contract
	v := visitor.NewTypeCheckVisitor(tc.symTable, contractSymbol)
	contractNode := tc.symTable.GetNodeBySymbol(contractSymbol).(*node.ContractNode)
	contractNode.Accept(v)
}