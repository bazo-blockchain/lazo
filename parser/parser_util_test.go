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
	assert.Equal(t, len(p.errors), 0, p.errors)
}

func assertErrorAt(t *testing.T, p *Parser, index int, errSubStr string) {
	assert.Assert(t, len(p.errors) > index)
	err := p.errors[index].Error()
	assert.Assert(t, strings.Contains(err, errSubStr), err)
}

func assertPosition(t *testing.T, actualPos token.Position, line int, col int) {
	assert.Equal(t, actualPos.Line, line)
	assert.Equal(t, actualPos.Column, col)
}

func assertProgram(t *testing.T, node *node.ProgramNode, hasContract bool) {
	assert.Equal(t, node.Contract != nil, hasContract)
}
func assertContract(t *testing.T, node *node.ContractNode, name string, totalVars int, totalFunctions int) {
	assert.Equal(t, node.Name, name)
	assert.Equal(t, len(node.Fields), totalVars)
	assert.Equal(t, len(node.Functions), totalFunctions)
}

func assertField(t *testing.T, node *node.FieldNode, varType string, id string, expr string) {
	assert.Equal(t, node.Type.String(), varType)
	assert.Equal(t, node.Identifier, id)
	if expr != "" {
		assertExpression(t, node.Expression, expr)
	}
}

func assertStruct(t *testing.T, node *node.StructNode, name string, totalFields int) {
	assert.Equal(t, node.Name, name)
	assert.Equal(t, len(node.Fields), totalFields)
}

func assertStructField(t *testing.T, node *node.StructFieldNode, varType string, id string) {
	assert.Equal(t, node.Type.String(), varType)
	assert.Equal(t, node.Identifier, id)
}

func assertConstructor(t *testing.T, node *node.ConstructorNode, totalPTypes int, totalStmts int) {
	assert.Equal(t, len(node.Parameters), totalPTypes)
	assert.Equal(t, len(node.Body), totalStmts)
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
	assert.Equal(t, node.Type.String(), varType)
	assert.Equal(t, node.Identifier, id)
	assertExpression(t, node.Expression, expr)
}

func assertAssignmentStatement(t *testing.T, node *node.AssignmentStatementNode, designator string, expr string) {
	assert.Equal(t, node.Left.String(), designator)
	assertExpression(t, node.Right, expr)
}

func assertShorthandAssignmentStatement(t *testing.T, node *node.ShorthandAssignmentStatementNode,
	designator string, expr string, operator token.Symbol) {
	assert.Equal(t, node.Designator.String(), designator)
	assertExpression(t, node.Expression, expr)
	assert.Equal(t, node.Operator, operator, token.SymbolLexeme[node.Operator])
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

func assertDesignator(t *testing.T, node *node.BasicDesignatorNode, value string) {
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

func assertType(t *testing.T, typeNode node.TypeNode, varType string) {
	assert.Equal(t, typeNode.String(), varType)
}

func assertExpression(t *testing.T, node node.ExpressionNode, expr string) {
	if expr == "" {
		assert.Equal(t, node, nil)
	} else {
		assert.Equal(t, node.String(), expr)
	}
}

func assertTernaryExpression(t *testing.T, n node.ExpressionNode, condition string, trueExpr string, falseExpr string) {
	ternary, ok := n.(*node.TernaryExpressionNode)

	assert.Equal(t, ok, true)
	assert.Equal(t, ternary.Condition.String(), condition)
	assert.Equal(t, ternary.Then.String(), trueExpr)
	assert.Equal(t, ternary.Else.String(), falseExpr)
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

func assertElementAccess(t *testing.T, n node.ExpressionNode, designator string, exp string) {
	elementAccess, ok := n.(*node.ElementAccessNode)

	assert.Assert(t, ok)
	assert.Equal(t, elementAccess.Designator.String(), designator)
	assert.Equal(t, elementAccess.Expression.String(), exp)
}

func assertMemberAccess(t *testing.T, n node.ExpressionNode, designator string, id string) {
	memberAccess, ok := n.(*node.MemberAccessNode)

	assert.Assert(t, ok)
	assert.Equal(t, memberAccess.Designator.String(), designator)
	assert.Equal(t, memberAccess.Identifier, id)
}

// assertArrayLengthCreation should receive the string lengths in the following format: "1,2,3,..."
func assertArrayLengthCreation(t *testing.T, n node.ExpressionNode, name string, lengths string) {
	arrayCreation, ok := n.(*node.ArrayLengthCreationNode)

	assert.Assert(t, ok)
	assert.Equal(t, arrayCreation.ElementType.String(), name)

	assertArrayLengths(t, arrayCreation, lengths)
}

func assertArrayLengths(t *testing.T, n *node.ArrayLengthCreationNode, lengths string) {
	var lengthStrings []string
	for _, length := range n.Lengths {
		lengthStrings = append(lengthStrings, length.String())
	}

	result := strings.Join(lengthStrings, ",")

	assert.Equal(t, result, lengths)
}

func assertArrayValueCreation(t *testing.T, n node.ExpressionNode, name string, values ...string) {
	arrayCreation, ok := n.(*node.ArrayValueCreationNode)

	assert.Assert(t, ok)
	assert.Equal(t, arrayCreation.Type.String(), name)

	assert.Equal(t, len(arrayCreation.Elements.Values), len(values))
	for i, v := range values {
		assert.Equal(t, arrayCreation.Elements.Values[i].String(), v)
	}
}

func assertStructCreation(t *testing.T, n node.ExpressionNode, name string, values ...string) {
	structCreation, ok := n.(*node.StructCreationNode)

	assert.Assert(t, ok)
	assert.Equal(t, structCreation.Name, name)

	assert.Equal(t, len(structCreation.FieldValues), len(values))
	for i, v := range values {
		assert.Equal(t, structCreation.FieldValues[i].String(), v)
	}
}

func assertStructNamedCreation(t *testing.T, n node.ExpressionNode, name string, values ...pair) {
	structCreation, ok := n.(*node.StructNamedCreationNode)

	assert.Assert(t, ok)
	assert.Equal(t, structCreation.Name, name)

	assert.Equal(t, len(structCreation.FieldValues), len(values))
	for i, pair := range values {
		field := structCreation.FieldValues[i]
		assert.Equal(t, field.Name, pair.key)
		assert.Equal(t, field.Expression.String(), pair.value)
		i++
	}
}

func assertTypeCast(t *testing.T, n node.ExpressionNode, typeName string, expr string) {
	typeCast, ok := n.(*node.TypeCastNode)

	assert.Assert(t, ok)
	assert.Equal(t, typeCast.Type.String(), typeName)
	assertExpression(t, typeCast.Expression, expr)
}

type pair struct {
	key   string
	value string
}
