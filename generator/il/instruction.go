package il

type OpCode int

const (
	TEST OpCode = iota
)

type Instruction struct {
	OpCode OpCode
	Operand struct{}
}
