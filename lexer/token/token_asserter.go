package token

import (
	"gotest.tools/assert"
	"math/big"
	"testing"
)

func AssertInteger(t *testing.T, token Token, value *big.Int){
	tok, ok := token.(*IntegerToken)

	assert.Equal(t, ok, true)
	assert.Equal(t, tok.Value.Cmp(value) == 0, true)
}

func AssertIdentifier(t *testing.T, token Token, value string){
	tok, ok := token.(*IdentifierToken)

	assert.Equal(t, ok, true)
	assert.Equal(t, tok.Literal(), value)
}

func AssertError(t *testing.T, token Token, value string) {
	tok, ok := token.(*ErrorToken)

	assert.Equal(t, ok, true)
	assert.Equal(t, tok.Literal(), value)
}