package parser

import (
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"testing"
)

// Binary Expressions
// ------------------------

func TestLogicOr(t *testing.T) {
	p := newParserFromInput("x || y || z")
	e := p.parseOr()

	node.AssertBinaryExpression(t, e, "(x || y)", "z", token.Or)
}

func TestLogicAnd(t *testing.T) {
	p := newParserFromInput("x && y && z")
	e := p.parseAnd()

	node.AssertBinaryExpression(t, e, "(x && y)", "z", token.And)
}

func TestMixedLogicOperatorsAndOr(t *testing.T) {
	p := newParserFromInput("x && y || z")
	e := p.parseExpression()

	node.AssertBinaryExpression(t, e, "(x && y)", "z", token.Or)
}

func TestMixedLogicOperatorOrAnd(t *testing.T) {
	p := newParserFromInput("x || y && z")
	e := p.parseExpression()

	node.AssertBinaryExpression(t, e, "x", "(y && z)", token.Or)
}
