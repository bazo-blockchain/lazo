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

type FunctionData struct {
	Identifier   string            `json:"ID"`
	ReturnTypes  []TypeData        `json:"ReturnTypes"`
	ParamTypes   []TypeData        `json:"ParamTypes"`
	LocalTypes   []TypeData        `json:"LocalTypes"`
	Instructions []*il.Instruction `json:"Instructions"`
}

type ContractData struct {
	Identifier   string            `json:"ID"`
	Fields       []TypeData        `json:"Fields"`
	Functions    []*FunctionData   `json:"Functions"`
	Instructions []*il.Instruction `json:"Instructions"`
}
