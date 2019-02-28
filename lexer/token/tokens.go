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

type AbstractToken struct {
	Position
	Lexeme string
}

func (t *AbstractToken) Pos() Position {
	return t.Position
}

func (t *AbstractToken) Literal() string {
	return t.Lexeme
}

// Concrete Tokens
// ----------------

type IdentifierToken struct {
	AbstractToken
}

func (t *IdentifierToken) String() string {
	return fmt.Sprintf("[%s] IDENTIFER %s", t.Pos(), t.Literal())
}

// --------------------------

type IntegerToken struct {
	AbstractToken
	Value *big.Int
}

func (t *IntegerToken) String() string {
	return fmt.Sprintf("[%s] INT %s", t.Pos(), t.Literal())
}

// --------------------------

type BooleanToken struct {
	AbstractToken
	Value bool
}

func (t *BooleanToken) String() string {
	return fmt.Sprintf("[%s] BOOLEAN %s", t.Pos(), t.Literal())
}

// --------------------------

type StringToken struct {
	AbstractToken
}

func (t *StringToken) String() string {
	return fmt.Sprintf("[%s] STRING %s", t.Pos(), t.Literal())
}

// --------------------------

type CharacterToken struct {
	AbstractToken
	Value rune
}

func (t *CharacterToken) String() string {
	return fmt.Sprintf("[%s] CHAR %s", t.Pos(), t.Literal())
}

// --------------------------

type FixToken struct {
	AbstractToken
	Value Symbol
}

func (t *FixToken) String() string {
	return fmt.Sprintf("[%s] SYMBOL %s", t.Pos(), t.Literal())
}

// --------------------------

type ErrorToken struct {
	AbstractToken
	Msg string
}

func (t *ErrorToken) String() string {
	return fmt.Sprintf("[%s] Error %s - %s", t.Pos(), t.Msg, t.Literal())
}
