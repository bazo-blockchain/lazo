package lexer

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"gotest.tools/assert"
	"math/big"
	"strings"
	"testing"
)

type lexerTestUtil struct {
	t      *testing.T
	lex    *Lexer
	tokens []token.Token
}

func newLexerTestUtil(t *testing.T, input string) *lexerTestUtil {
	tester := &lexerTestUtil{
		t:   t,
		lex: New(bufio.NewReader(strings.NewReader(input))),
	}

	for !tester.lex.isEnd {
		tok := tester.lex.NextToken()
		if ftok, ok := tok.(*token.FixToken); !ok || ftok.Value != token.NewLine {
			tester.tokens = append(tester.tokens, tok)
		}
	}
	return tester
}

func (tester *lexerTestUtil) assertTotal(total int) {
	assert.Equal(tester.t, len(tester.tokens), total)
}

func (tester *lexerTestUtil) assertInteger(index int, value *big.Int) {
	token.AssertInteger(tester.t, tester.tokens[index], value)
}

func (tester *lexerTestUtil) assertIdentifer(index int, value string) {
	token.AssertIdentifier(tester.t, tester.tokens[index], value)
}

func (tester *lexerTestUtil) assertString(index int, value string) {
	token.AssertString(tester.t, tester.tokens[index], value)
}

func (tester *lexerTestUtil) assertCharacter(index int, value rune) {
	token.AssertCharacter(tester.t, tester.tokens[index], value)
}

func (tester *lexerTestUtil) assertFixToken(index int, value token.Symbol) {
	token.AssertFixToken(tester.t, tester.tokens[index], value)
}

func (tester *lexerTestUtil) assertError(index int, value string) {
	token.AssertError(tester.t, tester.tokens[index], value)
}
