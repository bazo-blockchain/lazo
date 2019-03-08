package symbol

type Symbol interface {
	NewSymbol() Symbol
	AllDeclarations() []Symbol
	String() string
}

type BaseSymbol struct {
	Scope Symbol
	Identifier string
}

func (sym *BaseSymbol) NewSymbol(scope Symbol, identifier string) BaseSymbol {
	return BaseSymbol{
		Scope: scope,
		Identifier: identifier,
	}
}

func (sym *BaseSymbol) AllDeclarations() []BaseSymbol {
	return []BaseSymbol{}
}

func (sym *BaseSymbol) String() string {
	return sym.Identifier
}

// TODO Implement further symbols
