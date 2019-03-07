package parser

import (
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"testing"
)

// Binary Expressions
// ==============================

// Logic operators
// ---------------

func TestLogicOr(t *testing.T) {
	e := parseExpressionFromInput(t, "x || y || z")
	node.AssertBinaryExpression(t, e, "(x || y)", "z", token.Or)
}

func TestLogicAnd(t *testing.T) {
	e := parseExpressionFromInput(t, "x && y && z")
	node.AssertBinaryExpression(t, e, "(x && y)", "z", token.And)
}

func TestMixedLogicOperatorsAndOr(t *testing.T) {
	e := parseExpressionFromInput(t, "x && y || z")
	node.AssertBinaryExpression(t, e, "(x && y)", "z", token.Or)
}

func TestMixedLogicOperatorOrAnd(t *testing.T) {
	e := parseExpressionFromInput(t, "x || y && z")
	node.AssertBinaryExpression(t, e, "x", "(y && z)", token.Or)
}

// Equality operators
// ------------------

func TestEquality(t *testing.T) {
	e := parseExpressionFromInput(t, "5 == 4")
	node.AssertBinaryExpression(t, e, "5", "4", token.Equal)
}

func TestUnequlity(t *testing.T) {
	e := parseExpressionFromInput(t, "5 != 4")
	node.AssertBinaryExpression(t, e, "5", "4", token.Unequal)
}

func TestEqualityPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "false == 5 > 3")
	node.AssertBinaryExpression(t, e, "false", "(5 > 3)", token.Equal)
}

func TestEqualityAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "x == y == true")
	node.AssertBinaryExpression(t, e, "(x == y)", "true", token.Equal)
}

// Relational operators
// --------------------

func TestLess(t *testing.T) {
	e := parseExpressionFromInput(t, "1 < 3")
	node.AssertBinaryExpression(t, e, "1", "3", token.Less)
}

func TestLessEqual(t *testing.T) {
	e := parseExpressionFromInput(t, "1 <= 3")
	node.AssertBinaryExpression(t, e, "1", "3", token.LessEqual)
}

func TestGreater(t *testing.T) {
	e := parseExpressionFromInput(t, "1 > 3")
	node.AssertBinaryExpression(t, e, "1", "3", token.Greater)
}

func TestGreaterEqual(t *testing.T) {
	e := parseExpressionFromInput(t, "1 >= 3")
	node.AssertBinaryExpression(t, e, "1", "3", token.GreaterEqual)
}

func TestRelationalComparisonAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "1 < 2 <= 3 > 4 >= 5")
	node.AssertBinaryExpression(t, e, "(((1 < 2) <= 3) > 4)", "5", token.GreaterEqual)
}

func TestExpr(t *testing.T) {
	e := parseExpressionFromInput(t, "1")
	node.AssertExpression(t, e, "1")
}

// --------------

func parseExpressionFromInput(t *testing.T, input string) node.ExpressionNode {
	p := newParserFromInput(input)
	assertNoErrors(t, p)
	return p.parseExpression()
}
