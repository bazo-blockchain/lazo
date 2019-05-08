// Package lexer performs lexical analysis and creates tokens.
// It reads the input source code character by character, recognizes the lexemes
// and outputs a sequence of tokens describing the lexemes.
package lexer

import (
	"bufio"
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/big"
)

// Lexer holds the current character and position from the given reader.
// It is used to scan characters and create tokens.
type Lexer struct {
	reader     *bufio.Reader
	current    rune
	currentPos token.Position
	tokenPos   token.Position
	isEnd      bool
}

// New creates a new Lexer struct with the given reader and initializes the current position.
// It also reads the first character and initializes the current character.
// It returns the created lexer struct
func New(reader *bufio.Reader) *Lexer {
	lex := &Lexer{
		reader:     reader,
		currentPos: token.NewPosition(),
	}
	lex.nextChar()
	return lex
}

// NextToken reads character by character from reader and creates a token when possible.
// White space is skipped and, therefore, no token is created for that.
// However, tokens are created for new lines, since they are part of the syntax.
// It returns the created token containing the token position (line and column), the literal itself and the token type.
func (lex *Lexer) NextToken() token.Token {
	lex.skipWhiteSpace()

	lex.tokenPos = lex.currentPos

	if lex.isEnd {
		return &token.FixToken{
			AbstractToken: lex.newAbstractToken(""),
			Value:         token.EOF,
		}
	}

	if lex.isDigit() {
		return lex.readInteger()
	}

	if lex.isLetter() || lex.isChar('_') {
		return lex.readName()
	}

	switch lex.current {
	case '"':
		return lex.readString()
	case '\'':
		return lex.readCharacter()
	default:
		return lex.readFixToken()
	}
}

func (lex *Lexer) skipWhiteSpace() {
	for !lex.isEnd && lex.current <= ' ' && !lex.isChar('\n') {
		lex.nextChar()
	}
}

func (lex *Lexer) readInteger() token.Token {
	var lexeme string
	value := new(big.Int)
	var ok bool
	var isHex bool

	peekChar, peekError := lex.peekChar()
	if lex.isChar('0') &&
		(peekChar == 'x' || peekChar == 'X') && peekError == nil {
		// skip 0x
		lex.nextChar()
		lex.nextChar()

		lexeme = lex.readLexeme(lex.isHexDigit)
		value, ok = value.SetString(lexeme, 16)
		isHex = true

		lexeme = "0x" + lexeme
	} else {
		lexeme = lex.readLexeme(lex.isDigit)
		value, ok = value.SetString(lexeme, 10)
	}

	abstractToken := lex.newAbstractToken(lexeme)
	if !ok {
		return lex.newErrorToken(abstractToken, "Error while parsing string to big int")
	}

	return &token.IntegerToken{
		AbstractToken: abstractToken,
		Value:         value,
		IsHex:         isHex,
	}
}

func (lex *Lexer) readName() token.Token {
	lexeme := lex.readLexeme(func() bool {
		return lex.isLetter() || lex.isChar('_') || lex.isDigit()
	})
	abstractToken := lex.newAbstractToken(lexeme)

	if symbol, ok := token.Keywords[lexeme]; ok {
		return &token.FixToken{
			AbstractToken: abstractToken,
			Value:         symbol,
		}
	}

	return &token.IdentifierToken{
		AbstractToken: abstractToken,
	}
}

func (lex *Lexer) readString() token.Token {
	// skip opening double quote
	lex.nextChar()

	lexeme, err := lex.readEscapedLexeme(func() bool {
		return !lex.isChar('"')
	}, allowedStringEscapedCodes)

	abstractToken := lex.newAbstractToken(lexeme)
	if err != nil {
		return lex.newErrorToken(abstractToken, err.Error())
	}

	if lex.isChar('"') {
		// skip closing double quote
		lex.nextChar()

		return &token.StringToken{
			AbstractToken: abstractToken,
		}
	}

	return lex.newErrorToken(abstractToken, "String not closed")
}

func (lex *Lexer) readCharacter() token.Token {
	// skip opening quote
	lex.nextChar()

	lexeme, err := lex.readEscapedLexeme(func() bool {
		return !lex.isChar('\'')
	}, allowedCharEscapedCodes)

	abstractToken := lex.newAbstractToken(lexeme)
	if err != nil {
		return lex.newErrorToken(abstractToken, err.Error())
	}

	if len(lexeme) > 1 {
		return lex.newErrorToken(abstractToken, "Characters cannot contain more than one symbol")
	}

	if lex.isChar('\'') {
		// skip closing quote
		lex.nextChar()

		if len(lexeme) == 0 {
			return lex.newErrorToken(abstractToken, `Empty character. Use '\0' instead.`)
		}

		return &token.CharacterToken{
			AbstractToken: abstractToken,
			Value:         []rune(lexeme)[0],
		}
	}
	return lex.newErrorToken(abstractToken, "Character not closed")
}

func (lex *Lexer) readFixToken() token.Token {
	// Check if the character could belong to a multi character operation
	if symbol, ok := token.PossibleMultiCharOperators[lex.current]; ok {
		buf := []rune{lex.current}
		lex.nextChar()

		// Check if the concatenated characters really build a multi character operation
		if multiCharSymbol, ok := token.MultiCharOperators[string(buf[0])+string(lex.current)]; ok {
			buf = append(buf, lex.current)
			symbol = multiCharSymbol
			lex.nextChar()
		}

		abstractToken := lex.newAbstractToken(string(buf))

		return &token.FixToken{
			AbstractToken: abstractToken,
			Value:         symbol,
		}
	} else if symbol, ok := token.SingleCharOperators[lex.current]; ok { // Check if the character is a single character operator
		abstractToken := lex.newAbstractToken(string(lex.current))
		lex.nextChar()

		return &token.FixToken{
			AbstractToken: abstractToken,
			Value:         symbol,
		}
	}

	if lex.isChar('\n') {
		lex.nextChar()
		return &token.FixToken{
			AbstractToken: lex.newAbstractToken(`\n`),
			Value:         token.NewLine,
		}
	}

	lexeme := string(lex.current)
	lex.nextChar()

	return &token.ErrorToken{
		AbstractToken: lex.newAbstractToken(lexeme),
		Msg:           "Invalid character",
	}
}

func (lex *Lexer) nextChar() {
	if char, _, err := lex.reader.ReadRune(); err != nil {
		lex.current = 0
		if err == io.EOF {
			lex.isEnd = true
		} else {
			log.Fatal(err)
		}
	} else {
		lex.current = char

		if char == '\n' {
			lex.currentPos.NextLine()
		} else {
			lex.currentPos.MoveRight()
		}
	}

}

func (lex *Lexer) peekChar() (rune, error) {
	char, _, readError := lex.reader.ReadRune()

	if readError == nil {
		unreadError := lex.reader.UnreadRune()
		return char, unreadError
	}
	return char, readError
}

// Helpers
// -----------------------

type predicate func() bool

func (lex *Lexer) readLexeme(pred predicate) string {
	var buf []rune

	for !lex.isEnd && pred() {
		buf = append(buf, lex.current)
		lex.nextChar()
	}

	return string(buf)
}

func contains(s []rune, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

var allowedCharEscapedCodes = []rune{'0', 'n', '\'', '\\'}
var allowedStringEscapedCodes = []rune{'n', '"', '\\'}

var escapedChars = map[rune]rune{
	'0':  0,
	'n':  '\n',
	'\'': '\'',
	'\\': '\\',
	'"':  '"',
}

func (lex *Lexer) readEscapedLexeme(pred predicate, allowedCodes []rune) (string, error) {
	var buf []rune

	for !lex.isEnd && pred() {
		// Escape codes
		if lex.isChar('\\') {
			lex.nextChar()
			charToEscape := lex.current

			if contains(allowedCodes, charToEscape) {
				if escapedChar, ok := escapedChars[charToEscape]; ok {
					buf = append(buf, escapedChar)
				} else {
					panic(fmt.Sprintf("No escape code is found for %c", charToEscape))
				}
			} else {
				lex.nextChar()
				return string(buf), fmt.Errorf("escape code %c is not allowed", charToEscape)
			}
		} else if lex.current > 126 {
			buf = append(buf, lex.current)
			lex.nextChar()
			return string(buf), errors.New("unicode char is not allowed")
		} else {
			buf = append(buf, lex.current)
		}

		lex.nextChar()
	}

	return string(buf), nil
}

func (lex *Lexer) isLetter() bool {
	return lex.current >= 'A' && lex.current <= 'Z' ||
		lex.current >= 'a' && lex.current <= 'z'
}

func (lex *Lexer) isDigit() bool {
	return lex.current >= '0' && lex.current <= '9'
}

func (lex *Lexer) isHexDigit() bool {
	return lex.isDigit() ||
		lex.current >= 'a' && lex.current <= 'f' ||
		lex.current >= 'A' && lex.current <= 'F'
}

func (lex *Lexer) isChar(char rune) bool {
	return lex.current == char
}

func (lex *Lexer) newAbstractToken(lexeme string) token.AbstractToken {
	return token.AbstractToken{
		Position: lex.tokenPos,
		Lexeme:   lexeme,
	}
}

func (lex *Lexer) newErrorToken(abstractToken token.AbstractToken, msg string) *token.ErrorToken {
	return &token.ErrorToken{
		AbstractToken: abstractToken,
		Msg:           msg,
	}
}
