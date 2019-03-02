package lexer

import (
	"gotest.tools/assert"
	"math/big"
	"testing"
)

func TestEmptySource(t *testing.T) {
	tester := newLexerTesterFromInput(t, "")

	tester.assertTotal(0)
	assert.Equal(t, tester.lex.IsEnd, true)
	assert.Equal(t, tester.lex.current, int32(0))
	assert.Equal(t, tester.lex.currentPos.String(), "1:0")
}

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

func TestMixedIntegers(t *testing.T) {
	tester := newLexerTesterFromInput(t, "123abc 0XAFG")

	tester.assertTotal(4)
	tester.assertInteger(0, big.NewInt(123))
	tester.assertIdentifer(1, "abc")
	tester.assertInteger(2, big.NewInt(0xaf))
	tester.assertIdentifer(3, "G")
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

// Character Tokens
// ----------------
func TestCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'c'")
	tester.assertCharacter(0, 'c')
}

func TestMultipleCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'c''b'")
	tester.assertCharacter(0, 'c')
	tester.assertCharacter(1, 'b')
}

func TestBackslashCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'\\\\'")
	tester.assertCharacter(0, '\\')
}

func TestQuoteCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'\\''")
	tester.assertCharacter(0, '\'')
}

func TestNewlineCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'\n'")
	tester.assertCharacter(0, '\n')
}

func TestInvalidCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'cc'")
	tester.assertError(0, "cc")
}

// String Tokens
// -------------

// Fix Tokens
// ----------------
