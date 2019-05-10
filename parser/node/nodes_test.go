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

// Type Nodes
// ----------

func TestBasicTypeNode_Type(t *testing.T) {
	basicType := &BasicTypeNode{
		Identifier: "int",
	}

	assert.Equal(t, basicType.Type(), "int")
}

func TestArrayTypeNode_Type(t *testing.T) {
	basicType := &BasicTypeNode{
		Identifier: "int",
	}
	arrayType := &ArrayTypeNode{
		ElementType: basicType,
	}

	assert.Equal(t, arrayType.Type(), "int[]")
}

func TestArrayTypeNode_Type_2D(t *testing.T) {
	basicType := &BasicTypeNode{
		Identifier: "int",
	}
	arrayType := &ArrayTypeNode{
		ElementType: basicType,
	}
	arrayType2D := &ArrayTypeNode{
		ElementType: arrayType,
	}

	assert.Equal(t, arrayType2D.Type(), "int[][]")
}

func TestMapTypeNode_Type(t *testing.T) {
	basicType := &BasicTypeNode{
		Identifier: "int",
	}
	mapType := &MapTypeNode{
		KeyType:   basicType,
		ValueType: basicType,
	}

	assert.Equal(t, mapType.Type(), "Map<int,int>")
}
