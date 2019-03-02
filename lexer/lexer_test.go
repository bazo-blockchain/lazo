package lexer

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"gotest.tools/assert"
	"math/big"
	"strings"
	"testing"
)

func TestStateWithoutCode(t *testing.T) {
	lex := New(bufio.NewReader(strings.NewReader("")))

	// initial state without reading the token
	assertLexerState(t, lex, true, 0, "1:0")

	tok := lex.NextToken()
	token.AssertFixToken(t, tok, token.EOF)

	// it shouldn't have changed the initial lexer state
	assertLexerState(t, lex, true, 0, "1:0")
}

func TestStateWithCode(t *testing.T) {
	lex := New(bufio.NewReader(strings.NewReader("test")))

	// initial state without reading the token
	assertLexerState(t, lex, false, 't', "1:1")

	tok := lex.NextToken()
	token.AssertIdentifier(t, tok, "test")
	assertLexerState(t, lex, true, 0, "1:4")

	tok = lex.NextToken()
	token.AssertFixToken(t, tok, token.EOF)
	assertLexerState(t, lex, true, 0, "1:4")
}

func assertLexerState(t *testing.T, lex *Lexer, isEnd bool, current rune, pos string) {
	assert.Equal(t, lex.IsEnd, isEnd)
	assert.Equal(t, lex.current, current)
	assert.Equal(t, lex.currentPos.String(), pos)
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


// String Tokens
// -------------

func TestEmptyString(t *testing.T) {
	tester := newLexerTesterFromInput(t, `""`)

	tester.assertTotal(1)
	tester.assertString(0, "")
}

func TestStrings(t *testing.T) {
	tester := newLexerTesterFromInput(t, `"test" "test2"`)

	tester.assertTotal(2)
	tester.assertString(0, "test")
	tester.assertString(1, "test2")
}

func TestStringWithEscapedChars(t *testing.T) {
	tester := newLexerTesterFromInput(t, `
		"new \n line"
		"double quote \""
		"backslash \\"`)

	tester.assertTotal(3)
	tester.assertString(0, "new \n line")
	tester.assertString(1, "double quote \"")
	tester.assertString(2, "backslash \\")
}

func TestStringWithNotAllowedEscapeChars(t *testing.T) {
	tester := newLexerTesterFromInput(t, `"single quote \' "`)

	tester.assertTotal(2)
	tester.assertError(0, "single quote ")
	tester.assertError(1, "")
}

func TestNotClosedString(t *testing.T) {
	tester := newLexerTesterFromInput(t, `"not closed`)

	tester.assertTotal(1)
	tester.assertError(0, "not closed")
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

func TestSingleQuoteCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'\\''")
	tester.assertCharacter(0, '\'')
}

func TestNewlineCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'\n'")
	tester.assertCharacter(0, '\n')
}

func TestEmptyCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'\\0'")
	tester.assertCharacter(0, 0)
}

func TestCharacterWithInvalidEscape(t *testing.T) {
	tester := newLexerTesterFromInput(t, `'\"'`)
	tester.assertError(0, "")
}

func TestInvalidCharacter(t *testing.T) {
	tester := newLexerTesterFromInput(t, "'cc'")
	tester.assertError(0, "cc")
}

// Fix Tokens
// ----------------

func TestReservedKeyword(t *testing.T) {
	tester := newLexerTesterFromInput(t, "contract Contract if IF")

	tester.assertTotal(4)
	tester.assertFixToken(0, token.Contract)
	tester.assertIdentifer(1, "Contract")
	tester.assertFixToken(2, token.If)
	tester.assertIdentifer(3, "IF")
}