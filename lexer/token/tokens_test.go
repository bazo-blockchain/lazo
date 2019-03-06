package token

import (
	"gotest.tools/assert"
	"testing"
)

func TestTokenTypes(t *testing.T) {
	assert.Equal(t, (&IdentifierToken{}).Type(), IDENTIFER)
	assert.Equal(t, (&IntegerToken{}).Type(), INTEGER)
	assert.Equal(t, (&BooleanToken{}).Type(), BOOLEAN)
	assert.Equal(t, (&StringToken{}).Type(), STRING)
	assert.Equal(t, (&CharacterToken{}).Type(), CHARACTER)
	assert.Equal(t, (&FixToken{}).Type(), SYMBOL)
	assert.Equal(t, (&ErrorToken{}).Type(), ERROR)
}

func TestIdentifierToken(t *testing.T) {
	i := IdentifierToken{
		AbstractToken{
			Position{
				Line:   1,
				Column: 1,
			},
			"test",
		},
	}

	assert.Equal(t, i.Pos().String(), "1:1")
	assert.Equal(t, i.Literal(), "test")
	assert.Equal(t, i.String(), "[1:1] IDENTIFER test")
}
