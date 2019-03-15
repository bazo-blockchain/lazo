package phase

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/checker/visitor"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
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

func (tc *TypeResolution) reportTypeCheckerError(sym symbol.Symbol, msg string) {
	tc.errors = append(tc.errors, errors.New(fmt.Sprintf("[%s] %s", tc.symTable.GetNodeBySymbol(sym).Pos(), msg)))
}