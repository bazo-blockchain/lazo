package lexer

import (
	"fmt"
	"math/big"
)

type Token interface {
	Pos() *Position
	Literal() string
	String() string
}

type abstractToken struct {
	Position
	Lexeme string
}

func (t *abstractToken) Pos() *Position {
	return &t.Position
}

func (t *abstractToken) Literal() string {
	return t.Lexeme
}

// Concrete Tokens
// ----------------

type IdentifierToken struct {
	abstractToken
}

func (t *IdentifierToken) String() string {
	return fmt.Sprintf("[%s] IDENTIFER %s", t.Pos(), t.Literal())
}

// --------------------------

type IntegerToken struct {
	abstractToken
	Value big.Int
}

func (t *IntegerToken) String() string {
	return fmt.Sprintf("[%s] INT %s", t.Pos(), t.Literal())
}

// --------------------------

// todo define tokens
