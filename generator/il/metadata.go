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
			bytes := []byte{byte(code.OpCode)}
			if code.Operand != nil {
				bytes = append(bytes, code.Operand.([]byte)...)
			}
			fmt.Println(bytes)
			f.Write(bytes)
		}
	}
	f.Write([]byte{byte(HALT)})
}
