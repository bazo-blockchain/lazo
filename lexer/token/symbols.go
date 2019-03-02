package token

type Symbol int

const (
	EOF Symbol = iota

	Addition
	Subtraction
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

var Keywords = map[string]Symbol{
	"contract": Contract,
	"return":   Return,
	"if":       If,
	"else":     Else,
	"function": Function,
}

var SingleCharOperations = map[string]Symbol{
	":": Colon,
	",": Comma,
	".": Period,
	"{": OpenBrace,
	"}": CloseBrace,
	"[": OpenBracket,
	"]": CloseBracket,
	"(": OpenParen,
	")": CloseParen,
	"+": Addition,
	"-": Subtraction,
	"/": Division,
	"*": Multiplication,
	"%": Modulo,
}

var PossibleMultiCharOperation = map[string]Symbol{
	"=": Assign,
	">": Greater,
	"<": Less,
	"!": Not,
}

var LogicalOperation = map[string]Symbol{
	"&&": And,
	"||": Or,
}

var MultiCharOperation = map[string]Symbol{
	"==": Equal,
	"!=": Unequal,
	">=": GreaterEqual,
	"<=": LessEqual,
}
