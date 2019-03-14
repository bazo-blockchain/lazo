package parser

import (
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"math/big"
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

// Term Expressions
// --------------------

func TestAddition(t *testing.T) {
	e := parseExpressionFromInput(t, "1 + 2")
	node.AssertBinaryExpression(t, e, "1", "2", token.Addition)
}

func TestConcatenation(t *testing.T) {
	e := parseExpressionFromInput(t, ` "hello" + "world" `)
	node.AssertBinaryExpression(t, e, "hello", "world", token.Addition)
}

func TestSubstraction(t *testing.T) {
	e := parseExpressionFromInput(t, "1 - 2")
	node.AssertBinaryExpression(t, e, "1", "2", token.Subtraction)
}

func TestTermAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "1 + 2 - 3")
	node.AssertBinaryExpression(t, e, "(1 + 2)", "3", token.Subtraction)
}

func TestTermPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "1 + 2 <= 3")
	node.AssertBinaryExpression(t, e, "(1 + 2)", "3", token.LessEqual)
}

// Factor Expressions
// ------------------

func TestMultiplication(t *testing.T) {
	e := parseExpressionFromInput(t, "3 * 4")
	node.AssertBinaryExpression(t, e, "3", "4", token.Multiplication)
}

func TestDivision(t *testing.T) {
	e := parseExpressionFromInput(t, "3 / 4")
	node.AssertBinaryExpression(t, e, "3", "4", token.Division)
}

func TestModulo(t *testing.T) {
	e := parseExpressionFromInput(t, "4 % 3")
	node.AssertBinaryExpression(t, e, "4", "3", token.Modulo)
}

func TestFactorAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "3 * 4 / 2 % 5")
	node.AssertBinaryExpression(t, e, "((3 * 4) / 2)", "5", token.Modulo)
}

func TestFactorPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "1 + 2 * 3")
	node.AssertBinaryExpression(t, e, "1", "(2 * 3)", token.Addition)
}

// Exponent Expressions
// --------------------

func TestExponent(t *testing.T) {
	e := parseExpressionFromInput(t, "2 ** 3 ** 4")
	node.AssertBinaryExpression(t, e, "2", "(3 ** 4)", token.Exponent)
}

func TestExponentWithFactor(t *testing.T) {
	e := parseExpressionFromInput(t, "2 ** 3 * 4")
	node.AssertBinaryExpression(t, e, "(2 ** 3)", "4", token.Multiplication)
}

func TestFactorWithExponent(t *testing.T) {
	e := parseExpressionFromInput(t, "2 / 3 ** 4")
	node.AssertBinaryExpression(t, e, "2", "(3 ** 4)", token.Division)
}

// Unary Expressions
// -----------------

func TestUnaryPlus(t *testing.T) {
	e := parseExpressionFromInput(t, "+x")
	node.AssertUnaryExpression(t, e, "x", token.Addition)
}

func TestUnaryMinus(t *testing.T) {
	e := parseExpressionFromInput(t, "2 - -3")
	node.AssertBinaryExpression(t, e, "2", "(-3)", token.Subtraction)

	unExpr := e.(*node.BinaryExpressionNode).Right
	node.AssertUnaryExpression(t, unExpr, "3", token.Subtraction)
}

func TestUnaryNot(t *testing.T) {
	e := parseExpressionFromInput(t, "!true")
	node.AssertUnaryExpression(t, e, "true", token.Not)
}

func TestUnaryPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "-4 + 2")
	node.AssertBinaryExpression(t, e, "(-4)", "2", token.Addition)
}

func TestUnaryWithFactor(t *testing.T) {
	e := parseExpressionFromInput(t, "-4 * 2")
	node.AssertUnaryExpression(t, e, "(4 * 2)", token.Subtraction)
}

func TestUnaryAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "-+-+x")
	node.AssertUnaryExpression(t, e, "(+(-(+x)))", token.Subtraction)
}

// Designator Expressions
// ----------------------

func TestDesignator(t *testing.T) {
	p := newParserFromInput("test")
	v := p.parseDesignator()
	node.AssertDesignator(t, v, "test")
	assertNoErrors(t, p)
}

func TestDesignatorWithNumbers(t *testing.T) {
	p := newParserFromInput("test123")
	v := p.parseDesignator()
	node.AssertDesignator(t, v, "test123")
	assertNoErrors(t, p)
}

func TestDesignatorWithUnderscore(t *testing.T) {
	p := newParserFromInput("_test")
	v := p.parseDesignator()
	node.AssertDesignator(t, v, "_test")
	assertNoErrors(t, p)
}

func TestInvalidDesignator(t *testing.T) {
	p := newParserFromInput("1ab")
	p.parseDesignator()
	assertHasError(t, p)
}

// Literal Expressions
// -------------------

func TestIntegerLiteral(t *testing.T) {
	p := newParserFromInput("1")
	i := p.parseInteger()
	node.AssertIntegerLiteral(t, i, big.NewInt(1))
	assertNoErrors(t, p)
}

func TestValidIntegerLiteral(t *testing.T) {
	p := newParserFromInput("0x1")
	i := p.parseInteger()
	node.AssertIntegerLiteral(t, i, big.NewInt(1))
	assertNoErrors(t, p)
}

func TestInvalidHexLiteral(t *testing.T) {
	p := newParserFromInput("0x")
	o := p.parseOperand()

	assertHasError(t, p)
	e := o.(*node.ErrorNode)
	node.AssertError(t, e, "Error while parsing string to big int")
}

func TestStringLiteral(t *testing.T) {
	p := newParserFromInput(`"test"`)
	s := p.parseString()
	node.AssertStringLiteral(t, s, "test")
	assertNoErrors(t, p)
}

func TestCharacterLiteral(t *testing.T) {
	p := newParserFromInput("'c'")
	c := p.parseCharacter()
	node.AssertCharacterLiteral(t, c, 'c')
	assertNoErrors(t, p)
}

func TestBoolLiteralTrue(t *testing.T) {
	p := newParserFromInput("true")
	b := p.parseBoolean()
	tok, _ := b.(*node.BoolLiteralNode)
	node.AssertBoolLiteral(t, tok, true)
	assertNoErrors(t, p)
}

func TestBoolLiteralFalse(t *testing.T) {
	p := newParserFromInput("false")
	b := p.parseBoolean()
	tok, _ := b.(*node.BoolLiteralNode)
	node.AssertBoolLiteral(t, tok, false)
	assertNoErrors(t, p)
}

func TestInvalidBoolLiteral(t *testing.T) {
	p := newParserFromInput("if")
	b := p.parseBoolean()

	assertHasError(t, p)
	tok, _ := b.(*node.ErrorNode)
	node.AssertError(t, tok, "Invalid boolean value if")
}

// --------------

func parseExpressionFromInput(t *testing.T, input string) node.ExpressionNode {
	p := newParserFromInput(input)
	assertNoErrors(t, p)
	return p.parseExpression()
}
