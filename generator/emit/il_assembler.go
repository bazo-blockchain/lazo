package emit

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/il"
	"math/big"
)

type Label uint

/**
 * IL Assembler creates IL Instructions
 */
type ILAssembler struct {
	instructions []*il.Instruction
	targets      map[Label]int
	labelCounter int
	byteCounter  int
}

func NewILAssembler(bytePos int) *ILAssembler {
	return &ILAssembler{
		targets:      map[Label]int{},
		labelCounter: -1,
		byteCounter:  bytePos,
	}
}

func (a *ILAssembler) Complete(halt bool) ([]*il.Instruction, int){
	if halt {
		a.Emit(il.HALT)
	} else {
		a.Emit(il.RET)
	}
	a.ResolveLabels()
	return a.instructions, a.byteCounter
}

func (a *ILAssembler) CreateLabel() Label {
	a.labelCounter++
	return Label(a.labelCounter)
}

func (a *ILAssembler) SetLabel(label Label) {
	a.targets[label] = a.byteCounter
}

func (a *ILAssembler) ResolveLabels() {
	for _, instruction := range a.instructions {
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

// OpCode helpers (Order in the same order as defined)
// --------------------------------------------------

func (a *ILAssembler) PushInt(value *big.Int) {
	var sign byte
	if value.Sign() == -1 {
		sign = 1
	}
	bytes := value.Bytes()
	total := len(bytes)
	operand := append([]byte{byte(total), sign}, bytes...)

	a.addInstruction(il.PUSH, operand)
	a.byteCounter += len(operand) + 1
}

func (a *ILAssembler) PushString(value string) {
	// TODO Implement
}

func (a *ILAssembler) PushCharacter(value rune) {
	// TODO Implement
}

func (a *ILAssembler) NegBool() {
	a.PushInt(big.NewInt(0))
	a.Emit(il.EQ)
}

func (a *ILAssembler) Jmp(label Label) {
	a.addInstruction(il.JMP, label)
	a.byteCounter += 3
}

func (a *ILAssembler) JmpIfTrue(label Label) {
	a.addInstruction(il.JMPIF, label)
	a.byteCounter += 3
}

func (a *ILAssembler) Call(function *symbol.FunctionSymbol) {
	a.addInstruction(il.CALL, function)
	a.byteCounter += 4
}

func (a *ILAssembler) addInstruction(opCode il.OpCode, operand interface{}) {
	a.instructions = append(a.instructions, &il.Instruction{
		OpCode:  opCode,
		Operand: operand,
	})
}
