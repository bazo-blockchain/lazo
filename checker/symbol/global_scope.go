package symbol

import "fmt"

type GlobalScope struct {
	AbstractSymbol
	Contract         *ContractSymbol
	Types            []*TypeSymbol
	BuiltInTypes     []*TypeSymbol
	BuiltInFunctions []*FunctionSymbol
	Constants        []*ConstantSymbol

	NullType   *TypeSymbol
	BoolType   *TypeSymbol
	CharType   *TypeSymbol
	StringType *TypeSymbol
	IntType    *TypeSymbol

	TrueConstant  *ConstantSymbol
	FalseConstant *ConstantSymbol
	NullConstant  *ConstantSymbol
}

func newGlobalScope() *GlobalScope {
	gs := &GlobalScope{}
	gs.NullType = NewTypeSymbol(gs, "@NULL")
	return gs
}

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

func (gs *GlobalScope) String() string {
	return fmt.Sprintf("\n Types: %s"+
		"\n Built-in Types: %s"+
		"\n Constants: %s"+
		"\n %s", gs.Types, gs.BuiltInTypes, gs.Constants, gs.Contract)
}
