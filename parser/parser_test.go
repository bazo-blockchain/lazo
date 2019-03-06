package parser

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"strings"
	"testing"
)

// Variable Nodes
// --------------

func TestVariable(t *testing.T) {
	p := newParserFromInput("int x \n")
	v := p.parseVariable()

	node.AssertVariable(t, v, "int", "x")
}

func TestVariableWithoutNewLine(t *testing.T) {
	p := newParserFromInput("int x")
	_ = p.parseVariable()
	assertHasError(t, p)
}

// ------------------------------

func newParserFromInput(input string) *Parser {
	return New(lexer.New(bufio.NewReader(strings.NewReader(input))))
}

func assertHasError(t *testing.T, p *Parser) {
	assert.Equal(t, len(p.errors) > 0, true)
}
