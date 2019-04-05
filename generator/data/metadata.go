package data

import (
	"fmt"
	"github.com/bazo-blockchain/bazo-miner/protocol"
	"github.com/bazo-blockchain/lazo/generator/il"
)

// Metadata contains the ContractData
type Metadata struct {
	Contract *ContractData
}

// CreateContract returns the byte code for the contract together with the contract fields
func (d *Metadata) CreateContract() ([]byte, []protocol.ByteArray) {
	return d.getByteCode(), d.getVariables()
}

func (d *Metadata) getByteCode() []byte {
	var byteCode []byte
	bytePos := 0

	for _, code := range d.Contract.Instructions {
		bytes := generateByteCode(code, bytePos)
		byteCode = append(byteCode, bytes...)
		bytePos += len(bytes)
	}

	for _, function := range d.Contract.Functions {
		fmt.Printf("%s: \n", function.Identifier)
		for _, code := range function.Instructions {
			bytes := generateByteCode(code, bytePos)
			byteCode = append(byteCode, bytes...)
			bytePos += len(bytes)
		}
	}
	return byteCode
}

func (d *Metadata) getVariables() []protocol.ByteArray {
	return make([]protocol.ByteArray, len(d.Contract.Fields))
}

func generateByteCode(code *il.Instruction, bytePos int) []byte {
	bytes := []byte{byte(code.OpCode)}
	if code.Operand != nil {
		bytes = append(bytes, code.Operand.([]byte)...)
	}
	fmt.Printf("%d: %s %v \n", bytePos, il.OpCodes[code.OpCode].Name, bytes)
	return bytes
}
