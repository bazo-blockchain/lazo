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
	assert.Equal(t, len(tester.errors) == 0, isValidCode, tester.errors)

	return tester
}

// Assert Functions
// ----------------

func (ct *CheckerTestUtil) assertTotalErrors(total int) {
	assert.Equal(ct.t, len(ct.errors), total)
}

func (ct *CheckerTestUtil) assertErrorAt(index int, errSubStr string) {
	assert.Assert(ct.t, len(ct.errors) > index)
	err := ct.errors[index].Error()
	assert.Assert(ct.t, strings.Contains(err, errSubStr), err)
}

func (ct *CheckerTestUtil) assertContract(totalVars int, totalFunctions int) {
	contractSymbol := ct.symbolTable.GlobalScope.Contract
	assert.Equal(ct.t, contractSymbol.Scope(), ct.symbolTable.GlobalScope)
	assert.Equal(ct.t, len(contractSymbol.Fields), totalVars)
	assert.Equal(ct.t, len(contractSymbol.Functions), totalFunctions)
	assert.Equal(ct.t, len(contractSymbol.AllDeclarations()), totalVars+totalFunctions)

	contractNode, ok := ct.symbolTable.GetNodeBySymbol(contractSymbol).(*node.ContractNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, contractSymbol.Identifier(), contractNode.Name)
}

func (ct *CheckerTestUtil) assertField(index int, expectedType symbol.TypeSymbol) {
	fieldSymbol := ct.symbolTable.GlobalScope.Contract.Fields[index]
	assert.Equal(ct.t, fieldSymbol.Scope(), ct.symbolTable.GlobalScope.Contract)
	assert.Equal(ct.t, fieldSymbol.Type, expectedType)
	assert.Equal(ct.t, len(fieldSymbol.AllDeclarations()), 0)

	fieldNode, ok := ct.symbolTable.GetNodeBySymbol(fieldSymbol).(*node.FieldNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, fieldSymbol.Identifier(), fieldNode.Identifier)

	if fieldNode.Expression != nil {
		ct.assertExpressionType(fieldNode.Expression, expectedType)
	}
}

func (ct *CheckerTestUtil) assertStruct(structName string, totalFields int) {
	structType := ct.globalScope.Structs[structName]
	assert.Equal(ct.t, structType.Scope(), ct.globalScope.Contract)
	assert.Equal(ct.t, len(structType.Fields), 2)
	assert.Equal(ct.t, len(structType.AllDeclarations()), totalFields)

	structNode, ok := ct.symbolTable.GetNodeBySymbol(structType).(*node.StructNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, structType.Identifier(), structNode.Name)
}

func (ct *CheckerTestUtil) assertStructField(structName string, fieldIndex int, expectedType symbol.TypeSymbol) {
	structType := ct.symbolTable.GlobalScope.Structs[structName]
	fieldSymbol := structType.Fields[fieldIndex]
	assert.Equal(ct.t, fieldSymbol.Scope(), structType)
	assert.Equal(ct.t, fieldSymbol.Type, expectedType)
	assert.Equal(ct.t, len(fieldSymbol.AllDeclarations()), 0)

	structFieldNode, ok := ct.symbolTable.GetNodeBySymbol(fieldSymbol).(*node.StructFieldNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, fieldSymbol.Identifier(), structFieldNode.Identifier)
}

func (ct *CheckerTestUtil) assertConstructor(totalParams int, totalVars int) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Constructor
	assert.Equal(ct.t, functionSymbol.Scope(), ct.symbolTable.GlobalScope.Contract)
	assert.Equal(ct.t, len(functionSymbol.ReturnTypes), 0)
	assert.Equal(ct.t, len(functionSymbol.Parameters), totalParams)
	assert.Equal(ct.t, len(functionSymbol.LocalVariables), totalVars)
	assert.Equal(ct.t, len(functionSymbol.AllDeclarations()), totalParams+totalVars)

	_, ok := ct.symbolTable.GetNodeBySymbol(functionSymbol).(*node.ConstructorNode)
	assert.Assert(ct.t, ok)
}

func (ct *CheckerTestUtil) assertFunction(index int, totalReturnTypes int, totalParams int, totalVars int) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[index]
	assert.Equal(ct.t, functionSymbol.Scope(), ct.symbolTable.GlobalScope.Contract)
	assert.Equal(ct.t, len(functionSymbol.ReturnTypes), totalReturnTypes)
	assert.Equal(ct.t, len(functionSymbol.Parameters), totalParams)
	assert.Equal(ct.t, len(functionSymbol.LocalVariables), totalVars)
	assert.Equal(ct.t, len(functionSymbol.AllDeclarations()), totalParams+totalVars)

	functionNode, ok := ct.symbolTable.GetNodeBySymbol(functionSymbol).(*node.FunctionNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, functionSymbol.Identifier(), functionNode.Name)
}

func (ct *CheckerTestUtil) assertReturnType(funcIndex int, returnTypeIndex int, expectedType symbol.TypeSymbol) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[funcIndex]

	returnTypeSymbol := functionSymbol.ReturnTypes[returnTypeIndex]
	assert.Equal(ct.t, returnTypeSymbol, expectedType)
}

func (ct *CheckerTestUtil) assertFuncParam(funcIndex int, paramIndex int, expectedType symbol.TypeSymbol) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[funcIndex]

	paramSymbol := functionSymbol.Parameters[paramIndex]
	ct.assertParam(paramSymbol, functionSymbol, expectedType)
}

func (ct *CheckerTestUtil) assertParam(paramSymbol *symbol.ParameterSymbol, scope symbol.Symbol,
	expectedType symbol.TypeSymbol) {
	assert.Equal(ct.t, paramSymbol.Scope(), scope)
	assert.Equal(ct.t, paramSymbol.Type, expectedType)
	assert.Equal(ct.t, len(paramSymbol.AllDeclarations()), 0)

	varNode, ok := ct.symbolTable.GetNodeBySymbol(paramSymbol).(*node.ParameterNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, paramSymbol.Identifier(), varNode.Identifier)
}

func (ct *CheckerTestUtil) assertFuncLocalVariable(funcIndex int, varIndex int,
	expectedType symbol.TypeSymbol, totalVisibleIn int) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[funcIndex]

	varSymbol := functionSymbol.LocalVariables[varIndex]
	ct.assertLocalVariable(varSymbol, functionSymbol, expectedType, totalVisibleIn)
}

func (ct *CheckerTestUtil) assertLocalVariable(varSymbol *symbol.LocalVariableSymbol, scope symbol.Symbol,
	expectedType symbol.TypeSymbol, totalVisibleIn int) {
	assert.Equal(ct.t, varSymbol.Scope(), scope)
	assert.Equal(ct.t, varSymbol.Type, expectedType)
	assert.Equal(ct.t, len(varSymbol.VisibleIn), totalVisibleIn)
	assert.Equal(ct.t, len(varSymbol.AllDeclarations()), 0)

	varNode, ok := ct.symbolTable.GetNodeBySymbol(varSymbol).(*node.VariableNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, varSymbol.Identifier(), varNode.Identifier)

	if varNode.Expression != nil {
		ct.assertExpressionType(varNode.Expression, expectedType)
	}
}

func (ct *CheckerTestUtil) assertMultiLocalVariable(funcIndex int, varIndex int, multiVarIndex int,
	expectedType symbol.TypeSymbol, totalVisibleIn int) {
	functionSymbol := ct.symbolTable.GlobalScope.Contract.Functions[funcIndex]

	varSymbol := functionSymbol.LocalVariables[varIndex]
	assert.Equal(ct.t, varSymbol.Scope(), functionSymbol)
	assert.Equal(ct.t, varSymbol.Type, expectedType)
	assert.Equal(ct.t, len(varSymbol.VisibleIn), totalVisibleIn)
	assert.Equal(ct.t, len(varSymbol.AllDeclarations()), 0)

	varNode, ok := ct.symbolTable.GetNodeBySymbol(varSymbol).(*node.MultiVariableNode)
	assert.Assert(ct.t, ok)
	assert.Equal(ct.t, varSymbol.Identifier(), varNode.Identifiers[multiVarIndex])

	calledFuncSym := ct.symbolTable.GetDeclByDesignator(varNode.FuncCall.Designator).(*symbol.FunctionSymbol)
	assert.Equal(ct.t, calledFuncSym.ReturnTypes[multiVarIndex], expectedType)
}

func (ct *CheckerTestUtil) assertAssignment(assignStmt *node.AssignmentStatementNode, expectedType symbol.TypeSymbol) {
	ct.assertExpressionType(assignStmt.Left, expectedType)
	ct.assertExpressionType(assignStmt.Right, expectedType)
}

func (ct *CheckerTestUtil) assertBasicDesignator(expr node.ExpressionNode, decl symbol.Symbol, expectedType symbol.TypeSymbol) {
	designator, ok := expr.(*node.BasicDesignatorNode)
	assert.Assert(ct.t, ok)
	ct.assertDesignator(designator, decl, expectedType)
}

func (ct *CheckerTestUtil) assertElementAccess(expr node.ExpressionNode, decl symbol.Symbol, expectedType symbol.TypeSymbol) {
	designator, ok := expr.(*node.ElementAccessNode)
	assert.Assert(ct.t, ok)
	ct.assertDesignator(designator, decl, expectedType)
}

func (ct *CheckerTestUtil) assertMemberAccess(expr node.ExpressionNode, decl symbol.Symbol, expectedType symbol.TypeSymbol) {
	designator, ok := expr.(*node.MemberAccessNode)
	assert.Assert(ct.t, ok)
	ct.assertDesignator(designator, decl, expectedType)
}

func (ct *CheckerTestUtil) assertDesignator(designator node.DesignatorNode, decl symbol.Symbol, expectedType symbol.TypeSymbol) {
	assert.Equal(ct.t, ct.symbolTable.GetDeclByDesignator(designator), decl)
	ct.assertExpressionType(designator, expectedType)
}

func (ct *CheckerTestUtil) assertExpressionType(expr node.ExpressionNode, expectedType symbol.TypeSymbol) {
	assert.Equal(ct.t, ct.symbolTable.GetTypeByExpression(expr), expectedType)
}

func (ct *CheckerTestUtil) assertArrayLengthCreation(expr node.ExpressionNode, expectedType symbol.TypeSymbol) {
	creation, ok := expr.(*node.ArrayLengthCreationNode)
	assert.Assert(ct.t, ok)
	ct.assertExpressionType(creation, expectedType)
	for _, length := range creation.Lengths {
		ct.assertExpressionType(length, ct.globalScope.IntType)
	}
}

func (ct *CheckerTestUtil) assertArrayValueCreation(expr node.ExpressionNode, expectedType symbol.TypeSymbol, expectedValueType symbol.TypeSymbol) {
	creation, ok := expr.(*node.ArrayValueCreationNode)
	assert.Assert(ct.t, ok)
	ct.assertExpressionType(creation, expectedType)
	for _, element := range creation.Elements.Values {
		ct.assertExpressionType(element, expectedValueType)
	}
}

// Helper Functions
// ----------------

func (ct *CheckerTestUtil) getFieldNode(index int) *node.FieldNode {
	return ct.syntaxTree.Contract.Fields[index]
}

func (ct *CheckerTestUtil) getConstructorStatementNode(stmtIndex int) node.StatementNode {
	return ct.syntaxTree.Contract.Constructor.Body[stmtIndex]
}

func (ct *CheckerTestUtil) getFuncStatementNode(funcIndex int, stmtIndex int) node.StatementNode {
	return ct.syntaxTree.Contract.Functions[funcIndex].Body[stmtIndex]
}

func (ct *CheckerTestUtil) getLocalVariableSymbol(funcIndex int, varIndex int) *symbol.LocalVariableSymbol {
	return ct.globalScope.Contract.Functions[funcIndex].LocalVariables[varIndex]
}
