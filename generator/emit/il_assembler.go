package emit

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/il"
	"github.com/bazo-blockchain/lazo/generator/util"
	"math/big"
)

type Label uint

/**
 * IL Assembler creates IL Instructions
 */
type ILAssembler struct {
	instructions []*il.Instruction
	targets      map[Label]uint16
	labelCounter int
	bytePos      *uint16
}

func NewILAssembler(bytePos *uint16) *ILAssembler {
	return &ILAssembler{
		targets:      map[Label]uint16{},
		labelCounter: -1,
		bytePos:      bytePos,
	}
}

func (a *ILAssembler) Complete(halt bool) []*il.Instruction {
	if halt {
		a.Emit(il.HALT)
	} else {
		a.Emit(il.RET)
	}
	a.ResolveLabels()
	return a.instructions
}

func (a *ILAssembler) CreateLabel() Label {
	a.labelCounter++
	return Label(a.labelCounter)
}

func (a *ILAssembler) SetLabel(label Label) {
	a.targets[label] = *a.bytePos
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
	a.addInstruction(opCode, nil, 0)
}

// OpCode helpers (Order in the same order as defined)
// --------------------------------------------------

func (a *ILAssembler) PushInt(value *big.Int) {
	sign := util.GetSignByte(value)
	bytes := value.Bytes()
	total := byte(len(bytes))
	operand := append([]byte{total, sign}, bytes...)

	a.addInstruction(il.PUSH, operand, byte(len(operand)))
}

func (a *ILAssembler) PushBool(value bool) {
	if value {
		a.PushInt(big.NewInt(1))
	} else {
		a.PushInt(big.NewInt(0))
	}
}

func (a *ILAssembler) PushString(value string) {
	// TODO Implement
}

func (a *ILAssembler) PushCharacter(value rune) {
	// TODO Implement
}

func (a *ILAssembler) NegBool() {
	a.PushBool(false)
	a.Emit(il.EQ)
}

func (a *ILAssembler) ConvertToBool() {
	a.PushBool(true)
	a.Emit(il.EQ)
}

func (a *ILAssembler) Jmp(label Label) {
	a.addInstruction(il.JMP, label, 2)
}

func (a *ILAssembler) JmpIfTrue(label Label) {
	a.addInstruction(il.JMPIF, label, 2)
}

func (a *ILAssembler) Call(function *symbol.FunctionSymbol) {
	a.addInstruction(il.CALL, function, 3)
}

func (a *ILAssembler) Store(index byte) {
	a.addInstruction(il.STORE, []byte{index}, 1)
}

func (a *ILAssembler) Load(index byte) {
	a.addInstruction(il.LOAD, []byte{index}, 1)
}

func (a *ILAssembler) addInstruction(opCode il.OpCode, operand interface{}, operandSize byte) {
	a.instructions = append(a.instructions, &il.Instruction{
		OpCode:  opCode,
		Operand: operand,
	})
	*a.bytePos += uint16(operandSize) + 1
}
