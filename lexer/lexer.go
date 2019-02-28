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

	if lex.isLetter() {
		return lex.readIdentifier()
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

func (lex *Lexer) readIdentifier() *token.IdentifierToken {
	lexeme := lex.readLexeme(lex.isLetter)

	return &token.IdentifierToken{
		AbstractToken: token.AbstractToken{
			Position: lex.currentPos,
			Lexeme: lexeme,
		},
	}
}

func (lex *Lexer) readInteger() token.Token {
	// TODO: Hex Numbers
	lexeme := lex.readLexeme(lex.isDigit)
	value := new(big.Int)
	value, ok := value.SetString(lexeme, 10)

	abstractToken := lex.newAbstractToken(lexeme)
	if !ok {
		return lex.newErrorToken(abstractToken, "Error while parsing string to big int")
	}

	return &token.IntegerToken{
		AbstractToken: abstractToken,
		Value: value,
	}

}

func (lex *Lexer) readFixToken() *token.FixToken {
	// TODO Implement correctly
	lex.nextChar()
	return nil
}

func (lex *Lexer) readString() token.Token {
	// skip opening double quote
	lex.nextChar()

	lexeme := lex.readLexeme(func() bool {
		return !lex.isChar('"')
	})

	abstractToken := lex.newAbstractToken(lexeme)

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

func (lex *Lexer) readCharacter() *token.CharacterToken {
	return nil
}

func (lex *Lexer) nextChar() {
	if char, _, err := lex.reader.ReadRune(); err != nil {
		if err == io.EOF {
			lex.EOF = true
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

// Helpers

type predicate func() bool

func (lex *Lexer) readLexeme(pred predicate) string {
	buf := []rune{lex.current}
	lex.nextChar()

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
