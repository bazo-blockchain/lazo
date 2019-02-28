package token

import (
	"fmt"
	"math/big"
)

type Token interface {
	Pos() Position
	Literal() string
	String() string
}

type abstractToken struct {
	Position
	Lexeme string
}

func (t *abstractToken) Pos() Position {
	return t.Position
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

type BooleanToken struct {
	abstractToken
	Value bool
}

func (t *BooleanToken) String() string {
	return fmt.Sprintf("[%s] BOOLEAN %s", t.Pos(), t.Literal())
}

// --------------------------

type StringToken struct {
	abstractToken
}

func (t *StringToken) String() string {
	return fmt.Sprintf("[%s] STRING %s", t.Pos(), t.Literal())
}

// --------------------------

type CharacterToken struct {
	abstractToken
	Value rune
}

func (t *CharacterToken) String() string {
	return fmt.Sprintf("[%s] CHAR %s", t.Pos(), t.Literal())
}

// --------------------------

type FixToken struct {
	abstractToken
	Value Symbol
}

func (t *FixToken) String() string {
	return fmt.Sprintf("[%s] SYMBOL %s", t.Pos(), t.Literal())
}
