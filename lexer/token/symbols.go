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

	Assign

	// Keywords

	Contract
	Return
	If
	Else
	Function
)

var Keywords = map[string]Symbol {
	"contract": Contract,
	"return": Return,
	"if": If,
	"else": Else,
	"function": Function,
}

var SingleCharOperations = map[string]Symbol {
	":": Colon,
	",": Comma,
	".": Period,
	"{": OpenBrace,
	"}": CloseBrace,
	"[": OpenBracket,
	"]": CloseBracket,
	"(": OpenParen,
	")": CloseParen,
}