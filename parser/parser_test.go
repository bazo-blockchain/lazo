package parser

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
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
}

func TestVariableWithoutNewLine(t *testing.T) {
	p := newParserFromInput("int x")
	_ = p.parseVariableStatement()
	assertHasError(t, p)
}

// Function Nodes
// --------------
func TestEmptyFunction(t *testing.T) {
	p := newParserFromInput(`function void test(){}`)
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 0, 0)
}

func TestFunctionWithParam(t *testing.T) {
	p := newParserFromInput(`function void test(int a){}`)
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 1, 0)
}

func TestFunctionWithParams(t *testing.T) {
	p := newParserFromInput(`function void test(int a, int b){}`)
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 2, 0)
}

func TestFunctionWithMultipleRTypes(t *testing.T) {
	p := newParserFromInput(`function (int, int) test(){}`)
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 2, 0, 0)
}

func TestFunctionWithParamsAndRTypes(t *testing.T) {
	p := newParserFromInput(`function (int, int) test(int a, int b){}`)
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 2, 2, 0)
}

func TestFunctionWithStatement(t *testing.T) {
	p := newParserFromInput(`function void test(){
		int a
		}`)
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 0, 1)
}

func TestFunctionWithMultipleStatements(t *testing.T) {
	p := newParserFromInput(`function void test(){
		int a
		int b
		}`)
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 1, 0, 2)
}

func TestFullFunction(t *testing.T) {
	p := newParserFromInput(`function (int, int) test(int a, int b){
		int a
		int b
		}`)
	f := p.parseFunction()
	node.AssertFunction(t, f, "test", 2, 2, 2)
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
