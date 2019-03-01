package lexer

import (
	"math/big"
	"testing"
)

// Integer Tokens
// --------------

func TestDecimalDigits(t *testing.T) {
	tester := newLexerTesterFromInput(t, "0 001 456")

	tester.assertTotal(3)
	tester.assertInteger(0, big.NewInt(0))
	tester.assertInteger(1, big.NewInt(1))
	tester.assertInteger(2, big.NewInt(456))
}

func TestHexDigits(t *testing.T) {
	tester := newLexerTesterFromInput(t, "0x123 0xaf 0x123af")

	tester.assertTotal(3)
	tester.assertInteger(0, big.NewInt(0x123))
	tester.assertInteger(0, big.NewInt(291))
	tester.assertInteger(1, big.NewInt(175))
	tester.assertInteger(2, big.NewInt(74671))
}

func TestInvalidHex(t *testing.T) {
	tester := newLexerTesterFromInput(t, "0xg")

	tester.assertTotal(2)
	tester.assertError(0, "0x")
	tester.assertIdentifer(1, "g")
}

// Identifier Tokens
// -----------------

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
