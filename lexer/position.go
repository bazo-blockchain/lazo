package lexer

import "fmt"

type Position struct {
	Line   int
	Column int
}

func New() *Position {
	return &Position{
		Line:   1,
		Column: 0,
	}
}

func (pos *Position) MoveRight() {
	pos.Column++
}

func (pos *Position) NextLine() {
	pos.Line++
	pos.Column = 0
}

func (pos Position) String() string {
	return fmt.Sprintf("%d:%d", pos.Line, pos.Column)
}
