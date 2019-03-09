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
	assert.Equal(t, tok.Pos().String(), "1:0")

	// it shouldn't have changed the initial lexer state
	assertLexerState(t, lex, true, 0, "1:0")
}

func TestStateWithCode(t *testing.T) {
	lex := New(bufio.NewReader(strings.NewReader("test")))

	// initial state without reading the token
	assertLexerState(t, lex, false, 't', "1:1")

	tok := lex.NextToken()
	token.AssertIdentifier(t, tok, "test")
	assert.Equal(t, tok.Pos().String(), "1:1")
	assertLexerState(t, lex, true, 0, "1:4")

	tok = lex.NextToken()
	token.AssertFixToken(t, tok, token.EOF)
	assert.Equal(t, tok.Pos().String(), "1:4")
	assertLexerState(t, lex, true, 0, "1:4")
}

func assertLexerState(t *testing.T, lex *Lexer, isEnd bool, current rune, pos string) {
	assert.Equal(t, lex.isEnd, isEnd)
	assert.Equal(t, lex.current, current)
	assert.Equal(t, lex.currentPos.String(), pos)
}

// Integer Tokens
// --------------

func TestDecimalDigits(t *testing.T) {
	tester := newLexerTestUtil(t, "0 001 456")

	tester.assertTotal(3)
	tester.assertInteger(0, big.NewInt(0))
	tester.assertInteger(1, big.NewInt(1))
	tester.assertInteger(2, big.NewInt(456))
}

func TestHexDigits(t *testing.T) {
	tester := newLexerTestUtil(t, "0x123 0xaf 0x123af")

	tester.assertTotal(3)
	tester.assertInteger(0, big.NewInt(0x123))
	tester.assertInteger(0, big.NewInt(291))
	tester.assertInteger(1, big.NewInt(175))
	tester.assertInteger(2, big.NewInt(74671))
}

func TestInvalidHex(t *testing.T) {
	tester := newLexerTestUtil(t, "0xg")

	tester.assertTotal(2)
	tester.assertError(0, "0x")
	tester.assertIdentifer(1, "g")
}

func TestMixedIntegers(t *testing.T) {
	tester := newLexerTestUtil(t, "123abc 0XAFG")

	tester.assertTotal(4)
	tester.assertInteger(0, big.NewInt(123))
	tester.assertIdentifer(1, "abc")
	tester.assertInteger(2, big.NewInt(0xaf))
	tester.assertIdentifer(3, "G")
}

// Identifier Tokens
// -----------------

func TestValidIdentifers(t *testing.T) {
	tester := newLexerTestUtil(t, "id test three3 under_score4 _5five")
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
	tester := newLexerTestUtil(t, `""`)

	tester.assertTotal(1)
	tester.assertString(0, "")
}

func TestStrings(t *testing.T) {
	tester := newLexerTestUtil(t, `"test" "test2"`)

	tester.assertTotal(2)
	tester.assertString(0, "test")
	tester.assertString(1, "test2")
}

func TestStringWithEscapedChars(t *testing.T) {
	tester := newLexerTestUtil(t, `
		"new \n line"
		"double quote \""
		"backslash \\"`)

	tester.assertTotal(3)
	tester.assertString(0, "new \n line")
	tester.assertString(1, "double quote \"")
	tester.assertString(2, "backslash \\")
}

func TestStringWithNotAllowedEscapeChars(t *testing.T) {
	tester := newLexerTestUtil(t, `"single quote \' "`)

	tester.assertTotal(2)
	tester.assertError(0, "single quote ")
	tester.assertError(1, "")
}

func TestNotClosedString(t *testing.T) {
	tester := newLexerTestUtil(t, `"not closed`)

	tester.assertTotal(1)
	tester.assertError(0, "not closed")
}

// Character Tokens
// ----------------
func TestCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'c'")
	tester.assertCharacter(0, 'c')
}

func TestMultipleCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'c''b'")
	tester.assertCharacter(0, 'c')
	tester.assertCharacter(1, 'b')
}

func TestBackslashCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'\\\\'")
	tester.assertCharacter(0, '\\')
}

func TestSingleQuoteCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'\\''")
	tester.assertCharacter(0, '\'')
}

func TestNewlineCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'\n'")
	tester.assertCharacter(0, '\n')
}

func TestEmptyCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'\\0'")
	tester.assertCharacter(0, 0)
}

func TestInvalidEmptyCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "''")
	tester.assertError(0, "")
}

func TestCharacterWithInvalidEscape(t *testing.T) {
	tester := newLexerTestUtil(t, `'\"'`)
	tester.assertError(0, "")
}

func TestInvalidCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'cc'")
	tester.assertError(0, "cc")
}

func TestUnicodeCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'£'")
	tester.assertError(0, "£")
}

func TestNotClosedCharacter(t *testing.T) {
	tester := newLexerTestUtil(t, "'c")
	tester.assertError(0, "c")
}

// Fix Tokens
// ----------------

func TestReservedKeyword(t *testing.T) {
	tester := newLexerTestUtil(t, "contract Contract if IF")

	tester.assertTotal(4)
	tester.assertFixToken(0, token.Contract)
	tester.assertIdentifer(1, "Contract")
	tester.assertFixToken(2, token.If)
	tester.assertIdentifer(3, "IF")
}

func TestContract(t *testing.T) {
	tester := newLexerTestUtil(t, "contract")
	tester.assertFixToken(0, token.Contract)
}

func TestReturn(t *testing.T) {
	tester := newLexerTestUtil(t, "return")
	tester.assertFixToken(0, token.Return)
}

func TestIf(t *testing.T) {
	tester := newLexerTestUtil(t, "if")
	tester.assertFixToken(0, token.If)
}

func TestElseIf(t *testing.T) {
	tester := newLexerTestUtil(t, "else if")
	tester.assertFixToken(0, token.Else)
	tester.assertFixToken(1, token.If)
}

func TestElse(t *testing.T) {
	tester := newLexerTestUtil(t, "else")
	tester.assertFixToken(0, token.Else)
}

func TestFunction(t *testing.T) {
	tester := newLexerTestUtil(t, "function")
	tester.assertFixToken(0, token.Function)
}

func TestColon(t *testing.T) {
	tester := newLexerTestUtil(t, ":")
	tester.assertFixToken(0, token.Colon)
}

func TestComma(t *testing.T) {
	tester := newLexerTestUtil(t, ",")
	tester.assertFixToken(0, token.Comma)
	assert.Equal(t, tester.tokens[0].(*token.FixToken).Lexeme, ",")
}

func TestPeriod(t *testing.T) {
	tester := newLexerTestUtil(t, ".")
	tester.assertFixToken(0, token.Period)
}

func TestOpenBrace(t *testing.T) {
	tester := newLexerTestUtil(t, "{")
	tester.assertFixToken(0, token.OpenBrace)
}

func TestCloseBrace(t *testing.T) {
	tester := newLexerTestUtil(t, "}")
	tester.assertFixToken(0, token.CloseBrace)
}

func TestOpenBracket(t *testing.T) {
	tester := newLexerTestUtil(t, "[")
	tester.assertFixToken(0, token.OpenBracket)
}

func TestCloseBracket(t *testing.T) {
	tester := newLexerTestUtil(t, "]")
	tester.assertFixToken(0, token.CloseBracket)
}

func TestOpenParen(t *testing.T) {
	tester := newLexerTestUtil(t, "(")
	tester.assertFixToken(0, token.OpenParen)
}

func TestCloseParen(t *testing.T) {
	tester := newLexerTestUtil(t, ")")
	tester.assertFixToken(0, token.CloseParen)
}

func TestAddition(t *testing.T) {
	tester := newLexerTestUtil(t, "+")
	tester.assertFixToken(0, token.Addition)
}

func TestSubtraction(t *testing.T) {
	tester := newLexerTestUtil(t, "-")
	tester.assertFixToken(0, token.Subtraction)
}

func TestMultiplication(t *testing.T) {
	tester := newLexerTestUtil(t, "*")
	tester.assertFixToken(0, token.Multiplication)
}

func TestDivision(t *testing.T) {
	tester := newLexerTestUtil(t, "/")
	tester.assertFixToken(0, token.Division)
}

func TestModulo(t *testing.T) {
	tester := newLexerTestUtil(t, "%")
	tester.assertFixToken(0, token.Modulo)
}

func TestAssign(t *testing.T) {
	tester := newLexerTestUtil(t, "=")
	tester.assertFixToken(0, token.Assign)
}

func TestGreater(t *testing.T) {
	tester := newLexerTestUtil(t, ">")
	tester.assertFixToken(0, token.Greater)
}

func TestLess(t *testing.T) {
	tester := newLexerTestUtil(t, "<")
	tester.assertFixToken(0, token.Less)
}

func TestNot(t *testing.T) {
	tester := newLexerTestUtil(t, "!")
	tester.assertFixToken(0, token.Not)
}

func TestAnd(t *testing.T) {
	tester := newLexerTestUtil(t, "&&")
	tester.assertFixToken(0, token.And)
}

func TestOr(t *testing.T) {
	tester := newLexerTestUtil(t, "||")
	tester.assertFixToken(0, token.Or)
}

func TestEqual(t *testing.T) {
	tester := newLexerTestUtil(t, "==")
	tester.assertFixToken(0, token.Equal)
}

func TestUnequal(t *testing.T) {
	tester := newLexerTestUtil(t, "!=")
	tester.assertFixToken(0, token.Unequal)
}

func TestGreaterEqual(t *testing.T) {
	tester := newLexerTestUtil(t, ">=")
	tester.assertFixToken(0, token.GreaterEqual)
}

func TestLessEqual(t *testing.T) {
	tester := newLexerTestUtil(t, "<=")
	tester.assertFixToken(0, token.LessEqual)
}

func TestInvalidFixToken(t *testing.T) {
	tester := newLexerTestUtil(t, "$")
	tester.assertError(0, "$")
}

func TestFuncDeclaration(t *testing.T) {
	tester := newLexerTestUtil(t, "function void test() {}")
	tester.assertFixToken(0, token.Function)
	tester.assertIdentifer(1, "void")
	tester.assertIdentifer(2, "test")
	tester.assertFixToken(3, token.OpenParen)
	tester.assertFixToken(4, token.CloseParen)
	tester.assertFixToken(5, token.OpenBrace)
	tester.assertFixToken(6, token.CloseBrace)
}

func TestLinuxNewLine(t *testing.T) {
	lex := New(bufio.NewReader(strings.NewReader("\n \n")))

	tok := lex.NextToken()
	token.AssertFixToken(t, tok, token.NewLine)
	tok = lex.NextToken()
	token.AssertFixToken(t, tok, token.NewLine)
}

func TestWindowsNewLine(t *testing.T) {
	lex := New(bufio.NewReader(strings.NewReader("\r\n")))

	tok := lex.NextToken()
	token.AssertFixToken(t, tok, token.NewLine)
	tok = lex.NextToken()
	token.AssertFixToken(t, tok, token.EOF)
}

func TestSystemNewLine(t *testing.T) {
	lex := New(bufio.NewReader(strings.NewReader(`
		1
		2
	`)))

	tok := lex.NextToken()
	token.AssertFixToken(t, tok, token.NewLine)
	tok = lex.NextToken() // skip 1
	tok = lex.NextToken()
	token.AssertFixToken(t, tok, token.NewLine)
	tok = lex.NextToken() // skip 2
	tok = lex.NextToken()
	token.AssertFixToken(t, tok, token.NewLine)
}
