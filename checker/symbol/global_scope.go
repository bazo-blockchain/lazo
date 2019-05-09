package symbol

import "fmt"

// GlobalScope encapsulates global information such as the contract, types, built-ins, constants a.s.o.
// It is used to lookup global information
type GlobalScope struct {
	AbstractSymbol
	Contract         *ContractSymbol
	Types            map[string]TypeSymbol
	BuiltInTypes     []*BasicTypeSymbol
	BuiltInFunctions []*FunctionSymbol
	Constants        []*ConstantSymbol
	Structs          map[string]*StructTypeSymbol

	NullType   *BasicTypeSymbol
	BoolType   *BasicTypeSymbol
	CharType   *BasicTypeSymbol
	StringType *BasicTypeSymbol
	IntType    *BasicTypeSymbol

	TrueConstant  *ConstantSymbol
	FalseConstant *ConstantSymbol
	NullConstant  *ConstantSymbol
	ArrayLength   *FieldSymbol
}

func newGlobalScope() *GlobalScope {
	gs := &GlobalScope{}
	gs.NullType = NewBasicTypeSymbol(gs, "@NULL")
	gs.Structs = make(map[string]*StructTypeSymbol)
	gs.Types = make(map[string]TypeSymbol)
	return gs
}

// AllDeclarations returns all declarations made within the global scope such as types, built-ins and constants
func (gs *GlobalScope) AllDeclarations() []Symbol {
	var symbols []Symbol
	for _, s := range gs.Types {
		symbols = append(symbols, s)
	}
	for _, s := range gs.BuiltInFunctions {
		symbols = append(symbols, s)
	}
	for _, s := range gs.Constants {
		symbols = append(symbols, s)
	}
	return symbols
}

// String creates a string representation of the global scope
func (gs *GlobalScope) String() string {
	return fmt.Sprintf("\n Types: %s"+
		"\n Built-in Types: %s"+
		"\n Constants: %s"+
		"\n %s", gs.Types, gs.BuiltInTypes, gs.Constants, gs.Contract)
}
