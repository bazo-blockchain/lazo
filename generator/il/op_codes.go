package il

type OpCode int

// Opcodes have been adopted from the VM
// See https://github.com/bazo-blockchain/bazo-smartcontract/blob/master/src/vm/op_codes.go
const (
	PUSH OpCode = iota
	DUP
	ROLL
	POP
	ADD
	SUB
	MULT
	DIV
	MOD
	NEG
	EQ
	NEQ
	LT
	GT
	LTE
	GTE
	SHIFTL
	SHIFTR
	NOP
	JMP
	JMPIF
	CALL
	CALLIF
	CALLEXT
	RET
	SIZE
	STORE
	SSTORE
	LOAD
	SLOAD
	ADDRESS
	ISSUER
	BALANCE
	CALLER
	CALLVAL
	CALLDATA
	NEWMAP
	MAPHASKEY
	MAPPUSH
	MAPGETVAL
	MAPSETVAL
	MAPPREMOVE
	NEWARR
	ARRAPPEND
	ARRINSERT
	ARRREMOVE
	ARRAT
	SHA3
	CHECKSIG
	ERRHALT
	HALT
)

var OpCodeLiterals = map[OpCode]string{
	PUSH:       "PUSH",
	DUP:        "DUP",
	ROLL:       "ROLL",
	POP:        "POP",
	ADD:        "ADD",
	SUB:        "SUB",
	MULT:       "MULT",
	DIV:        "DIV",
	MOD:        "MOD",
	NEG:        "NEG",
	EQ:         "EQ",
	NEQ:        "NEQ",
	LT:         "LT",
	GT:         "GT",
	LTE:        "LTE",
	GTE:        "GTE",
	SHIFTL:     "SHIFTL",
	SHIFTR:     "SHIFTR",
	NOP:        "NOP",
	JMP:        "JMP",
	JMPIF:      "JMPIF",
	CALL:       "CALL",
	CALLIF:     "CALLIF",
	CALLEXT:    "CALLEXT",
	RET:        "RET",
	SIZE:       "SIZE",
	STORE:      "STORE",
	SSTORE:     "SSTORE",
	LOAD:       "LOAD",
	SLOAD:      "SLOAD",
	ADDRESS:    "ADDRESS",
	ISSUER:     "ISSUER",
	BALANCE:    "BALANCE",
	CALLER:     "CALLER",
	CALLVAL:    "CALLVAL",
	CALLDATA:   "CALLDATA",
	NEWMAP:     "NEWMAP",
	MAPHASKEY:  "MAPHASKEY",
	MAPPUSH:    "MAPPUSH",
	MAPGETVAL:  "MAPGETVAL",
	MAPSETVAL:  "MAPSETVAL",
	MAPPREMOVE: "MAPPREMOVE",
	NEWARR:     "NEWARR",
	ARRAPPEND:  "ARRAPPEND",
	ARRINSERT:  "ARRINSERT",
	ARRREMOVE:  "ARRREMOVE",
	ARRAT:      "ARRAT",
	SHA3:       "SHA3",
	CHECKSIG:   "CHECKSIG",
	ERRHALT:    "ERRHALT",
	HALT:       "HALT",
}
