package token

import "fmt"

// Position holds the line and column number
type Position struct {
	Line   int
	Column int
}

// NewPosition creates a new position with line 1 and column 0
func NewPosition() Position {
	return Position{
		Line:   1,
		Column: 0,
	}
}

// MoveRight increments the column by 1
func (pos *Position) MoveRight() {
	pos.Column++
}

// NextLine increments the line by 1 and reset the column to 0
func (pos *Position) NextLine() {
	pos.Line++
	pos.Column = 0
}

func (pos Position) String() string {
	return fmt.Sprintf("%d:%d", pos.Line, pos.Column)
}
