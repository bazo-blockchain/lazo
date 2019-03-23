package data

import (
	"encoding/json"
	"fmt"
	"github.com/bazo-blockchain/lazo/generator/il"
	"io/ioutil"
	"os"
)

type Metadata struct {
	Contract *ContractData
}

func (d *Metadata) Save(destinationFile string) {
	// TODO Error Handling
	contract, err := json.MarshalIndent(d.Contract, "", " ")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile(destinationFile, contract, 0644)
}

func (d *Metadata) SaveByteCode(outputFile string) {
	f, err := os.Create(outputFile)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	f.Write(d.GetByteCode())
}

func (d *Metadata) GetByteCode() []byte {
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

func generateByteCode(code *il.Instruction, bytePos int) []byte {
	bytes := []byte{byte(code.OpCode)}
	if code.Operand != nil {
		bytes = append(bytes, code.Operand.([]byte)...)
	}
	fmt.Printf("%d: %s %v \n", bytePos, il.OpCodes[code.OpCode].Name, bytes)
	return bytes
}
