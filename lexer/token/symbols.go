package token

type Symbol int

const (
	Plus Symbol = iota
	Minus
	Division
	Multiplication
	Modulo

	Less
	LessEqual
	GreaterEqual
	Greater

	Equal
	Unequal

	OpenBrace
	CloseBrace
	OpenBracket
	CloseBracket
	OpenParen
	CloseParen

	Colon
	Comma
	Period

	Not
	And
	Or

	// Keywords

	Contract
	Assign
	Return
	If
	Else
	Function
)
