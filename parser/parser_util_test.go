package parser

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"math/big"
	"strings"
	"testing"
)

// Helpers
// ------------------------------

func newParserFromInput(input string) *Parser {
	return New(lexer.New(bufio.NewReader(strings.NewReader(input))))
}

func assertHasError(t *testing.T, p *Parser) {
	assert.Equal(t, len(p.errors) > 0, true)
}

func assertNoErrors(t *testing.T, p *Parser) {
	assert.Equal(t, len(p.errors), 0)
}

func assertErrorAt(t *testing.T, p *Parser, index int, errSubStr string) {
	assert.Assert(t, len(p.errors) > index)
	err := p.errors[index].Error()
	assert.Assert(t, strings.Contains(err, errSubStr), err)
}

func assertProgram(t *testing.T, node *node.ProgramNode, hasContract bool) {
	assert.Equal(t, node.Contract != nil, hasContract)
}
func assertContract(t *testing.T, node *node.ContractNode, name string, totalVars int, totalFunctions int) {
	assert.Equal(t, node.Name, name)
	assert.Equal(t, len(node.Variables), totalVars)
	assert.Equal(t, len(node.Functions), totalFunctions)
}

func assertVariable(t *testing.T, node *node.VariableNode, varType string, id string) {
	assert.Equal(t, node.Type.Identifier, varType)
	assert.Equal(t, node.Identifier, id)
}

func assertFunction(t *testing.T, node *node.FunctionNode, name string, totalRTypes int, totalPTypes int, totalStmts int) {
	assert.Equal(t, node.Name, name)
	assert.Equal(t, len(node.ReturnTypes), totalRTypes)
	assert.Equal(t, len(node.Parameters), totalPTypes)
	assert.Equal(t, len(node.Body), totalStmts)
}

// Statements
// ----------

func assertStatementBlock(t *testing.T, node []node.StatementNode, totalStmt int) {
	assert.Equal(t, len(node), totalStmt)
}

func assertVariableStatement(t *testing.T, node *node.VariableNode, varType string, id string, expr string) {
	assert.Equal(t, node.Type.Identifier, varType)
	assert.Equal(t, node.Identifier, id)
	assertExpression(t, node.Expression, expr)
}

func assertAssignmentStatement(t *testing.T, node *node.AssignmentStatementNode, designator string, expr string) {
	assert.Equal(t, node.Left.Value, designator)
	assertExpression(t, node.Right, expr)
}

func assertIfStatement(t *testing.T, node *node.IfStatementNode, cond string, totalThen int, totalElse int) {
	assertExpression(t, node.Condition, cond)
	assert.Equal(t, len(node.Then), totalThen)
	assert.Equal(t, len(node.Else), totalElse)

}

func assertReturnStatement(t *testing.T, node *node.ReturnStatementNode, totalExpr int) {
	assert.Equal(t, len(node.Expressions), totalExpr)
}

func assertStatement(t *testing.T, node node.StatementNode, stmt string) {
	assert.Equal(t, node.String(), stmt)
}

// ----------

func assertDesignator(t *testing.T, node *node.DesignatorNode, value string) {
	assert.Equal(t, node.Value, value)
}

func assertIntegerLiteral(t *testing.T, node *node.IntegerLiteralNode, value *big.Int) {
	assert.Equal(t, node.Value.Cmp(value), 0)
}

func assertStringLiteral(t *testing.T, node *node.StringLiteralNode, value string) {
	assert.Equal(t, node.Value, value)
}

func assertCharacterLiteral(t *testing.T, node *node.CharacterLiteralNode, value rune) {
	assert.Equal(t, node.Value, value)
}

func assertBoolLiteral(t *testing.T, node *node.BoolLiteralNode, value bool) {
	assert.Equal(t, node.Value, value)
}

func assertError(t *testing.T, node *node.ErrorNode, message string) {
	assert.Equal(t, node.Message, message)
}

func assertType(t *testing.T, typeNode *node.TypeNode, varType string) {
	assert.Equal(t, typeNode.Identifier, varType)
}

func assertExpression(t *testing.T, node node.ExpressionNode, expr string) {
	assert.Equal(t, node.String(), expr)
}

func assertBinaryExpression(t *testing.T, n node.ExpressionNode, left string, right string, op token.Symbol) {
	binExpr, ok := n.(*node.BinaryExpressionNode)

	assert.Equal(t, ok, true)
	assert.Equal(t, binExpr.Left.String(), left)
	assert.Equal(t, binExpr.Right.String(), right)
	assert.Equal(t, binExpr.Operator, op)
}

func assertUnaryExpression(t *testing.T, n node.ExpressionNode, expr string, op token.Symbol) {
	unExpr, ok := n.(*node.UnaryExpressionNode)

	assert.Equal(t, ok, true)
	assert.Equal(t, unExpr.Expression.String(), expr)
	assert.Equal(t, unExpr.Operator, op)
}

func assertFuncCall(t *testing.T, n node.ExpressionNode, designator string, args ...string) {
	funcCall, ok := n.(*node.FuncCallNode)

	assert.Assert(t, ok)
	assert.Equal(t, funcCall.Designator.String(), designator)

	assert.Equal(t, len(funcCall.Args), len(args))
	for i, a := range args {
		assert.Equal(t, funcCall.Args[i].String(), a)
	}
}
