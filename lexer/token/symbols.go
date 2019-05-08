package token

// Symbol is the type of allowed fixed symbols in the source code
type Symbol int

// Allowed fixed symbols
const (
	EOF Symbol = iota
	NewLine

	Plus
	Minus
	Division
	Multiplication
	Modulo
	Exponent

	ShiftLeft
	ShiftRight

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

	BitwiseNot
	BitwiseAnd
	BitwiseOr
	BitwiseXor

	Assign

	// Keywords

	Contract
	Struct
	Map
	Delete
	New
	Constructor
	If
	Else
	Function
	Return
	True
	False
)

// SymbolLexeme maps the Symbol type to its lexeme value
var SymbolLexeme = map[Symbol]string{
	EOF:     "EOF",
	NewLine: `\n`,

	Plus:           "+",
	Minus:          "-",
	Multiplication: "*",
	Division:       "/",
	Modulo:         "%",
	Exponent:       "**",

	ShiftLeft:  "<<",
	ShiftRight: ">>",

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

	BitwiseNot: "~",
	BitwiseAnd: "&",
	BitwiseOr:  "|",
	BitwiseXor: "^",

	Assign: "=",

	// Keywords

	Contract:    "contract",
	Struct:      "struct",
	Map:         "Map",
	Delete:      "delete",
	New:         "new",
	Constructor: "constructor",
	If:          "if",
	Else:        "else",
	Function:    "function",
	Return:      "return",
	True:        "true",
	False:       "false",
}

// Keywords maps reserved literal values to the Symbol type
var Keywords = map[string]Symbol{
	"contract":    Contract,
	"struct":      Struct,
	"Map":         Map,
	"delete":      Delete,
	"new":         New,
	"constructor": Constructor,
	"if":          If,
	"else":        Else,
	"function":    Function,
	"return":      Return,
	"true":        True,
	"false":       False,
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
	'+': Plus,
	'-': Minus,
	'/': Division,
	'%': Modulo,
	'~': BitwiseNot,
	'^': BitwiseXor,
}

// PossibleMultiCharOperators maps the first character of a possible operation to the Symbol type.
var PossibleMultiCharOperators = map[rune]Symbol{
	'=': Assign,
	'>': Greater,
	'<': Less,
	'!': Not,

	'*': Multiplication,

	'|': BitwiseOr,
	'&': BitwiseAnd,
}

// MultiCharOperators maps the full literal value of the operation to the Symbol type.
var MultiCharOperators = map[string]Symbol{
	"==": Equal,
	"!=": Unequal,
	">=": GreaterEqual,
	"<=": LessEqual,

	"**": Exponent,

	"<<": ShiftLeft,
	">>": ShiftRight,

	"||": Or,
	"&&": And,
}
