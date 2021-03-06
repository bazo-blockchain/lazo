package parser

import (
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"math/big"
	"testing"
)

// Ternary Expressions
// ==============================

func TestTernaryExpression(t *testing.T) {
	e := parseExpressionFromInput(t, "x == y ? true : false")
	assertTernaryExpression(t, e, "(x == y)", "true", "false")
}

func TestNestedTernary(t *testing.T) {
	p := newParserFromInput("x ? y ? 1 : 2 : z")
	_ = p.parseExpression()

	assertErrorAt(t, p, 0, "Symbol : expected, but got ?")
}

// Binary Expressions
// ==============================

// Logic operators
// ---------------

func TestLogicOr(t *testing.T) {
	e := parseExpressionFromInput(t, "x || y || z")
	assertBinaryExpression(t, e, "(x || y)", "z", token.Or)
}

func TestLogicAnd(t *testing.T) {
	e := parseExpressionFromInput(t, "x && y && z")
	assertBinaryExpression(t, e, "(x && y)", "z", token.And)
}

func TestMixedLogicOperatorsAndOr(t *testing.T) {
	e := parseExpressionFromInput(t, "x && y || z")
	assertBinaryExpression(t, e, "(x && y)", "z", token.Or)
}

func TestMixedLogicOperatorOrAnd(t *testing.T) {
	e := parseExpressionFromInput(t, "x || y && z")
	assertBinaryExpression(t, e, "x", "(y && z)", token.Or)
}

// Bitwise Logic operators
// -----------------------

func TestBitwiseOr(t *testing.T) {
	e := parseExpressionFromInput(t, "x | y | z")
	assertBinaryExpression(t, e, "(x | y)", "z", token.BitwiseOr)
}

func TestBitwiseXOr(t *testing.T) {
	e := parseExpressionFromInput(t, "x ^ y ^ z")
	assertBinaryExpression(t, e, "(x ^ y)", "z", token.BitwiseXOr)
}

func TestBitwiseAnd(t *testing.T) {
	e := parseExpressionFromInput(t, "x & y & z")
	assertBinaryExpression(t, e, "(x & y)", "z", token.BitwiseAnd)
}

func TestMixedBitwiseLogicOperatorsAndOrXor(t *testing.T) {
	e := parseExpressionFromInput(t, "w | x ^ y & z")
	assertBinaryExpression(t, e, "w", "(x ^ (y & z))", token.BitwiseOr)
}

func TestMixedBitwiseLogicOperatorsAndOrXor2(t *testing.T) {
	e := parseExpressionFromInput(t, "w & x | y ^ z")
	assertBinaryExpression(t, e, "(w & x)", "(y ^ z)", token.BitwiseOr)
}

// Equality operators
// ------------------

func TestEquality(t *testing.T) {
	e := parseExpressionFromInput(t, "5 == 4")
	assertBinaryExpression(t, e, "5", "4", token.Equal)
}

func TestUnequlity(t *testing.T) {
	e := parseExpressionFromInput(t, "5 != 4")
	assertBinaryExpression(t, e, "5", "4", token.Unequal)
}

func TestEqualityAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "x == y == true")
	assertBinaryExpression(t, e, "(x == y)", "true", token.Equal)
}

// Relational operators
// --------------------

func TestLess(t *testing.T) {
	e := parseExpressionFromInput(t, "1 < 3")
	assertBinaryExpression(t, e, "1", "3", token.Less)
}

func TestLessEqual(t *testing.T) {
	e := parseExpressionFromInput(t, "1 <= 3")
	assertBinaryExpression(t, e, "1", "3", token.LessEqual)
}

func TestGreater(t *testing.T) {
	e := parseExpressionFromInput(t, "1 > 3")
	assertBinaryExpression(t, e, "1", "3", token.Greater)
}

func TestGreaterEqual(t *testing.T) {
	e := parseExpressionFromInput(t, "1 >= 3")
	assertBinaryExpression(t, e, "1", "3", token.GreaterEqual)
}

func TestRelationalComparisonAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "1 < 2 <= 3 > 4 >= 5")
	assertBinaryExpression(t, e, "(((1 < 2) <= 3) > 4)", "5", token.GreaterEqual)
}

// Bitwise Shift Operators
// -----------------------

func TestShiftLeft(t *testing.T) {
	e := parseExpressionFromInput(t, "13 << 2")
	assertBinaryExpression(t, e, "13", "2", token.ShiftLeft)
}

func TestShiftRight(t *testing.T) {
	e := parseExpressionFromInput(t, "13 >> 2")
	assertBinaryExpression(t, e, "13", "2", token.ShiftRight)
}

func TestShiftAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "13 << 3 >> 1")
	assertBinaryExpression(t, e, "(13 << 3)", "1", token.ShiftRight)
}

// Term Expressions
// --------------------

func TestAddition(t *testing.T) {
	e := parseExpressionFromInput(t, "1 + 2")
	assertBinaryExpression(t, e, "1", "2", token.Plus)
}

func TestConcatenation(t *testing.T) {
	e := parseExpressionFromInput(t, ` "hello" + "world" `)
	assertBinaryExpression(t, e, "hello", "world", token.Plus)
}

func TestSubstraction(t *testing.T) {
	e := parseExpressionFromInput(t, "1 - 2")
	assertBinaryExpression(t, e, "1", "2", token.Minus)
}

func TestTermAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "1 + 2 - 3")
	assertBinaryExpression(t, e, "(1 + 2)", "3", token.Minus)
}

// Factor Expressions
// ------------------

func TestMultiplication(t *testing.T) {
	e := parseExpressionFromInput(t, "3 * 4")
	assertBinaryExpression(t, e, "3", "4", token.Multiplication)
}

func TestDivision(t *testing.T) {
	e := parseExpressionFromInput(t, "3 / 4")
	assertBinaryExpression(t, e, "3", "4", token.Division)
}

func TestModulo(t *testing.T) {
	e := parseExpressionFromInput(t, "4 % 3")
	assertBinaryExpression(t, e, "4", "3", token.Modulo)
}

func TestFactorAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "3 * 4 / 2 % 5")
	assertBinaryExpression(t, e, "((3 * 4) / 2)", "5", token.Modulo)
}

// Exponent Expressions
// --------------------

func TestExponent(t *testing.T) {
	e := parseExpressionFromInput(t, "2 ** 3 ** 4")
	assertBinaryExpression(t, e, "2", "(3 ** 4)", token.Exponent)
}

func TestExponentWithFactor(t *testing.T) {
	e := parseExpressionFromInput(t, "2 ** 3 * 4")
	assertBinaryExpression(t, e, "(2 ** 3)", "4", token.Multiplication)
}

func TestFactorWithExponent(t *testing.T) {
	e := parseExpressionFromInput(t, "2 / 3 ** 4")
	assertBinaryExpression(t, e, "2", "(3 ** 4)", token.Division)
}

// Type Cast Expressions
// ---------------------

func TestIntTypeCast(t *testing.T) {
	e := parseExpressionFromInput(t, "(String) 5")
	assertTypeCast(t, e, "String", "5")
}

func TestCharTypeCast(t *testing.T) {
	e := parseExpressionFromInput(t, "(String) 'c'")
	assertTypeCast(t, e, "String", "c")
}

func TestBoolTypeCast(t *testing.T) {
	e := parseExpressionFromInput(t, "(String) true")
	assertTypeCast(t, e, "String", "true")
}

func TestDesignatorTypeCast(t *testing.T) {
	e := parseExpressionFromInput(t, "(String) x")
	assertTypeCast(t, e, "String", "x")
}

func TestTypeCastRightAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "(String) x.y.z")
	assertTypeCast(t, e, "String", "x.y.z")
}

func TestTypeCastStringConcat(t *testing.T) {
	e := parseExpressionFromInput(t, `"Hello " + (String) 5 + (String) true`)
	assertBinaryExpression(t, e, "(Hello  + (String) 5)", "(String) true", token.Plus)
}

// Unary Expressions
// -----------------

func TestUnaryPlus(t *testing.T) {
	e := parseExpressionFromInput(t, "+x")
	assertUnaryExpression(t, e, "x", token.Plus)
}

func TestUnaryMinus(t *testing.T) {
	e := parseExpressionFromInput(t, "2 - -3")
	assertBinaryExpression(t, e, "2", "(-3)", token.Minus)

	unExpr := e.(*node.BinaryExpressionNode).Right
	assertUnaryExpression(t, unExpr, "3", token.Minus)
}

func TestUnaryNot(t *testing.T) {
	e := parseExpressionFromInput(t, "!true")
	assertUnaryExpression(t, e, "true", token.Not)
}

func TestUnaryBinaryNot(t *testing.T) {
	e := parseExpressionFromInput(t, "~x")
	assertUnaryExpression(t, e, "x", token.BitwiseNot)
}

func TestUnaryWithFactor(t *testing.T) {
	e := parseExpressionFromInput(t, "-4 * 2")
	assertUnaryExpression(t, e, "(4 * 2)", token.Minus)
}

func TestUnaryAssociativity(t *testing.T) {
	e := parseExpressionFromInput(t, "-+-+x")
	assertUnaryExpression(t, e, "(+(-(+x)))", token.Minus)
}

// Designator Expressions
// ----------------------

func TestDesignator(t *testing.T) {
	p := newParserFromInput("test")
	v := p.parseDesignator()
	assertDesignator(t, v.(*node.BasicDesignatorNode), "test")
	assertNoErrors(t, p)
}

func TestDesignatorWithNumbers(t *testing.T) {
	p := newParserFromInput("test123")
	v := p.parseDesignator()
	assertDesignator(t, v.(*node.BasicDesignatorNode), "test123")
	assertNoErrors(t, p)
}

func TestDesignatorWithUnderscore(t *testing.T) {
	p := newParserFromInput("_test")
	v := p.parseDesignator()
	assertDesignator(t, v.(*node.BasicDesignatorNode), "_test")
	assertNoErrors(t, p)
}

func TestInvalidDesignator(t *testing.T) {
	p := newParserFromInput("1ab")
	p.parseDesignator()
	assertHasError(t, p)
}

// Function Call Expressions
// --------------------------

func TestFuncCall(t *testing.T) {
	e := parseExpressionFromInput(t, "test()")
	assertFuncCall(t, e, "test")
}

func TestFuncCallWithParam(t *testing.T) {
	e := parseExpressionFromInput(t, "test(1)")
	assertFuncCall(t, e, "test", "1")
}

func TestFuncCallWithParams(t *testing.T) {
	e := parseExpressionFromInput(t, "test(x, 2 * 3)")
	assertFuncCall(t, e, "test", "x", "(2 * 3)")
}

func TestFuncCallMissingCloseParen(t *testing.T) {
	p := newParserFromInput("test(")
	_ = p.parseExpression()
	assertErrorAt(t, p, 0, ") expected, but got EOF")
}

func TestFuncCallMissingComma(t *testing.T) {
	p := newParserFromInput("test(1 2)")
	_ = p.parseExpression()
	assertErrorAt(t, p, 0, ", expected, but got 2")
}

func TestFuncCallOnMember(t *testing.T) {
	e := parseExpressionFromInput(t, "a.f()")
	assertFuncCall(t, e, "a.f")
	assertPosition(t, e.Pos(), 1, 1)
}

// Element Access Expression
// -------------------------

func TestElementAccess(t *testing.T) {
	e := parseExpressionFromInput(t, "arr[0]")
	assertElementAccess(t, e, "arr", "0")
	assertPosition(t, e.Pos(), 1, 1)
}

func TestElementAccessWithDesignator(t *testing.T) {
	e := parseExpressionFromInput(t, "arr[pos]")
	assertElementAccess(t, e, "arr", "pos")
}

func TestElementAccessWithExpression(t *testing.T) {
	e := parseExpressionFromInput(t, "arr[1 + 2]")
	assertElementAccess(t, e, "arr", "(1 + 2)")
}

func TestElementAccessOnMember(t *testing.T) {
	e := parseExpressionFromInput(t, "a.arr[pos]")
	assertElementAccess(t, e, "a.arr", "pos")
}

func TestInvalidElementAccess(t *testing.T) {
	p := newParserFromInput("a.arr[if]")
	p.parseExpression()
	assertErrorAt(t, p, 0, "Unsupported expression symbol if")
}

func TestMapElementAccess(t *testing.T) {
	e := parseExpressionFromInput(t, `map["key"]`)
	assertElementAccess(t, e, "map", "key")
	assertPosition(t, e.Pos(), 1, 1)
}

// Member Access Expression
// ------------------------

func TestMemberAccess(t *testing.T) {
	e := parseExpressionFromInput(t, "a.x")
	assertMemberAccess(t, e, "a", "x")
}

func TestMultipleMemberAccess(t *testing.T) {
	e := parseExpressionFromInput(t, "a.x.y")
	assertMemberAccess(t, e, "a.x", "y")
}

func TestInvalidMemberAccess(t *testing.T) {
	p := newParserFromInput("a.0")
	p.parseExpression()
	assertErrorAt(t, p, 0, "Identifier expected")
}

// Struct Creation Expressions
// ---------------------------

func TestUnsupportedCreation(t *testing.T) {
	p := newParserFromInput("new Person{}")
	p.parseCreation()

	assertErrorAt(t, p, 0, "Unsupported creation type with {")
	assert.Equal(t, len(p.errors), 1)
}

func TestStructCreation(t *testing.T) {
	s := parseExpressionFromInput(t, "new Person()")
	assertPosition(t, s.Pos(), 1, 1)
	assertStructCreation(t, s, "Person")
}

func TestStructCreationWithValues(t *testing.T) {
	s := parseExpressionFromInput(t, "new Person(120, 1 == 1)")
	assertStructCreation(t, s, "Person", "120", "(1 == 1)")
}

func TestStructCreationWithNamedFieldValues(t *testing.T) {
	s := parseExpressionFromInput(t, "new Person(x=120, y=1==1)")

	expectedFieldValues := []pair{
		{"x", "120"},
		{"y", "(1 == 1)"},
	}
	assertPosition(t, s.Pos(), 1, 1)
	assertStructNamedCreation(t, s, "Person", expectedFieldValues...)
}

func TestStructCreationWithNewlines(t *testing.T) {
	s := parseExpressionFromInput(t, `new Person(
		x=120, 
		y=1==1
	)`)

	expectedFieldValues := []pair{
		{"x", "120"},
		{"y", "(1 == 1)"},
	}
	assertStructNamedCreation(t, s, "Person", expectedFieldValues...)
}

// Array Nodes
// -----------

func TestArrayNewArrayAssignment1(t *testing.T) {
	p := parseExpressionFromInput(t, "new int[2]")

	assertArrayLengthCreation(t, p, "int", "2")
}

func TestArrayNewArrayAssignment2(t *testing.T) {
	p := parseExpressionFromInput(t, "new int[]{1,2}")

	assertArrayValueCreation(t, p, "int[]", "1", "2")
}

func TestNestedArrayNewArrayAssignment1(t *testing.T) {
	p := parseExpressionFromInput(t, "new int[1][2]")

	assertArrayLengthCreation(t, p, "int[]", "1,2")
}

func TestNestedArrayNewArrayAssignment2(t *testing.T) {
	p := parseExpressionFromInput(t, "new int[][]{{1, 2}, {3, 4}}")

	assertArrayValueCreation(t, p, "int[][]", "[[1 2]]", "[[3 4]]")
}

func TestNestedArrayNewArrayAssignment3(t *testing.T) {
	p := parseExpressionFromInput(t, "new int[][]{{1, 2}, {3}}")

	assertArrayValueCreation(t, p, "int[][]", "[[1 2]]", "[[3]]")
}

func TestNestedArrayNewArrayAssignment4(t *testing.T) {
	p := newParserFromInput("new int[]{}")
	p.parseExpression()

	assertErrorAt(t, p, 0, "Unsupported expression symbol }")
}

func TestNestedArrayNewArrayAssignment5(t *testing.T) {
	p := newParserFromInput("new int[][]{}")
	p.parseExpression()

	assertErrorAt(t, p, 0, "Unsupported expression symbol }")
}

func TestNestedArrayNewArrayAssignment6(t *testing.T) {
	p := newParserFromInput("new int[][]{{},{}}")
	p.parseExpression()

	assertErrorAt(t, p, 0, "Unsupported expression symbol }")
}

func TestNestedArrayNewArrayAssignment7(t *testing.T) {
	p := newParserFromInput("new int[][]{{1},{}}")
	p.parseExpression()

	assertErrorAt(t, p, 0, "Unsupported expression symbol }")
}

func TestInvalidLengthArrayNewArrayAssignment2(t *testing.T) {
	p := newParserFromInput("new int[]")
	p.parseCreation()
	assertErrorAt(t, p, 0, "Symbol { expected, but got EOF")
	assertErrorAt(t, p, 1, "Invalid array initialization")
}

func TestArrayValueAssignment(t *testing.T) {
	p := newParserFromInput("a[0] = 1\n")

	stmt := p.parseStatementWithIdentifier()
	assignment := stmt.(*node.AssignmentStatementNode)
	assertAssignmentStatement(t, assignment, "a[0]", "1")
	assertNoErrors(t, p)
}

func TestArrayValueAssignmentNegativeIndex(t *testing.T) {
	p := newParserFromInput("a[-1] = 1\n")

	stmt := p.parseStatementWithIdentifier()
	assignment := stmt.(*node.AssignmentStatementNode)
	assertAssignmentStatement(t, assignment, "a[(-1)]", "1")
	assertNoErrors(t, p)
}

// Parentheses
// ------------

func TestParentheses(t *testing.T) {
	e := parseExpressionFromInput(t, "(x + y)")
	assertBinaryExpression(t, e, "x", "y", token.Plus)
}

func TestNestedParentheses(t *testing.T) {
	e := parseExpressionFromInput(t, "((2) * (3 + 4))")
	assertBinaryExpression(t, e, "2", "(3 + 4)", token.Multiplication)
}

// Literal Expressions
// -------------------

func TestIntegerLiteral(t *testing.T) {
	p := newParserFromInput("1")
	i := p.parseInteger()
	assertIntegerLiteral(t, i, big.NewInt(1))
	assertNoErrors(t, p)
}

func TestValidIntegerLiteral(t *testing.T) {
	p := newParserFromInput("0x1")
	i := p.parseInteger()
	assertIntegerLiteral(t, i, big.NewInt(1))
	assertNoErrors(t, p)
}

func TestInvalidHexLiteral(t *testing.T) {
	p := newParserFromInput("0x")
	o := p.parseOperand()

	assertHasError(t, p)
	e := o.(*node.ErrorNode)
	assertError(t, e, "Error while parsing string to big int")
}

func TestStringLiteral(t *testing.T) {
	p := newParserFromInput(`"test"`)
	s := p.parseString()
	assertStringLiteral(t, s, "test")
	assertNoErrors(t, p)
}

func TestCharacterLiteral(t *testing.T) {
	p := newParserFromInput("'c'")
	c := p.parseCharacter()
	assertCharacterLiteral(t, c, 'c')
	assertNoErrors(t, p)
}

func TestBoolLiteralTrue(t *testing.T) {
	p := newParserFromInput("true")
	b := p.parseOperandSymbol()
	tok, _ := b.(*node.BoolLiteralNode)
	assertBoolLiteral(t, tok, true)
	assertNoErrors(t, p)
}

func TestBoolLiteralFalse(t *testing.T) {
	p := newParserFromInput("false")
	b := p.parseOperandSymbol()
	tok, _ := b.(*node.BoolLiteralNode)
	assertBoolLiteral(t, tok, false)
	assertNoErrors(t, p)
}

func TestInvalidBoolLiteral(t *testing.T) {
	p := newParserFromInput("if")
	b := p.parseBoolean(p.currentToken.(*token.FixToken))

	assertHasError(t, p)
	tok, _ := b.(*node.ErrorNode)
	assertError(t, tok, "Invalid boolean value if")
}

// Precedence
// -----------------

func TestElementAccessPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "2 * a[0]")
	assertBinaryExpression(t, e, "2", "a[0]", token.Multiplication)
}

func TestMemberAccessPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "2 * a.b")
	assertBinaryExpression(t, e, "2", "a.b", token.Multiplication)
}

func TestParenthesesPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "2 * (3 + 4)")
	assertBinaryExpression(t, e, "2", "(3 + 4)", token.Multiplication)
}

func TestFuncCallPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "2 * a()")
	assertBinaryExpression(t, e, "2", "a([])", token.Multiplication)
}

func TestUnaryPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "-4 + 2")
	assertBinaryExpression(t, e, "(-4)", "2", token.Plus)
}

func TestTypeCastPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "(String) x + y")
	assertBinaryExpression(t, e, "(String) x", "y", token.Plus)
}

func TestExponentiationPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "2 * 3 ** 4")
	assertBinaryExpression(t, e, "2", "(3 ** 4)", token.Multiplication)
}

func TestFactorPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "1 + 2 * 3")
	assertBinaryExpression(t, e, "1", "(2 * 3)", token.Plus)
}

func TestTermPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "1 + 2 <= 3")
	assertBinaryExpression(t, e, "(1 + 2)", "3", token.LessEqual)
}

func TestShiftPrecedence(t *testing.T) {
	// precedence order: +, <<, <
	e := parseExpressionFromInput(t, "2 < 3 << 3 + 4")
	assertBinaryExpression(t, e, "2", "(3 << (3 + 4))", token.Less)
}

func TestComparisonPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "false == 5 > 3")
	assertBinaryExpression(t, e, "false", "(5 > 3)", token.Equal)
}

func TestEqualityPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "false & 3 == 2")
	assertBinaryExpression(t, e, "false", "(3 == 2)", token.BitwiseAnd)
}

func TestBitwiseAndPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "2 ^ 5 & 3")
	assertBinaryExpression(t, e, "2", "(5 & 3)", token.BitwiseXOr)
}

func TestBitwiseXOrPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "2 | 5 ^ 3")
	assertBinaryExpression(t, e, "2", "(5 ^ 3)", token.BitwiseOr)
}

func TestBitwiseOrPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "2 && 5 | 3")
	assertBinaryExpression(t, e, "2", "(5 | 3)", token.And)
}

func TestLogicalPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "false || true && false")
	assertBinaryExpression(t, e, "false", "(true && false)", token.Or)
}

func TestTernaryExpressionPrecedence(t *testing.T) {
	e := parseExpressionFromInput(t, "x == y ? x + 1 : 2 * 3 + 4")
	assertTernaryExpression(t, e, "(x == y)", "(x + 1)", "((2 * 3) + 4)")
}

// --------------

func parseExpressionFromInput(t *testing.T, input string) node.ExpressionNode {
	p := newParserFromInput(input)
	expr := p.parseExpression()
	assertNoErrors(t, p)
	return expr
}
