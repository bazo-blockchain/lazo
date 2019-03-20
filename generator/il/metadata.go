package il

import (
	"encoding/json"
	"io/ioutil"
)

type MetaData struct {
	Contract *ContractData
}

func (d *MetaData) Save(destinationFile string) {
	// TODO Error Handling
	contract, err := json.MarshalIndent(d.Contract, "", " ")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile(destinationFile, contract, 0644)
}
