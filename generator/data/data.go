package data

import "github.com/bazo-blockchain/lazo/generator/il"

type TypeData int

const (
	_                 = iota
	BoolType TypeData = iota * -1
	IntType
	CharType
	StringType
)

type ContractData struct {
	Identifier   string
	Fields       []TypeData
	Functions    []*FunctionData
	Instructions []*il.Instruction
}

type FunctionData struct {
	Identifier   string
	ReturnTypes  []TypeData
	ParamTypes   []TypeData
	LocalTypes   []TypeData
	Instructions []*il.Instruction
	Hash         [4]byte
}
