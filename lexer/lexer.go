package lexer

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"io"
	"log"
	"math/big"
)

type Lexer struct {
	reader     *bufio.Reader
	current    rune
	currentPos token.Position
	EOF        bool
}

func New(reader *bufio.Reader) *Lexer {
	lex := &Lexer{
		reader: reader,
		currentPos: token.NewPosition(),
	}
	lex.nextChar()
	return lex
}

func (lex *Lexer) NextToken() token.Token {
	lex.skipWhiteSpace()

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
	for !lex.EOF && lex.current <= ' ' {
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
		Value: value,
		IsHex: isHex,
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

func (lex *Lexer) readFixToken() token.Token {

	if lex.isPossibleMultiCharFixToken() {
		buf := []rune {lex.current}

		symbol, _ := token.PossibleMultiCharOperation[string(buf)]

		lex.nextChar()

		if lex.isMultiCharFixToken(buf[0]) {
			buf = append(buf, lex.current)
			symbol, _ = token.MultiCharOperation[string(buf)]
			lex.nextChar()
		}

		abstractToken := lex.newAbstractToken(string(buf))

		return &token.FixToken{
			AbstractToken: abstractToken,
			Value: symbol,
		}
	}

	if lex.isSingleCharFixToken() {
		return lex.readSingleCharFixToken()
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
			 Value: symbol,
		}
	} else {
		return lex.newErrorToken(abstractToken, "Unknown Symbol")
	}
}

func (lex *Lexer) readSingleCharFixToken() *token.FixToken {
	lexeme := string(lex.current)

	symbol, _ := token.SingleCharOperations[lexeme]

	lex.nextChar()

	return &token.FixToken{
		AbstractToken: lex.newAbstractToken(lexeme),
		Value: symbol,
	}
}

func (lex *Lexer) readString() token.Token {
	// skip opening double quote
	lex.nextChar()

	var buf []rune

	for !lex.EOF && !lex.isChar('"'){
		// Escaping
		if lex.isChar('\\') {
			escapedChar := lex.current
			lex.nextChar()

			if lex.current == 'n' {
				escapedChar = '\n'
			}

			if lex.current == '\\' {
				escapedChar = '\\'
			}

			if lex.current == '"' {
				escapedChar = '"'
			}

			buf = append(buf, escapedChar)

		} else {

			buf = append(buf, lex.current)

			}

		lex.nextChar()
	}

	abstractToken := lex.newAbstractToken(string(buf))

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
	var buf []rune

	for !lex.EOF && !lex.isChar('\''){
		// Escaping
		if lex.isChar('\\') {
			escapedChar := lex.current
			lex.nextChar()

			if lex.current == 'n' {
				escapedChar = '\n'
			}

			if lex.current == '\\' {
				escapedChar = '\\'
			}

			if lex.current == '"' {
				escapedChar = '"'
			}

			buf = append(buf, escapedChar)

		} else {

			buf = append(buf, lex.current)

		}

		lex.nextChar()
	}

	abstractToken := lex.newAbstractToken(string(buf))

	if len(buf) > 1 {
		return lex.newErrorToken(abstractToken, "Characters cannot contain more than one symbol")
	}

	if lex.isChar('\'') {
		// skip closing quote
		lex.nextChar()

		return &token.CharacterToken{
			AbstractToken: abstractToken,
		}
	} else {
		return lex.newErrorToken(abstractToken, "Character not closed")
	}

}

func (lex *Lexer) nextChar() {
	char, _, err := lex.reader.ReadRune()

	if  err != nil {
		if err == io.EOF {
			lex.EOF = true
		} else {
			log.Fatal(err)
		}
	} else {
		if char == '\n' {
			lex.currentPos.NextLine()
		} else {
			lex.currentPos.MoveRight()
		}
	}

	lex.current = char
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

	for !lex.EOF && pred() {
		buf = append(buf, lex.current)
		lex.nextChar()
	}

	return string(buf)
}

func (lex *Lexer) isLetter() bool {
	return lex.current >= 'A' && lex.current <= 'Z' ||
		lex.current >= 'a' && lex.current <= 'z'
}

func (lex *Lexer) isDigit() bool {
	return lex.current >= '0' && lex.current <= '9'
}

// -----TODO Fix Duplicate Code ------
func (lex *Lexer) isSingleCharFixToken() bool {
	_, ok := token.SingleCharOperations[string(lex.current)]
	return ok
}

func (lex *Lexer) isPossibleMultiCharFixToken() bool {
	_, ok := token.PossibleMultiCharOperation[string(lex.current)]
	return ok
}
// ------------------------------------

func (lex *Lexer) isMultiCharFixToken(char rune) bool {
	_, ok := token.MultiCharOperation[string(char) + string(lex.current)]
	return ok
}

func (lex *Lexer) isHexDigit() bool {
	return lex.isDigit() ||
		lex.current >='a' && lex.current <='f' ||
		lex.current >= 'A' && lex.current <='F'
}

func (lex *Lexer) isChar(char rune) bool {
	return lex.current == char
}

func (lex *Lexer) newAbstractToken(lexeme string) token.AbstractToken {
	return token.AbstractToken {
		Position: lex.currentPos,
		Lexeme: lexeme,
	}
}

func (lex *Lexer) newErrorToken(abstractToken token.AbstractToken, msg string) *token.ErrorToken {
	return &token.ErrorToken {
		AbstractToken: abstractToken,
		Msg: msg,
	}
}
