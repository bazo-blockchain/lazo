package il

type TypeData int

const (
	_ = iota
	BoolType TypeData = iota * -1
	IntType
	CharType
	StringType
)

type FunctionData struct {
	Identifier string `json:"ID"`
	ReturnTypes []TypeData `json:"ReturnType"`
	ParamTypes []TypeData `json:"ParamType"`
	LocalTypes []TypeData `json:"LocalType"`
	Code []*Instruction `json:"Instruction"`
}

type ContractData struct {
	Identifier string `json:"ID"`
	Fields []TypeData `json:"Fields"`
	Functions []TypeData `json:"Functions"`
}
