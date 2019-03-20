package emit

import (
	"github.com/bazo-blockchain/lazo/generator/il"
)

/**
 * IL Assembler creates IL Instructions
 */
type ILAssembler struct {
	function *il.FunctionData
}

func NewILAssembler(function *il.FunctionData) *ILAssembler {
	return &ILAssembler{
		function: function,
	}
}

func (a *ILAssembler) Complete() {
	// a.Emit(il.RET)
	// resolve labels
}

func (a *ILAssembler) Emit(opCode il.OpCode) {
	a.EmitOperand(opCode, nil)
}

func (a *ILAssembler) EmitOperand(opCode il.OpCode, operand interface{}) {
	a.function.Instructions = append(a.function.Instructions, &il.Instruction{
		OpCode:  opCode,
		Operand: operand,
	})
}
