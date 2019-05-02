// Package data contains all the supported metadata types.
package data

import "github.com/bazo-blockchain/lazo/generator/il"

// ContractData contains the identifier, total fields, functions and instructions
type ContractData struct {
	Identifier   string
	TotalFields  uint16
	Functions    []*FunctionData
	Instructions []*il.Instruction
}

// FunctionData contains an the identifier, return types, parameter types, local variable types, instructions and the
// function hash
type FunctionData struct {
	Identifier   string
	Instructions []*il.Instruction
	Hash         [4]byte
}
