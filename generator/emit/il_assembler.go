package emit

import (
	"github.com/bazo-blockchain/lazo/generator/il"
)


type Label struct {}

/**
 * IL Assembler creates IL Instructions
 */
type ILAssembler struct {
	function *il.FunctionData
	targets map[*Label]int
}

func NewILAssembler(function *il.FunctionData) *ILAssembler {
	return &ILAssembler{
		function: function,
		targets: map[*Label]int{},
	}
}

func (a *ILAssembler) Complete() {
	// a.Emit(il.RET)
	// resolve labels
}

func (a *ILAssembler) CreateLabel() *Label {
	return &Label{}
}

func (a *ILAssembler) SetLabel(label *Label) {
	a.targets[label] = len(a.function.Instructions)
}

func (a *ILAssembler) FixLabels() {
	for i, instruction := range a.function.Instructions {
		operand := instruction.Operand
		if op, ok := operand.(*Label); ok {
			target := a.targets[op]
			delta := target - i - 1
			instruction.Operand = delta
		}
	}
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
