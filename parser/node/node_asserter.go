package node

import (
	"gotest.tools/assert"
	"testing"
)

func AssertVariable(t *testing.T, node *VariableNode, varType string, id string) {
	assert.Equal(t, node.Type.Identifier, varType)
	assert.Equal(t, node.Identifier, id)
}
