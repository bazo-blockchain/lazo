package checker

import (
	"bufio"
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser"
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"strings"
	"testing"
)

type CheckerTestUtil struct {
	t           *testing.T
	symbolTable *symbol.SymbolTable
	errors      []error
}

func newCheckerTestUtil(t *testing.T, contractCode string, isValidCode bool) *CheckerTestUtil {
	return newCheckerTestUtilWithRawInput(
		t,
		fmt.Sprintf("contract Test {\n %s \n }", contractCode),
		isValidCode,
	)
}

func newCheckerTestUtilWithRawInput(t *testing.T, code string, isValidCode bool) *CheckerTestUtil {
	p := parser.New(lexer.New(bufio.NewReader(strings.NewReader(code))))
	program, err := p.ParseProgram()
	assert.Equal(t, len(err), 0, "Program has syntax errors", err)

	tester := &CheckerTestUtil{
		t: t,
	}
	tester.symbolTable, tester.errors = New(program).Run()
	assert.Equal(t, len(tester.errors) == 0, isValidCode)

	return tester
}

func (ct *CheckerTestUtil) assertContract(totalVars int, totalFunctions int) {
	contractSymbol := ct.symbolTable.GlobalScope.Contract
	assert.Equal(ct.t, contractSymbol.GetScope(), ct.symbolTable.GlobalScope)
	assert.Equal(ct.t, len(contractSymbol.Fields), totalVars)
	assert.Equal(ct.t, len(contractSymbol.Functions), totalFunctions)
	assert.Equal(ct.t, len(contractSymbol.AllDeclarations()), totalVars+totalFunctions)

	contractNode, ok := ct.symbolTable.GetNodeBySymbol(contractSymbol).(*node.ContractNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, contractSymbol.GetIdentifier(), contractNode.Name)
}

func (ct *CheckerTestUtil) assertField(index int, expectedType *symbol.TypeSymbol){
	fieldSymbol := ct.symbolTable.GlobalScope.Contract.Fields[index]
	assert.Equal(ct.t, fieldSymbol.GetScope(), ct.symbolTable.GlobalScope.Contract)
	assert.Equal(ct.t, fieldSymbol.Type, expectedType)
	assert.Equal(ct.t, len(fieldSymbol.AllDeclarations()), 0)

	fieldNode, ok := ct.symbolTable.GetNodeBySymbol(fieldSymbol).(*node.VariableNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, fieldSymbol.GetIdentifier(), fieldNode.Identifier)
}

func (ct *CheckerTestUtil) assertFunction(index int, totalReturnTypes int, totalParams int, totalVars int) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[0]
	assert.Equal(ct.t, functionSymbol.GetScope(), ct.symbolTable.GlobalScope.Contract)
	assert.Equal(ct.t, len(functionSymbol.ReturnTypes), totalReturnTypes)
	assert.Equal(ct.t, len(functionSymbol.Parameters), totalParams)
	assert.Equal(ct.t, len(functionSymbol.LocalVariables), totalVars)
	assert.Equal(ct.t, len(functionSymbol.AllDeclarations()), totalParams + totalVars)

	functionNode, ok := ct.symbolTable.GetNodeBySymbol(functionSymbol).(*node.FunctionNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, functionSymbol.GetIdentifier(), functionNode.Name)
}
