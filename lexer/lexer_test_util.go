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

func (tester *lexerTester) assertInteger(tokenIndex int, value *big.Int){
	tok, ok := tester.tokens[tokenIndex].(*token.IntegerToken)

	assert.Equal(tester.t, ok, true)
	assert.Equal(tester.t, tok.Value.Cmp(value) == 0, true)
}

func (tester *lexerTester) assertIdentifer(tokenIndex int, value string) {
	tok, ok := tester.tokens[tokenIndex].(*token.IdentifierToken)

	assert.Equal(tester.t, ok, true)
	assert.Equal(tester.t, tok.Literal(), value)
}
