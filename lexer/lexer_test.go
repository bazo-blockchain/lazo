package lexer

import (
	"math/big"
	"testing"
)

// Integer Tokens
// --------------

func TestDecimalDigits(t *testing.T) {
	tester := newLexerTesterFromInput(t, "0 123 456")

	tester.assertTotal(3)
	tester.assertInteger(0, big.NewInt(0))
	tester.assertInteger(1, big.NewInt(123))
	tester.assertInteger(2, big.NewInt(456))
}

// Integer Tokens
// --------------

func TestValidIdentifers(t *testing.T) {
	tester := newLexerTesterFromInput(t, "id test three3 under_score4 _5five")
	expectedIds := []string{
		"id",
		"test",
		"three3",
		"under_score4",
		"_5five",
	}

	tester.assertTotal(len(expectedIds))
	for i, id := range expectedIds {
		tester.assertIdentifer(i, id)
	}
}
