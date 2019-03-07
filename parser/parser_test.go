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
	node.AssertProgram(t, program, false)
	assertNoErrors(t, p)
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

	node.AssertContract(t, c, "Test", 0, 0)
	assertNoErrors(t, p)
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
