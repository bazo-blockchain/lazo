package il

// OpCode is the type of byte code supported on Bazo VM
type OpCode byte

// Supported OpCodes. The OpCodes are copied from Bazo VM.
// See https://github.com/bazo-blockchain/bazo-vm/blob/master/vm/op_codes.go
const (
	PushInt OpCode = iota
	PushBool
	PushChar
	PushStr
	Push
	Dup
	Roll
	Swap
	Pop
	Add
	Sub
	Mul
	Div
	Mod
	Exp
	Neg
	Eq
	NotEq
	Lt
	Gt
	LtEq
	GtEq
	ShiftL
	ShiftR
	NoOp
	Jmp
	JmpTrue
	JmpFalse
	Call
	CallTrue
	CallExt
	Ret
	Size
	StoreLoc
	StoreSt
	LoadLoc
	LoadSt
	Address // Address of account
	Issuer  // Owner of smart contract account
	Balance // Balance of account
	Caller
	CallVal  // Amount of bazo coins transacted in transaction
	CallData // Parameters and function signature hash
	NewMap
	MapHasKey
	MapGetVal
	MapSetVal
	MapRemove
	NewArr
	ArrAppend
	ArrInsert
	ArrRemove
	ArrAt
	ArrLen
	NewStr
	StoreFld
	LoadFld
	SHA3
	CheckSig
	ErrHalt
	Halt
)
