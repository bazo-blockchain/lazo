package il

// Instruction consists of an OpCode and the Operand
type Instruction struct {
	OpCode  OpCode
	Operand interface{}
}
