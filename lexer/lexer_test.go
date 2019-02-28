package lexer

import (
	"testing"
)

// Integer Token
// -------------

func TestDecimalDigits(t *testing.T) {
	input := "0 123 456"

	lex := createLexerFromInput(input)
	assert.Equal(t, lex.NextToken().Literal(), "0")
}
