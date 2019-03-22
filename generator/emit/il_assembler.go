package emit

import (
	"github.com/bazo-blockchain/lazo/generator/il"
	"math/big"
)

type Label uint

/**
 * IL Assembler creates IL Instructions
 */
type ILAssembler struct {
	function     *il.FunctionData
	targets      map[Label]int
	labelCounter int
	byteCounter  int
}

func NewILAssembler(function *il.FunctionData) *ILAssembler {
	return &ILAssembler{
		function:     function,
		targets:      map[Label]int{},
		labelCounter: -1,
		byteCounter:  0,
	}
}

func (a *ILAssembler) Complete() {
	a.Emit(il.RET)
	a.ResolveLabels()
}

func (a *ILAssembler) CreateLabel() Label {
	a.labelCounter++
	return Label(a.labelCounter)
}

func (a *ILAssembler) SetLabel(label Label) {
	a.targets[label] = a.byteCounter
}

func (a *ILAssembler) ResolveLabels() {
	for _, instruction := range a.function.Instructions {
		operand := instruction.Operand
		if op, ok := operand.(Label); ok {
			instruction.Operand = []byte{0, byte(a.targets[op])}
		}
	}
}

func (a *ILAssembler) Emit(opCode il.OpCode) {
	a.addInstruction(opCode, nil)
	a.byteCounter++
}

func (a *ILAssembler) PushInt(value *big.Int) {
	var sign byte
	if value.Sign() == -1 {
		sign = 1
	}
	bytes := value.Bytes()
	total := len(bytes)
	a.addInstruction(il.PUSH, append([]byte{byte(total), sign}, bytes...))
	a.byteCounter += 3 + total
}

func (a *ILAssembler) Jmp(label Label) {
	a.addInstruction(il.JMP, label)
	a.byteCounter += 3
}

func (a *ILAssembler) JmpIfTrue(label Label) {
	a.addInstruction(il.JMPIF, label)
	a.byteCounter += 3
}

func (a *ILAssembler) NegBool() {
	a.PushInt(big.NewInt(0))
	a.Emit(il.EQ)
}

// TODO: Make it a private function
// Remove direct access to this function in the visitor. Use Emit or explicit functions
func (a *ILAssembler) addInstruction(opCode il.OpCode, operand interface{}) {
	a.function.Instructions = append(a.function.Instructions, &il.Instruction{
		OpCode:  opCode,
		Operand: operand,
	})
}
