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
	tok, ok := tester.tokens[index].(*token.IntegerToken)

	assert.Equal(tester.t, ok, true)
	assert.Equal(tester.t, tok.Value.Cmp(value) == 0, true)
}

func (tester *lexerTestUtil) assertIdentifer(index int, value string) {
	assertIdentifier(tester.t, tester.tokens[index], value)
}

func (tester *lexerTestUtil) assertString(index int, value string) {
	tok, ok := tester.tokens[index].(*token.StringToken)

	assert.Equal(tester.t, ok, true)
	assert.Equal(tester.t, tok.Literal(), value)
}

func (tester *lexerTestUtil) assertCharacter(index int, value rune) {
	tok, ok := tester.tokens[index].(*token.CharacterToken)

	assert.Equal(tester.t, ok, true)
	assert.Equal(tester.t, tok.Value, value)
}

func (tester *lexerTestUtil) assertFixToken(index int, value token.Symbol) {
	assertFixToken(tester.t, tester.tokens[index], value)
}

func (tester *lexerTestUtil) assertError(index int, value string) {
	tok, ok := tester.tokens[index].(*token.ErrorToken)

	assert.Equal(tester.t, ok, true)
	assert.Equal(tester.t, tok.Literal(), value)
}

func assertIdentifier(t *testing.T, tok token.Token, value string) {
	tok, ok := tok.(*token.IdentifierToken)

	assert.Equal(t, ok, true)
	assert.Equal(t, tok.Literal(), value)
}

func assertFixToken(t *testing.T, tok token.Token, value token.Symbol) {
	ftok, ok := tok.(*token.FixToken)

	assert.Equal(t, ok, true)
	assert.Equal(t, ftok.Value, value)
}
