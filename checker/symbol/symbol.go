package symbol

import "github.com/bazo-blockchain/lazo/parser/node"


// TODO Refactor:
// Remove NewSymbol from the Interface, every Symbol Type can have a NewSymbol method with different params
// as they are called primarly from the phases and in the phases we know what symbol it is... No dynamic calls required...
type Symbol interface {
	NewSymbol(scope Symbol, identifier string) Symbol
	AllDeclarations() []Symbol
	String() string
	GetScope() Symbol
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

func (sym *BaseSymbol) GetScope() Symbol {
	return sym.Scope
}

// Concrete Symbols
//----------------

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
		Symbol: BaseSymbol{}.NewSymbol(nil, ""),
		Types: []*TypeSymbol{},
		Functions: []*FunctionSymbol{},
		Constants: []*ConstantSymbol{},
		NullType: TypeSymbol{}.NewTypeSymbol(cu, "@NULL"),
	}
}

func (cu *CompilationUnit) AllDeclarations() []Symbol {
	// TODO implement
	return nil
}

//----------------

type ContractSymbol struct {
	*TypeSymbol
	Fields []*FieldSymbol
	Functions []*FunctionSymbol
}

func (sym *ContractSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	return &ContractSymbol{
		TypeSymbol: TypeSymbol{}.NewTypeSymbol(scope, identifier),
		Fields: []*FieldSymbol{},
		Functions: []*FunctionSymbol{},
	}
}

func (sym *ContractSymbol) AllDeclarations() []Symbol {
	// TODO Implement
	return nil
}

func (sym *ContractSymbol) AllFields() []*FieldSymbol {
	// TODO Implement
	return nil
}

//----------------

type FieldSymbol struct {
	Symbol
}

func (sym *FieldSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	return &FieldSymbol{
		Symbol: VariableSymbol{}.NewSymbol(scope, identifier),
	}
}

//----------------

type FunctionSymbol struct {
	Symbol
	ReturnTypes []*TypeSymbol
	Parameters []*ParameterSymbol
	LocalVariables []*LocalVariableSymbol

}

func (sym *FunctionSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	if scope == nil {
		// TODO Error
	}
	return &FunctionSymbol{
		Symbol: BaseSymbol{}.NewSymbol(scope, identifier),
	}
}

func (sym *FunctionSymbol) AllDeclarations() []Symbol {
	// TODO implement
	return nil
}

//----------------

type ParameterSymbol struct {
	Symbol
}

func (sym *ParameterSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	return &ParameterSymbol{
		Symbol: VariableSymbol{}.NewSymbol(scope, identifier),
	}
}

//----------------

type LocalVariableSymbol struct {
	Symbol
	VisibleIn map[node.StatementNode]struct{}
}

func (sym *LocalVariableSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	return &LocalVariableSymbol{
		Symbol: VariableSymbol{}.NewSymbol(scope, identifier),
		VisibleIn: map[node.StatementNode]struct{}{},
	}
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