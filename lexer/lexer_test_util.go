package lexer

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"gotest.tools/assert"
	"math/big"
	"strings"
	"testing"
)

type lexerTester struct {
	t      *testing.T
	lex    *Lexer
	tokens []token.Token
}

func newLexerTesterFromInput(t *testing.T, input string) *lexerTester {
	tester := &lexerTester{
		t: t,
		lex: New(bufio.NewReader(strings.NewReader(input))),
	}

	for !tester.lex.EOF {
		tester.tokens = append(tester.tokens, tester.lex.NextToken())
	}
	return tester
}

func (tester *lexerTester) assertTotal(total int) {
	assert.Equal(tester.t, len(tester.tokens), total)
}

func (tester *lexerTester) assertInteger(index int, value *big.Int){
	token.AssertInteger(tester.t, tester.tokens[index], value)
}

func (tester *lexerTester) assertIdentifer(index int, value string) {
	token.AssertIdentifier(tester.t, tester.tokens[index], value)
}

func (tester *lexerTester) assertError(index int, value string) {
	token.AssertError(tester.t, tester.tokens[index], value)
}
