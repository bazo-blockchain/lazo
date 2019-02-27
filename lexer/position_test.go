package lexer

import (
	"gotest.tools/assert"
	"testing"
)

func TestMoveRight(t *testing.T) {
	pos := New()
	assert.Equal(t, pos.Line, 1)
	assert.Equal(t, pos.Column, 0)

	pos.MoveRight()
	assert.Equal(t, pos.Line, 1)
	assert.Equal(t, pos.Column, 1)
}

func TestNextLine(t *testing.T) {
	pos := New()
	pos.MoveRight()
	pos.NextLine()
	assert.Equal(t, pos.Line, 2)
	assert.Equal(t, pos.Column, 0)
}

func TestString(t *testing.T) {
	pos := New()
	assert.Equal(t, pos.String(), "1:0")

	pos.MoveRight()
	assert.Equal(t, pos.String(), "1:1")
}
