package symbol

import "fmt"

type CompilationUnit struct {
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

func newCompilationUnit() *CompilationUnit {
	cu := &CompilationUnit{}
	cu.NullType = NewTypeSymbol(cu, "@NULL")
	return cu
}

func (cu *CompilationUnit) AllDeclarations() []Symbol {
	var symbols []Symbol
	for _, s := range cu.Types {
		symbols = append(symbols, s)
	}
	for _, s := range cu.BuiltInFunctions {
		symbols = append(symbols, s)
	}
	for _, s := range cu.Constants {
		symbols = append(symbols, s)
	}
	return symbols
}

func (cu *CompilationUnit) String() string {
	return fmt.Sprintf("\n Types: %s"+
		"\n Built-in Types: %s"+
		"\n Constants: %s"+
		"\n %s", cu.Types, cu.BuiltInTypes, cu.Constants, cu.Contract)
}
