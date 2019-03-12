package symbol

type Symbol interface {
	NewSymbol(scope Symbol, identifier string) Symbol
	AllDeclarations() []Symbol
	String() string
}

type BaseSymbol struct {
	Scope Symbol
	Identifier string
}

func (sym *BaseSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	return &BaseSymbol{
		Scope: scope,
		Identifier: identifier,
	}
}

func (sym *BaseSymbol) AllDeclarations() []Symbol {
	return []Symbol{}
}

func (sym *BaseSymbol) String() string {
	return sym.Identifier
}

// Concrete Symbols
//----------------

type CompilationUnit struct {
	Symbol
	Contract *ContractSymbol
	Types *[]TypeSymbol
	Functions *[]FunctionSymbol
	Constants *[]ConstantSymbol

	NullType *TypeSymbol
	BoolType *TypeSymbol
	CharType *TypeSymbol
	StringType *TypeSymbol
	IntType *TypeSymbol

	BuiltInTypes *[]TypeSymbol

	TrueConstant *ConstantSymbol
	FalseConstant *ConstantSymbol
	NullConstant *ConstantSymbol
}

func (cu *CompilationUnit) NewCompilationUnit() Symbol{
	return &CompilationUnit{
		Symbol: BaseSymbol{}.NewSymbol(nil, ""),
		NullType: TypeSymbol{}.NewTypeSymbol(cu, "@NULL"),
	}
}

func (cu *CompilationUnit) AllDeclarations() []Symbol {
	// TODO implement
	return nil
}

//----------------

type ContractSymbol struct {
	Symbol
	Fields *[]FieldSymbol
	Functions *[]FunctionSymbol
}

func (sym *ContractSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	return &ContractSymbol{
		Symbol: TypeSymbol{}.NewSymbol(scope, identifier),
	}
}

func (sym *ContractSymbol) AllDeclarations() []Symbol {
	// TODO Implement
	return nil
}

func (sym *ContractSymbol) AllFields() *[]FieldSymbol {
	// TODO Implement
	return nil
}

//----------------

type FieldSymbol struct {
	// TODO Implement
}

func (sym *FieldSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	// TODO Implement
	return nil
}

//----------------

type FunctionSymbol struct {
	// TODO Implement
}

func (sym *FunctionSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	// TODO Implement
	return nil
}

//----------------

type ParameterSymbol struct {
	// TODO Implement
}

func (sym *ParameterSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	// TODO Implement
	return nil
}

//----------------

type LocalVariableSymbol struct {
	// TODO Implement
}

func (sym *LocalVariableSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	// TODO Implement
	return nil
}

//----------------

type VariableSymbol struct {
	Symbol
	Type *TypeSymbol
}

func (sym *VariableSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	if scope == nil {
		// TODO Error
	}
	return &VariableSymbol{
		Symbol: BaseSymbol{}.NewSymbol(scope, identifier),
	}
}

//----------------

type ConstantSymbol struct {
	Symbol
	Type *TypeSymbol
}

func (sym *ConstantSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	if scope == nil {
		// TODO Error
	}
	if identifier == "" {
		// TODO Error
	}

	return &ConstantSymbol{
		Symbol: BaseSymbol{}.NewSymbol(scope, identifier),
	}
}

func (sym *ConstantSymbol) NewConstantSymbol(scope Symbol, identifier string, typeSymbol *TypeSymbol) *ConstantSymbol {
	constantSymbol := sym.NewSymbol(scope, identifier)
	if typeSymbol == nil {
		// TODO Error
	}
	symbol, _ := constantSymbol.(*ConstantSymbol)
	symbol.Type = typeSymbol
	return symbol
}

//----------------

type TypeSymbol struct {
	Symbol
}

func (sym *TypeSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	if scope == nil {
		// TODO Error
	}
	if identifier == "" {
		// TODO Error
	}
	return &TypeSymbol{
		Symbol: BaseSymbol{}.NewSymbol(scope, identifier),
	}
}

func (sym *TypeSymbol) NewTypeSymbol(scope Symbol, identifier string) *TypeSymbol {
	baseSymbol := sym.NewSymbol(scope, identifier)
	symbol, _ := baseSymbol.(*TypeSymbol)
	return symbol
}

//----------------