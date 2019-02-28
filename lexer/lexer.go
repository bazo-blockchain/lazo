package lexer

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"io"
	"log"
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

/*	if lex.isDigit() {
		return lex.readInteger()
	}*/

	if lex.isLetter() {
		return lex.readIdentifier()
	} else {
		lex.nextChar()
	}

	// read string
	// read character
	// read other fix tokens

	// todo remove - scan all characters and print

	// ------------

	return nil
}

func (lex *Lexer) skipWhiteSpace() {
	for !lex.EOF && lex.current <= ' ' {
		lex.nextChar()
	}
}

func (lex *Lexer) readIdentifier() *token.IdentifierToken {
	buf := []rune{lex.current}
	lex.nextChar()

	for !lex.EOF && lex.isLetter() {
		buf = append(buf, lex.current)
		lex.nextChar()
	}

	return &token.IdentifierToken{
		AbstractToken: token.AbstractToken{
			Position: lex.currentPos,
			Lexeme: string(buf),
		},
	}
}

func (lex *Lexer) readInteger() *token.IntegerToken {
	return nil
}

func (lex *Lexer) readFixToken() *token.FixToken {
	return nil
}

func (lex *Lexer) readString() *token.StringToken {
	return nil
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

// Helper Functions

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
