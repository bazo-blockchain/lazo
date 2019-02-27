package lexer

type Token interface {
	Pos() Position
	String() string
	Literal() string
}
