package emit

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/il"
	"github.com/bazo-blockchain/lazo/generator/util"
	"math/big"
)

// Label contain an unsigned int. They are used to jump to certain points within the byte code.
type Label uint

// ILAssembler creates IL instructions
type ILAssembler struct {
	instructions []*il.Instruction
	targets      map[Label]uint16
	labelCounter int
	bytePos      *uint16
}

// NewILAssembler creates a new ILAssembler
func NewILAssembler(bytePos *uint16) *ILAssembler {
	return &ILAssembler{
		targets:      map[Label]uint16{},
		labelCounter: -1,
		bytePos:      bytePos,
	}
}

// Complete appends a HALT or RET instruction to the byte code and resolves the labels.
// Returns the instructions
func (a *ILAssembler) Complete(halt bool) []*il.Instruction {
	if halt {
		a.Emit(il.Halt)
	} else {
		a.Emit(il.Ret)
	}
	a.ResolveLabels()
	return a.instructions
}

// CreateLabel increases the label counter and creates a new label
func (a *ILAssembler) CreateLabel() Label {
	a.labelCounter++
	return Label(a.labelCounter)
}

// SetLabel sets the label to a certain position in the byte code
func (a *ILAssembler) SetLabel(label Label) {
	a.targets[label] = *a.bytePos
}

// ResolveLabels resolves the label and sets the instruction operand
func (a *ILAssembler) ResolveLabels() {
	for _, instruction := range a.instructions {
		operand := instruction.Operand
		if op, ok := operand.(Label); ok {
			instruction.Operand = util.GetBytesFromUInt16(a.targets[op])
		}
	}
}

// Emit adds a new instruction to the byte code
func (a *ILAssembler) Emit(opCode il.OpCode) {
	a.addInstruction(opCode, nil, 0)
}

// OpCode helpers (Order in the same order as defined)
// --------------------------------------------------

// PushInt is a helper function that emits byte code to push an integer to the stack.
func (a *ILAssembler) PushInt(value *big.Int) {
	sign := util.GetSignByte(value)
	bytes := value.Bytes()
	total := byte(len(bytes))

	var operand []byte
	if total == 0 {
		operand = []byte{0}
	} else {
		operand = append([]byte{total, sign}, bytes...)
	}

	a.addInstruction(il.PushInt, operand, byte(len(operand)))
}

// PushBool is a helper function that emits byte code to push a boolean to the stack.
func (a *ILAssembler) PushBool(value bool) {
	var byteVal byte
	if value {
		byteVal = 1
	} else {
		byteVal = 0
	}
	a.addInstruction(il.PushBool, []byte{byteVal}, 1)
}

// PushString is a helper that emits byte code to push a string to the stack
func (a *ILAssembler) PushString(value string) {
	bytes := []byte(value)
	total := byte(len(bytes))
	operand := append([]byte{total}, bytes...)
	a.addInstruction(il.PushStr, operand, byte(len(operand)))
}

// PushCharacter is a helper that emits byte code to push a character to the stack
func (a *ILAssembler) PushCharacter(value rune) {
	operand := []byte(string(value))
	a.addInstruction(il.PushChar, operand, 1)
}

// PushFuncHash is a helper that emits byte code to push a function hash to the stack
func (a *ILAssembler) PushFuncHash(hash [4]byte) {
	var operand = make([]byte, 5)
	operand[0] = 4
	copy(operand[1:], hash[:])

	a.addInstruction(il.Push, operand, byte(len(operand)))
}

// Jmp is a helper that adds a JMP instruction to the byte code
// Is used to jump to labels within the byte code
func (a *ILAssembler) Jmp(label Label) {
	a.addInstruction(il.Jmp, label, 2)
}

// JmpTrue is a helper that adds a JmpTrue instruction to the byte code
// Is used to jump to labels within the byte code if the value at the top of the stack is 1 (true)
func (a *ILAssembler) JmpTrue(label Label) {
	a.addInstruction(il.JmpTrue, label, 2)
}

// JmpFalse is a helper that adds a JmpFalse instruction to the byte code
// Is used to jump to labels within the byte code if the value at the top of the stack is 0 (false)
func (a *ILAssembler) JmpFalse(label Label) {
	a.addInstruction(il.JmpFalse, label, 2)
}

// CallFunc is a helper that adds a CALL instruction to the byte code
// Is used to call functions
func (a *ILAssembler) CallFunc(function *symbol.FunctionSymbol) {
	a.addInstruction(il.Call, function, 4)
}

// CallTrue is a helper that adds a CALLIF instruction to the byte code
// Is used to call functions if the value at the top of the stack is 1 (true)
func (a *ILAssembler) CallTrue(function *symbol.FunctionSymbol) {
	a.addInstruction(il.CallTrue, function, 4)
}

// StoreLocal is a helper that adds a STORE instruction to the byte code
// Is used to store the value at the top of the stack in the call stack.
func (a *ILAssembler) StoreLocal(index byte) {
	a.addInstruction(il.StoreLoc, []byte{index}, 1)
}

// StoreState is a helper that adds a SSTORE instruction to the byte code
// Is used to store the value at the top of the stack within the contract variables (fields) at a certain index
func (a *ILAssembler) StoreState(index byte) {
	a.addInstruction(il.StoreSt, []byte{index}, 1)
}

// LoadLocal is a helper that adds a LOAD instruction to the byte code
// Is used to load the variable at the given index from the call stack to the evaluation stack
func (a *ILAssembler) LoadLocal(index byte) {
	a.addInstruction(il.LoadLoc, []byte{index}, 1)
}

// LoadState is a helper that adds a SLOAD instruction to the byte code
// Is used to load the variable at the given index from the contract variables (fields) to the evaluation stack
func (a *ILAssembler) LoadState(index byte) {
	a.addInstruction(il.LoadSt, []byte{index}, 1)
}

// NewStruct creates a new struct with the given number of fields
func (a *ILAssembler) NewStruct(totalFields uint16) {
	a.addInstruction(il.NewStr, util.GetBytesFromUInt16(totalFields), 2)
}

// StoreField pops element and struct from evaluation stack and stores the element at the given field index.
func (a *ILAssembler) StoreField(index uint16) {
	a.addInstruction(il.StoreFld, util.GetBytesFromUInt16(index), 2)
}

func (a *ILAssembler) addInstruction(opCode il.OpCode, operand interface{}, operandSize byte) {
	a.instructions = append(a.instructions, &il.Instruction{
		OpCode:  opCode,
		Operand: operand,
	})
	*a.bytePos += uint16(operandSize) + 1
}
