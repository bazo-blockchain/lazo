package emit

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/generator/il"
	"math/big"
)

type Label struct{}

/**
 * IL Assembler creates IL Instructions
 */
type ILAssembler struct {
	function    *il.FunctionData
	targets     map[*Label]int
	byteCounter int
}

func NewILAssembler(function *il.FunctionData) *ILAssembler {
	return &ILAssembler{
		function:    function,
		targets:     map[*Label]int{},
		byteCounter: 0,
	}
}

func (a *ILAssembler) Complete() {
	a.Emit(il.RET)
	a.ResolveLabels()
}

func (a *ILAssembler) CreateLabel() *Label {
	return &Label{}
}

func (a *ILAssembler) SetLabel(label *Label) {
	a.targets[label] = a.byteCounter
}

func (a *ILAssembler) ResolveLabels() {
	for _, instruction := range a.function.Instructions {
		operand := instruction.Operand
		if op, ok := operand.(*Label); ok {
			fmt.Println(a.targets[op])
			instruction.Operand = []byte{0, byte(a.targets[op])}
		}
	}
}

func (a *ILAssembler) Emit(opCode il.OpCode) {
	a.EmitOperand(opCode, nil)
	a.byteCounter++
}

func (a *ILAssembler) EmitOperand(opCode il.OpCode, operand interface{}) {
	a.function.Instructions = append(a.function.Instructions, &il.Instruction{
		OpCode:  opCode,
		Operand: operand,
	})
}

func (a *ILAssembler) PushInt(value *big.Int) {
	var sign byte
	if value.Sign() == -1 {
		sign = 1
	}
	bytes := value.Bytes()
	total := len(bytes)
	a.EmitOperand(il.PUSH, append([]byte{byte(total), sign}, bytes...))
	a.byteCounter += 3 + total
}

func (a *ILAssembler) JmpIf(label *Label) {
	a.EmitOperand(il.JMPIF, label)
	a.byteCounter += 3
}
