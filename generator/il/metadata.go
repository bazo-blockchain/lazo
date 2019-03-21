package il

import (
	"encoding/json"
	"fmt"
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

	for _, function := range d.Contract.Functions {
		fmt.Printf("%s: \n", function.Identifier)
		// w.WriteString(fmt.Sprintf("%s: \n", function.Identifier)) // function calls does not work
		byteCounter := 0
		for _, code := range function.Instructions {
			if code.OpCode == RET {
				continue
			}
			bytes := []byte{byte(code.OpCode)}
			if code.Operand != nil {
				bytes = append(bytes, code.Operand.([]byte)...)
			}
			fmt.Printf("%d: %s %v \n", byteCounter, OpCodes[code.OpCode].Name, bytes)
			f.Write(bytes)
			byteCounter += len(bytes)
		}
	}
	f.Write([]byte{byte(HALT)})
}
