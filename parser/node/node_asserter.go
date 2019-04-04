package node

import (
	"github.com/bazo-blockchain/lazo/lexer/token"
	"gotest.tools/assert"
	"math/big"
	"testing"
)

func AssertProgram(t *testing.T, node *ProgramNode, hasContract bool) {
	assert.Equal(t, node.Contract != nil, hasContract)
}

func AssertContract(t *testing.T, node *ContractNode, name string, totalVars int, totalFunctions int) {
	assert.Equal(t, node.Name, name)
	assert.Equal(t, len(node.Variables), totalVars)
	assert.Equal(t, len(node.Functions), totalFunctions)
}

func AssertVariable(t *testing.T, node *VariableNode, varType string, id string) {
	assert.Equal(t, node.Type.Identifier, varType)
	assert.Equal(t, node.Identifier, id)
}

func AssertFunction(t *testing.T, node *FunctionNode, name string, totalRTypes int, totalPTypes int, totalStmts int) {
	assert.Equal(t, node.Name, name)
	assert.Equal(t, len(node.ReturnTypes), totalRTypes)
	assert.Equal(t, len(node.Parameters), totalPTypes)
	assert.Equal(t, len(node.Body), totalStmts)
}

// Statements
// ----------

func AssertStatementBlock(t *testing.T, node []StatementNode, totalStmt int) {
	assert.Equal(t, len(node), totalStmt)
}

func AssertVariableStatement(t *testing.T, node *VariableNode, varType string, id string, expr string) {
	assert.Equal(t, node.Type.Identifier, varType)
	assert.Equal(t, node.Identifier, id)
	AssertExpression(t, node.Expression, expr)
}

func AssertAssignmentStatement(t *testing.T, node *AssignmentStatementNode, designator string, expr string) {
	assert.Equal(t, node.Left.Value, designator)
	AssertExpression(t, node.Right, expr)
}

func AssertIfStatement(t *testing.T, node *IfStatementNode, cond string, totalThen int, totalElse int) {
	AssertExpression(t, node.Condition, cond)
	assert.Equal(t, len(node.Then), totalThen)
	assert.Equal(t, len(node.Else), totalElse)

}

func AssertReturnStatement(t *testing.T, node *ReturnStatementNode, totalExpr int) {
	assert.Equal(t, len(node.Expressions), totalExpr)
}

func AssertStatement(t *testing.T, node StatementNode, stmt string) {
	assert.Equal(t, node.String(), stmt)
}

// ----------

func AssertDesignator(t *testing.T, node *DesignatorNode, value string) {
	assert.Equal(t, node.Value, value)
}

func AssertIntegerLiteral(t *testing.T, node *IntegerLiteralNode, value *big.Int) {
	assert.Equal(t, node.Value.Cmp(value), 0)
}

func AssertStringLiteral(t *testing.T, node *StringLiteralNode, value string) {
	assert.Equal(t, node.Value, value)
}

func AssertCharacterLiteral(t *testing.T, node *CharacterLiteralNode, value rune) {
	assert.Equal(t, node.Value, value)
}

func AssertBoolLiteral(t *testing.T, node *BoolLiteralNode, value bool) {
	assert.Equal(t, node.Value, value)
}

func AssertError(t *testing.T, node *ErrorNode, message string) {
	assert.Equal(t, node.Message, message)
}

func AssertType(t *testing.T, typeNode *TypeNode, varType string) {
	assert.Equal(t, typeNode.Identifier, varType)
}

func AssertExpression(t *testing.T, node ExpressionNode, expr string) {
	assert.Equal(t, node.String(), expr)
}

func AssertBinaryExpression(t *testing.T, node ExpressionNode, left string, right string, op token.Symbol) {
	binExpr, ok := node.(*BinaryExpressionNode)

	assert.Equal(t, ok, true)
	assert.Equal(t, binExpr.Left.String(), left)
	assert.Equal(t, binExpr.Right.String(), right)
	assert.Equal(t, binExpr.Operator, op)
}

func AssertUnaryExpression(t *testing.T, node ExpressionNode, expr string, op token.Symbol) {
	unExpr, ok := node.(*UnaryExpressionNode)

	assert.Equal(t, ok, true)
	assert.Equal(t, unExpr.Expression.String(), expr)
	assert.Equal(t, unExpr.Operator, op)
}
