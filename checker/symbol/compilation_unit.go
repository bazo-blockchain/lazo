package symbol

import "fmt"

type CompilationUnit struct {
	Symbol
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

func (cu *CompilationUnit) NewCompilationUnit() Symbol {
	return &CompilationUnit{
		NullType: NewTypeSymbol(cu, "@NULL"),
	}
}

func (cu *CompilationUnit) AllDeclarations() []Symbol {
	// TODO implement
	return nil
}

func (cu *CompilationUnit) String() string {
	return fmt.Sprintf("\n Types: %s"+
		"\n Built-in Types: %s"+
		"\n Constants: %s"+
		"\n Contract: %s", cu.Types, cu.BuiltInTypes, cu.Constants, cu.Contract)
}
