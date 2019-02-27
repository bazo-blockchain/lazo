package lexer

import (
	"gotest.tools/assert"
	"testing"
)

func TestIdentifierToken(t *testing.T) {
	i := IdentifierToken{
		abstractToken{
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
