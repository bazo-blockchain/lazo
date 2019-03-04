package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"io"
	"log"
	"math/big"
)

type Lexer struct {
	reader     *bufio.Reader
	current    rune
	currentPos token.Position
	tokenPos   token.Position
	IsEnd      bool
}

func New(reader *bufio.Reader) *Lexer {
	lex := &Lexer{
		reader:     reader,
		currentPos: token.NewPosition(),
	}
	lex.nextChar()
	return lex
}

func (lex *Lexer) NextToken() token.Token {
	lex.skipWhiteSpace()

	lex.tokenPos = lex.currentPos

	if lex.IsEnd {
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
	for !lex.IsEnd && lex.current <= ' ' {
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
		return !lex.isChar(' ') && !lex.isChar('\n')
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
	} else {
		return lex.newErrorToken(abstractToken, "String not closed")
	}
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

		return &token.CharacterToken{
			AbstractToken: abstractToken,
			Value:         []rune(lexeme)[0],
		}
	} else {
		return lex.newErrorToken(abstractToken, "Character not closed")
	}
}

func (lex *Lexer) readFixToken() token.Token {

	// Check if the character could belong to a multi character operation
	if symbol, ok := token.PossibleMultiCharOperation[string(lex.current)]; ok {
		buf := []rune{lex.current}

		lex.nextChar()

		// Check if the concatenated characters really build a multi character operation
		if multiCharSymbol, ok := token.MultiCharOperation[string(buf[0])+string(lex.current)]; ok {
			buf = append(buf, lex.current)
			symbol = multiCharSymbol
			lex.nextChar()
		}

		abstractToken := lex.newAbstractToken(string(buf))

		return &token.FixToken{
			AbstractToken: abstractToken,
			Value:         symbol,
		}
	}

	// Check if the character is a single character operator
	if symbol, ok := token.SingleCharOperations[string(lex.current)]; ok {

		lex.nextChar()

		return &token.FixToken{
			AbstractToken: lex.newAbstractToken(string(lex.current)),
			Value:         symbol,
		}
	}

	if lex.current == '&' || lex.current == '|' {
		return lex.readLogicalFixToken()
	}

	lex.nextChar()
	return nil

}

func (lex *Lexer) readLogicalFixToken() token.Token {
	buf := []rune{lex.current}
	lex.nextChar()
	buf = append(buf, lex.current)

	abstractToken := lex.newAbstractToken(string(buf))

	if symbol, ok := token.LogicalOperation[string(buf)]; ok {
		lex.nextChar()

		return &token.FixToken{
			AbstractToken: abstractToken,
			Value:         symbol,
		}
	} else {
		return lex.newErrorToken(abstractToken, "Unknown Symbol")
	}
}

func (lex *Lexer) nextChar() {
	if char, _, err := lex.reader.ReadRune(); err != nil {
		lex.current = 0
		if err == io.EOF {
			lex.IsEnd = true
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
	} else {
		return char, readError
	}
}

// Helpers
// -----------------------

type predicate func() bool

func (lex *Lexer) readLexeme(pred predicate) string {
	var buf []rune

	for !lex.IsEnd && pred() {
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

	for !lex.IsEnd && pred() {
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
				return string(buf), errors.New(fmt.Sprintf("Escape code %c is not allowed", charToEscape))
			}
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
