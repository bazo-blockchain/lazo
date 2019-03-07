package node

import (
	"gotest.tools/assert"
	"testing"
)

func AssertProgram(t *testing.T, node *ProgramNode, hasContract bool) {
	assert.Equal(t, node.Contract != nil, hasContract)
}

func AssertContract(t *testing.T, node *ContractNode, name string, totalVars int, totalFunctions int) {
	assert.Equal(t, node.Name, name)
	assert.Equal(t, len(node.Variables), totalVars)
	assert.Equal(t, len(node.Functions), totalFunctions)
}

func AssertVariable(t *testing.T, node *VariableNode, varType string, id string) {
	assert.Equal(t, node.Type.Identifier, varType)
	assert.Equal(t, node.Identifier, id)
	// TODO Assert Expressions
}

func AssertFunction(t *testing.T, node *FunctionNode, name string, totalRTypes int, totalPTypes int, totalStmts int) {
	assert.Equal(t, node.Name, name)
	assert.Equal(t, len(node.ReturnTypes), totalRTypes)
	assert.Equal(t, len(node.Parameters), totalPTypes)
	assert.Equal(t, len(node.Body), totalStmts)
}

func AssertType(t *testing.T, typeNode *TypeNode, varType string) {
	assert.Equal(t, typeNode.Identifier, varType)
}
