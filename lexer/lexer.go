package lexer

import (
	"bufio"
	"fmt"
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
	return &Lexer{
		reader: reader,
		currentPos: token.NewPosition(),
	}
}

func (lex *Lexer) NextToken() token.Token {
	// skip whitespaces

	// read identifier
	// read integer
	// read string
	// read character
	// read other fix tokens

	// todo remove - scan all characters and print
	// ------------
	for !lex.EOF {
		lex.nextChar()
		fmt.Printf("%s %c\n", lex.currentPos, lex.current)
	}
	// ------------

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
