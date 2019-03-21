package il

import (
	"bufio"
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

func (d *Metadata) SaveBazoByteCode(outputFile string) {
	f, err := os.Create(outputFile)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	for _, function := range d.Contract.Functions {
		// w.WriteString(fmt.Sprintf("%s: \n", function.Identifier)) // function calls does not work
		for _, code := range function.Instructions {
			if code.OpCode == RET {
				continue
			}
			var operand interface{}
			if code.Operand != nil {
				operand = fmt.Sprintf("%v", code.Operand)
			} else {
				operand = ""
			}
			w.WriteString(fmt.Sprintf("%s %v \n", OpCodeLiterals[code.OpCode], operand))
		}
	}
}
