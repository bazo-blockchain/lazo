package lexer

import (
	"bufio"
	"strings"
)

func createLexerFromInput(input string) *Lexer{
	return New(bufio.NewReader(strings.NewReader(input)))
}
