package data

import "github.com/bazo-blockchain/lazo/generator/il"

// TypeData is an alias for int and used for type representation
type TypeData int

const (
	_ = iota
	// BoolType represents the boolean type
	BoolType TypeData = iota * -1
	// IntType represents the integer type
	IntType
	// CharType represents the character type
	CharType
	// StringType represents the string type
	StringType
)

// ContractData contains the identifier, fields, functions and instructions
type ContractData struct {
	Identifier   string
	Fields       []TypeData
	Functions    []*FunctionData
	Instructions []*il.Instruction
}

// FunctionData contains an the identifier, return types, parameter types, local variable types, instructions and the
// function hash
type FunctionData struct {
	Identifier   string
	ReturnTypes  []TypeData
	ParamTypes   []TypeData
	LocalTypes   []TypeData
	Instructions []*il.Instruction
	Hash         [4]byte
}
