package lexer

type Token interface {
	Pos() Position
}
