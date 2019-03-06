package token

type Symbol int

const (
	EOF Symbol = iota
	NewLine

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
	True
	False
)

var SymbolLexeme = map[Symbol]string{
	EOF: "EOF",
	NewLine: `\n`,

	Addition: "+",
	Subtraction: "-",
	Multiplication: "*",
	Division: "/",
	Modulo: "%",

	Less: "<",
	LessEqual: "<=",
	GreaterEqual: ">=",
	Greater: ">",

	Equal: "==",
	Unequal: "!=",

	OpenBrace: "{",
	CloseBrace: "}",
	OpenBracket: "[",
	CloseBracket: "]",
	OpenParen: "(",
	CloseParen: ")",

	Colon: ":",
	Comma: ",",
	Period: ".",

	Not: "!",
	And: "&&",
	Or: "||",

	Assign: "=",

	// Keywords

	Contract: "contract",
	Return: "return",
	If: "if",
	Else: "else",
	Function: "function",
	True: "true",
	False: "false",
}

var Keywords = map[string]Symbol{
	"contract": Contract,
	"return":   Return,
	"if":       If,
	"else":     Else,
	"function": Function,
	"true":     True,
	"false":    False,
}

var BooleanConstants = map[Symbol]bool{
	True: true,
	False: false,
}

var SingleCharOperations = map[rune]Symbol{
	':': Colon,
	',': Comma,
	'.': Period,
	'{': OpenBrace,
	'}': CloseBrace,
	'[': OpenBracket,
	']': CloseBracket,
	'(': OpenParen,
	')': CloseParen,
	'+': Addition,
	'-': Subtraction,
	'/': Division,
	'*': Multiplication,
	'%': Modulo,
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
