package token

// Symbol is the type of allowed fixed symbols in the source code
type Symbol int

// Allowed fixed symbols
const (
	EOF Symbol = iota
	NewLine

	Addition
	Subtraction
	Division
	Multiplication
	Modulo
	Exponent

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

// SymbolLexeme maps the Symbol type to its lexeme value
var SymbolLexeme = map[Symbol]string{
	EOF:     "EOF",
	NewLine: `\n`,

	Addition:       "+",
	Subtraction:    "-",
	Multiplication: "*",
	Division:       "/",
	Modulo:         "%",
	Exponent:       "**",

	Less:         "<",
	LessEqual:    "<=",
	GreaterEqual: ">=",
	Greater:      ">",

	Equal:   "==",
	Unequal: "!=",

	OpenBrace:    "{",
	CloseBrace:   "}",
	OpenBracket:  "[",
	CloseBracket: "]",
	OpenParen:    "(",
	CloseParen:   ")",

	Colon:  ":",
	Comma:  ",",
	Period: ".",

	Not: "!",
	And: "&&",
	Or:  "||",

	Assign: "=",

	// Keywords

	Contract: "contract",
	Return:   "return",
	If:       "if",
	Else:     "else",
	Function: "function",
	True:     "true",
	False:    "false",
}

// Keywords maps reserved literal values to the Symbol type
var Keywords = map[string]Symbol{
	"contract": Contract,
	"return":   Return,
	"if":       If,
	"else":     Else,
	"function": Function,
	"true":     True,
	"false":    False,
}

// BooleanConstants maps Symbol type to built-in boolean value
var BooleanConstants = map[Symbol]bool{
	True:  true,
	False: false,
}

// SingleCharOperators maps single character to Symbol type
var SingleCharOperators = map[rune]Symbol{
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
	'%': Modulo,
}

// PossibleMultiCharOperators maps the first character of a possible operation to the Symbol type.
var PossibleMultiCharOperators = map[rune]Symbol{
	'=': Assign,
	'>': Greater,
	'<': Less,
	'!': Not,

	'*': Multiplication,
}

// MultiCharOperators maps the full literal value of the operation to the Symbol type.
var MultiCharOperators = map[string]Symbol{
	"==": Equal,
	"!=": Unequal,
	">=": GreaterEqual,
	"<=": LessEqual,

	"**": Exponent,
}

// LogicalOperators maps the logical operator to Symbol type.
var LogicalOperators = map[string]Symbol{
	"&&": And,
	"||": Or,
}
