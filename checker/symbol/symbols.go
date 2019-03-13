package symbol

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type Symbol interface {
	GetScope() Symbol
	GetName() string
	AllDeclarations() []Symbol
	String() string
}

type AbstractSymbol struct {
	Scope Symbol
	Identifier string
}

func NewAbstractSymbol(scope Symbol, identifier string) AbstractSymbol {
	return AbstractSymbol{
		Scope: scope,
		Identifier: identifier,
	}
}

func (sym *AbstractSymbol) GetScope() Symbol {
	return sym.Scope
}

func (sym *AbstractSymbol) GetName() string {
	return sym.Identifier
}

func (sym *AbstractSymbol) AllDeclarations() []Symbol {
	return []Symbol{}
}

func (sym *AbstractSymbol) String() string {
	return fmt.Sprintf("Abstract Symbol: %s", sym.GetName())
}

// Concrete Symbols
//----------------

type ContractSymbol struct {
	*TypeSymbol
	Fields []*FieldSymbol
	Functions []*FunctionSymbol
}

func (sym *ContractSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	return &ContractSymbol{
		TypeSymbol: NewTypeSymbol(scope, identifier),
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
		Symbol: (&VariableSymbol{}).NewSymbol(scope, identifier),
	}
}

//----------------

type FunctionSymbol struct {
	AbstractSymbol
	ReturnTypes []*TypeSymbol
	Parameters []*ParameterSymbol
	LocalVariables []*LocalVariableSymbol

}

func (sym *FunctionSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	if scope == nil {
		// TODO Error
	}
	return &FunctionSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
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
		Symbol: (&VariableSymbol{}).NewSymbol(scope, identifier),
	}
}

//----------------

type LocalVariableSymbol struct {
	Symbol
	VisibleIn map[node.StatementNode]struct{}
}

func (sym *LocalVariableSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	return &LocalVariableSymbol{
		Symbol: (&VariableSymbol{}).NewSymbol(scope, identifier),
		VisibleIn: map[node.StatementNode]struct{}{},
	}
}

//----------------

type VariableSymbol struct {
	AbstractSymbol
	Type *TypeSymbol
}

func (sym *VariableSymbol) NewSymbol(scope Symbol, identifier string) Symbol {
	if scope == nil {
		// TODO Error
	}
	return &VariableSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

//----------------

type ConstantSymbol struct {
	AbstractSymbol
	Type *TypeSymbol
}

func NewConstantSymbol(scope Symbol, identifier string, typeSymbol *TypeSymbol) *ConstantSymbol {
	return &ConstantSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
		Type: typeSymbol,
	}
}

func (sym *ConstantSymbol) String() string {
	return fmt.Sprintf("Constant %s", sym.Identifier)
}

//----------------

type TypeSymbol struct {
	AbstractSymbol
}

func NewTypeSymbol(scope Symbol, identifier string) *TypeSymbol {
	return &TypeSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

func (sym *TypeSymbol) String() string {
	return fmt.Sprintf("Type %s", sym.Identifier)
}
