package typecheck

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type typeChecker struct {
	symTable *symbol.SymbolTable
	errors   []error
}

// Run performs type checks
// Returns errors that occurred during type checking
func Run(symTable *symbol.SymbolTable) []error {
	check := typeChecker{
		symTable: symTable,
	}
	check.checkTypes()
	return check.errors
}

func (tc *typeChecker) checkTypes() {
	contractSymbol := tc.symTable.GlobalScope.Contract
	v := newTypeCheckVisitor(tc.symTable, contractSymbol)
	contractNode := tc.symTable.GetNodeBySymbol(contractSymbol).(*node.ContractNode)

	contractNode.Accept(v)
	tc.errors = v.Errors
}
