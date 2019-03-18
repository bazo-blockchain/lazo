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
	syntaxTree  *node.ProgramNode
	globalScope *symbol.GlobalScope
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
		t:          t,
		syntaxTree: program,
	}
	tester.symbolTable, tester.errors = New(program).Run()
	tester.globalScope = tester.symbolTable.GlobalScope
	assert.Equal(t, len(tester.errors) == 0, isValidCode)

	return tester
}

// Assert Functions
// ----------------

func (ct *CheckerTestUtil) assertTotalErrors(total int) {
	assert.Equal(ct.t, len(ct.errors), total)
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

func (ct *CheckerTestUtil) assertField(index int, expectedType *symbol.TypeSymbol) {
	fieldSymbol := ct.symbolTable.GlobalScope.Contract.Fields[index]
	assert.Equal(ct.t, fieldSymbol.GetScope(), ct.symbolTable.GlobalScope.Contract)
	assert.Equal(ct.t, fieldSymbol.Type, expectedType)
	assert.Equal(ct.t, len(fieldSymbol.AllDeclarations()), 0)

	fieldNode, ok := ct.symbolTable.GetNodeBySymbol(fieldSymbol).(*node.VariableNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, fieldSymbol.GetIdentifier(), fieldNode.Identifier)

	if fieldNode.Expression != nil {
		ct.assertExpressionType(fieldNode.Expression, expectedType)
	}
}

func (ct *CheckerTestUtil) assertFunction(index int, totalReturnTypes int, totalParams int, totalVars int) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[index]
	assert.Equal(ct.t, functionSymbol.GetScope(), ct.symbolTable.GlobalScope.Contract)
	assert.Equal(ct.t, len(functionSymbol.ReturnTypes), totalReturnTypes)
	assert.Equal(ct.t, len(functionSymbol.Parameters), totalParams)
	assert.Equal(ct.t, len(functionSymbol.LocalVariables), totalVars)
	assert.Equal(ct.t, len(functionSymbol.AllDeclarations()), totalParams+totalVars)

	functionNode, ok := ct.symbolTable.GetNodeBySymbol(functionSymbol).(*node.FunctionNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, functionSymbol.GetIdentifier(), functionNode.Name)
}

func (ct *CheckerTestUtil) assertReturnType(funcIndex int, returnTypeIndex int, expectedType *symbol.TypeSymbol) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[funcIndex]

	returnTypeSymbol := functionSymbol.ReturnTypes[returnTypeIndex]
	assert.Equal(ct.t, returnTypeSymbol, expectedType)
}

func (ct *CheckerTestUtil) assertFuncParam(funcIndex int, paramIndex int, expectedType *symbol.TypeSymbol) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[funcIndex]

	paramSymbol := functionSymbol.Parameters[paramIndex]
	assert.Equal(ct.t, paramSymbol.GetScope(), functionSymbol)
	assert.Equal(ct.t, paramSymbol.Type, expectedType)
	assert.Equal(ct.t, len(paramSymbol.AllDeclarations()), 0)

	varNode, ok := ct.symbolTable.GetNodeBySymbol(paramSymbol).(*node.VariableNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, paramSymbol.GetIdentifier(), varNode.Identifier)
}

func (ct *CheckerTestUtil) assertLocalVariable(funcIndex int, varIndex int,
	expectedType *symbol.TypeSymbol, totalVisibleIn int) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[funcIndex]

	varSymbol := functionSymbol.LocalVariables[varIndex]
	assert.Equal(ct.t, varSymbol.GetScope(), functionSymbol)
	assert.Equal(ct.t, varSymbol.Type, expectedType)
	assert.Equal(ct.t, len(varSymbol.VisibleIn), totalVisibleIn)
	assert.Equal(ct.t, len(varSymbol.AllDeclarations()), 0)

	varNode, ok := ct.symbolTable.GetNodeBySymbol(varSymbol).(*node.VariableNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, varSymbol.GetIdentifier(), varNode.Identifier)

	if varNode.Expression != nil {
		ct.assertExpressionType(varNode.Expression, expectedType)
	}
}

func (ct *CheckerTestUtil) assertAssignment(assignStmt *node.AssignmentStatementNode, expectedType *symbol.TypeSymbol) {
	ct.assertExpressionType(assignStmt.Left, expectedType)
	ct.assertExpressionType(assignStmt.Right, expectedType)
}

func (ct *CheckerTestUtil) assertDesignator(expr node.ExpressionNode, decl symbol.Symbol, expectedType *symbol.TypeSymbol) {
	designator, ok := expr.(*node.DesignatorNode)
	assert.Assert(ct.t, ok)

	assert.Equal(ct.t, ct.symbolTable.GetDeclByDesignator(designator), decl)
	ct.assertExpressionType(expr, expectedType)
}

func (ct *CheckerTestUtil) assertExpressionType(expr node.ExpressionNode, expectedType *symbol.TypeSymbol) {
	assert.Equal(ct.t, ct.symbolTable.GetTypeByExpression(expr), expectedType)
}

// Helper Functions
// ----------------

func (ct *CheckerTestUtil) getFieldNode(index int) *node.VariableNode {
	return ct.syntaxTree.Contract.Variables[index]
}

func (ct *CheckerTestUtil) getFuncStatementNode(funcIndex int, stmtIndex int) node.StatementNode{
	return ct.syntaxTree.Contract.Functions[funcIndex].Body[stmtIndex]
}

func (ct *CheckerTestUtil) getLocalVariableSymbol(funcIndex int, varIndex int) *symbol.LocalVariableSymbol {
	return ct.globalScope.Contract.Functions[funcIndex].LocalVariables[varIndex]
}


