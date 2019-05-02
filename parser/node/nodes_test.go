package node

import (
	"github.com/bazo-blockchain/lazo/lexer/token"
	"gotest.tools/assert"
	"testing"
)

func TestContractNode_String(t *testing.T) {
	contract := &ContractNode{
		AbstractNode: AbstractNode{
			Position: token.Position{
				Line:   1,
				Column: 1},
		},
		Name: "Test",
	}
	assert.Equal(t, contract.String(),
		"[1:1] CONTRACT Test \n FIELDS: [] \n\n STRUCTS: [] \n\n CONSTRUCTOR:  \n\n FUNCS: []")
}
