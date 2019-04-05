// Package token contains all the supported token types and their functions.
package token

import (
	"fmt"
	"math/big"
)

// TokenType is the type of supported token types
type TokenType int

// Supported TokenTypes
const (
	IDENTIFER TokenType = iota
	INTEGER
	STRING
	CHARACTER
	SYMBOL
	ERROR
)

// Token is the interface that wraps the basic Token functions
type Token interface {
	Pos() Position
	Literal() string
	String() string
	Type() TokenType
}

// AbstractToken contains token position and lexeme, which all concrete tokens have.
type AbstractToken struct {
	Position
	Lexeme string
}

// Pos returns the token position
func (t *AbstractToken) Pos() Position {
	return t.Position
}

// Literal returns the actual token lexeme
func (t *AbstractToken) Literal() string {
	return t.Lexeme
}

// Concrete Tokens
// ----------------

// IdentifierToken holds the identifier and compose abstract token
type IdentifierToken struct {
	AbstractToken
}

// Type returns the token type
func (t *IdentifierToken) Type() TokenType {
	return IDENTIFER
}

func (t *IdentifierToken) String() string {
	return fmt.Sprintf("[%s] IDENTIFER %s", t.Pos(), t.Literal())
}

// --------------------------

// IntegerToken holds integer number and compose abstract token
type IntegerToken struct {
	AbstractToken
	Value *big.Int
	IsHex bool
}

// Type returns the token type
func (t *IntegerToken) Type() TokenType {
	return INTEGER
}

func (t *IntegerToken) String() string {
	return fmt.Sprintf("[%s] INT %s", t.Pos(), t.Literal())
}

// --------------------------

// StringToken holds a string literal and compose abstract token
type StringToken struct {
	AbstractToken
}

// Type returns the token type
func (t *StringToken) Type() TokenType {
	return STRING
}

func (t *StringToken) String() string {
	return fmt.Sprintf("[%s] STRING %s", t.Pos(), t.Literal())
}

// --------------------------

// CharacterToken holds a character literal and compose abstract token
type CharacterToken struct {
	AbstractToken
	Value rune
}

// Type returns the token type
func (t *CharacterToken) Type() TokenType {
	return CHARACTER
}

func (t *CharacterToken) String() string {
	return fmt.Sprintf("[%s] CHAR %s", t.Pos(), t.Literal())
}

// --------------------------

// FixToken holds fix symbols and compose abstract token
type FixToken struct {
	AbstractToken
	Value Symbol
}

// Type returns the token type
func (t *FixToken) Type() TokenType {
	return SYMBOL
}

func (t *FixToken) String() string {
	return fmt.Sprintf("[%s] SYMBOL %s", t.Pos(), t.Literal())
}

// --------------------------

// ErrorToken holds lexer error and compose abstract token
type ErrorToken struct {
	AbstractToken
	Msg string
}

// Type returns the token type
func (t *ErrorToken) Type() TokenType {
	return ERROR
}

func (t *ErrorToken) String() string {
	return fmt.Sprintf("[%s] Error: %s - %s", t.Pos(), t.Msg, t.Literal())
}
