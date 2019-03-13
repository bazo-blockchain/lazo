package symbol

type CompilationUnit struct {
	Symbol
	Contract *ContractSymbol
	Types []*TypeSymbol
	Functions []*FunctionSymbol
	Constants []*ConstantSymbol

	NullType *TypeSymbol
	BoolType *TypeSymbol
	CharType *TypeSymbol
	StringType *TypeSymbol
	IntType *TypeSymbol

	BuiltInTypes []*TypeSymbol

	TrueConstant *ConstantSymbol
	FalseConstant *ConstantSymbol
	NullConstant *ConstantSymbol
}

func (cu *CompilationUnit) NewCompilationUnit() Symbol{
	return &CompilationUnit{
		Types: []*TypeSymbol{},
		Functions: []*FunctionSymbol{},
		Constants: []*ConstantSymbol{},
		NullType: (&TypeSymbol{}).NewTypeSymbol(cu, "@NULL"),
	}
}

func (cu *CompilationUnit) AllDeclarations() []Symbol {
	// TODO implement
	return nil
}
