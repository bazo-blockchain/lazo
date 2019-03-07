package parser

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"math/big"
	"strings"
	"testing"
)

// Program Nodes
// --------------

func TestEmptyProgram(t *testing.T) {
	p := newParserFromInput("")
	program, _ := p.ParseProgram()

	assertNoErrors(t, p)
	node.AssertProgram(t, program, false)
}

func TestProgramWithNewlines(t *testing.T) {
	p := newParserFromInput("\n \n  \n \n")
	_, _ = p.ParseProgram()

	assertNoErrors(t, p)
}

func TestInvalidProgram(t *testing.T) {
	p := newParserFromInput("hello")
	_, _ = p.ParseProgram()

	assertHasError(t, p)
}

// Contract Nodes
// --------------

func TestEmptyContract(t *testing.T) {
	p := newParserFromInput("contract Test {\n \n}")
	c := p.parseContract()

	assertNoErrors(t, p)
	node.AssertContract(t, c, "Test", 0, 0)
}

func TestContractWithVariable(t *testing.T) {
	p := newParserFromInput(`contract Test {
		int x
		int y
	}`)
	c := p.parseContract()

	assertNoErrors(t, p)
	node.AssertContract(t, c, "Test", 2, 0)
	node.AssertVariable(t, c.Variables[0], "int", "x")
	node.AssertVariable(t, c.Variables[1], "int", "y")
}

// Variable Nodes
// --------------

func TestVariable(t *testing.T) {
	p := newParserFromInput("int x \n")
	v := p.parseVariableStatement()

	node.AssertVariable(t, v, "int", "x")
	assertNoErrors(t, p)
}

func TestVariableWithoutNewLine(t *testing.T) {
	p := newParserFromInput("int x")
	_ = p.parseVariableStatement()
	assertHasError(t, p)
}

// Function Nodes
// --------------
func TestEmptyFunction(t *testing.T) {
	p := newParserFromInput("function void test(){\n}\n")
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 0, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithParam(t *testing.T) {
	p := newParserFromInput("function void test(int a){\n}\n")
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 1, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithParams(t *testing.T) {
	p := newParserFromInput("function void test(int a, int b){\n}\n")
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 2, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithMultipleRTypes(t *testing.T) {
	p := newParserFromInput("function (int, int) test(){\n}\n")
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 2, 0, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithParamsAndRTypes(t *testing.T) {
	p := newParserFromInput("function (int, int) test(int a, int b){\n}\n")
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 2, 2, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithStatement(t *testing.T) {
	p := newParserFromInput("function void test(){\nint a\n}\n")
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 0, 1)
	assertNoErrors(t, p)
}

func TestFunctionWithMultipleStatements(t *testing.T) {
	p := newParserFromInput("function void test(){\nint a\nint b\n}\n")
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 0, 2)
	assertNoErrors(t, p)
}

func TestFullFunction(t *testing.T) {
	p := newParserFromInput("function (int, int) test(int a, int b){\nint a\nint b\n}\n")
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 2, 2, 2)
	assertNoErrors(t, p)
}

func TestFunctionWORType(t *testing.T) {
	p := newParserFromInput("function test(int a, int b){\n}\n")
	p.parseFunction()
	assertHasError(t, p)
}

func TestFunctionMissingNewline(t *testing.T) {
	p := newParserFromInput("function test(int a, int b){\n}")
	p.parseFunction()
	assertHasError(t, p)
}

func TestFunctionMissingNewlineInBody(t *testing.T) {
	p := newParserFromInput("function test(int a, int b){}\n")
	p.parseFunction()
	assertHasError(t, p)
}

// Statement Nodes
// ---------------

func TestEmptyStatementBlock(t *testing.T) {
	p := newParserFromInput("{\n}\n")
	v := p.parseStatementBlock()

	node.AssertStatementBlock(t, v, 0)
	assertNoErrors(t, p)
}

func TestStatementBlock(t *testing.T) {
	p := newParserFromInput("{\nint a = 5\n}\n")
	v := p.parseStatementBlock()

	node.AssertStatementBlock(t, v, 1)
	assertNoErrors(t, p)
}

func TestMultipleStatementBlock(t *testing.T) {
	p := newParserFromInput("{\nint a = 5\nint b = 4\n}\n")
	v := p.parseStatementBlock()

	node.AssertStatementBlock(t, v, 2)
	assertNoErrors(t, p)
}

func TestReturnStatementMissingNewline(t *testing.T) {
	p := newParserFromInput("return")
	p.parseReturnStatement()
	assertHasError(t, p)
}

func TestEmptyReturnStatement(t *testing.T) {
	p := newParserFromInput("return \n")
	v := p.parseReturnStatement()

	node.AssertReturnStatement(t, v, 0)
	assertNoErrors(t, p)
}

func TestSingleReturnStatement(t *testing.T) {
	p := newParserFromInput("return 1\n")
	v := p.parseReturnStatement()

	node.AssertReturnStatement(t, v, 1)
	assertNoErrors(t, p)
}

func TestMultipleReturnStatement(t *testing.T) {
	p := newParserFromInput("return 1, 2\n")
	v := p.parseReturnStatement()

	node.AssertReturnStatement(t, v, 2)
	assertNoErrors(t, p)
}

func TestIfStatement(t *testing.T) {
	p := newParserFromInput("if(true){\n} else{\n}\n")
	v := p.parseIfStatement()

	node.AssertIfStatement(t, v, "true", 0, 0)
	assertNoErrors(t, p)
}

func TestIfStatementSingleStatement(t *testing.T) {
	p := newParserFromInput("if(true){\nint a \n} else{\nint b\n}\n")
	v := p.parseIfStatement()

	node.AssertIfStatement(t, v, "true", 1, 1)
	assertNoErrors(t, p)
}

func TestIfStatementSingleThenStatement(t *testing.T) {
	p := newParserFromInput("if(true){\nint a \n} else{\n}\n")
	v := p.parseIfStatement()

	node.AssertIfStatement(t, v, "true", 1, 0)
	assertNoErrors(t, p)
}

func TestIfStatementSingleElseStatement(t *testing.T) {
	p := newParserFromInput("if(true){\n} else{\n int a \n}\n")
	v := p.parseIfStatement()

	node.AssertIfStatement(t, v, "true", 0, 1)
	assertNoErrors(t, p)
}

func TestIfStatementMultipleStatement(t *testing.T) {
	p := newParserFromInput("if(true){\nint a\n int b\n} else{\nint c\n int d\n}\n")
	v := p.parseIfStatement()

	node.AssertIfStatement(t, v, "true", 2, 2)
	assertNoErrors(t, p)
}

func TestIfStatementMultipleThenStatement(t *testing.T) {
	p := newParserFromInput("if(true){\nint a\n int b\n} else{\n}\n")
	v := p.parseIfStatement()

	node.AssertIfStatement(t, v, "true", 2, 0)
	assertNoErrors(t, p)
}

func TestIfStatementMultipleElseStatement(t *testing.T) {
	p := newParserFromInput("if(true){\n} else{\nint c\n int d\n}\n")
	v := p.parseIfStatement()

	node.AssertIfStatement(t, v, "true", 0, 2)
	assertNoErrors(t, p)
}

func TestIfStatementWOElse(t *testing.T) {
	p := newParserFromInput("if(true){\n}\n")
	v := p.parseIfStatement()

	node.AssertIfStatement(t, v, "true", 0, 0)
	assertNoErrors(t, p)
}

func TestIfStatementWOElseWONewline(t *testing.T) {
	p := newParserFromInput("if(true){\n}")
	p.parseIfStatement()

	assertHasError(t, p)
}

func TestAssignmentStatement(t *testing.T) {
	p := newParserFromInput("a = 5\n")
	i := p.readIdentifier()
	v := p.parseAssignmentStatement(i)

	node.AssertAssignmentStatement(t, v, "a", "5")
	assertNoErrors(t, p)
}

func TestAssignmentStatementChar(t *testing.T) {
	p := newParserFromInput("a = 'c'\n")
	i := p.readIdentifier()
	v := p.parseAssignmentStatement(i)

	node.AssertAssignmentStatement(t, v, "a", "c")
	assertNoErrors(t, p)
}

func TestAssignmentStatementWONewline(t *testing.T) {
	p := newParserFromInput("a = 'c'")
	i := p.readIdentifier()
	p.parseAssignmentStatement(i)

	assertHasError(t, p)
}

// Designator Nodes
// ----------------

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

// Literal Nodes
// -------------

func TestIntegerLiteral(t *testing.T) {
	p := newParserFromInput("1")
	i := p.parseInteger()
	node.AssertIntegerLiteral(t, i, big.NewInt(1))
	assertNoErrors(t, p)
}

func TestInvalidIntegerLiteral(t *testing.T) {
	p := newParserFromInput("0x1")
	i := p.parseInteger()
	node.AssertIntegerLiteral(t, i, big.NewInt(1))
	assertNoErrors(t, p)
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
	tok, _ := b.(*node.ErrorNode)
	node.AssertError(t, tok, "Invalid boolean value if")
	assertHasError(t, p)
}

// Type Nodes
//-----------

func TestTypeNode(t *testing.T) {
	p := newParserFromInput("int")
	v := p.parseType()
	node.AssertType(t, v, "int")
}

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
