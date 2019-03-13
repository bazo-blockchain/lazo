package symbol

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type Symbol interface {
	GetScope() Symbol
	GetIdentifier() string
	AllDeclarations() []Symbol
	String() string
}

type AbstractSymbol struct {
	Scope      Symbol
	Identifier string
}

func NewAbstractSymbol(scope Symbol, identifier string) AbstractSymbol {
	return AbstractSymbol{
		Scope:      scope,
		Identifier: identifier,
	}
}

func (sym *AbstractSymbol) GetScope() Symbol {
	return sym.Scope
}

func (sym *AbstractSymbol) GetIdentifier() string {
	return sym.Identifier
}

func (sym *AbstractSymbol) AllDeclarations() []Symbol {
	return []Symbol{}
}

func (sym *AbstractSymbol) String() string {
	return fmt.Sprintf("Abstract Symbol: %s", sym.Identifier)
}

// Concrete Symbols
//-----------------

type ContractSymbol struct {
	AbstractSymbol
	Fields    []*FieldSymbol
	Functions []*FunctionSymbol
}

func NewContractSymbol(scope Symbol, identifier string) *ContractSymbol {
	return &ContractSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

func (sym *ContractSymbol) AllDeclarations() []Symbol {
	// TODO Implement
	return nil
}

func (sym *ContractSymbol) String() string {
	return fmt.Sprintf("Contract: %s, Fields: %s, Functions %s", sym.Identifier, sym.Fields, sym.Functions)
}

//----------------

type FieldSymbol struct {
	AbstractSymbol
	Type *TypeSymbol
}

func NewFieldSymbol(scope Symbol, identifier string) *FieldSymbol {
	return &FieldSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

func (sym *FieldSymbol) String() string {
	return fmt.Sprintf("%s %s", sym.Type, sym.Identifier)
}

//----------------

type FunctionSymbol struct {
	AbstractSymbol
	ReturnTypes    []*TypeSymbol
	Parameters     []*ParameterSymbol
	LocalVariables []*LocalVariableSymbol
}

func NewFunctionSymbol(scope Symbol, identifier string) *FunctionSymbol {
	return &FunctionSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

func (sym *FunctionSymbol) AllDeclarations() []Symbol {
	// TODO implement
	return nil
}

func (sym *FunctionSymbol) String() string {
	return fmt.Sprintf("%s %s(%s): %d vars", sym.ReturnTypes, sym.Identifier, sym.Parameters, len(sym.LocalVariables))
}

//----------------

type ParameterSymbol struct {
	VariableSymbol
}

func NewParameterSymbol(scope Symbol, identifier string) *ParameterSymbol {
	return &ParameterSymbol{
		VariableSymbol: *NewVariableSymbol(scope, identifier),
	}
}

//----------------

type LocalVariableSymbol struct {
	VariableSymbol
	VisibleIn []*node.StatementNode
}

func NewLocalVariableSymbol(scope Symbol, identifier string) *LocalVariableSymbol {
	return &LocalVariableSymbol{
		VariableSymbol: *NewVariableSymbol(scope, identifier),
	}
}

//----------------

type VariableSymbol struct {
	AbstractSymbol
	Type *TypeSymbol
}

func NewVariableSymbol(scope Symbol, identifier string) *VariableSymbol {
	return &VariableSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

func (sym *VariableSymbol) String() string {
	return fmt.Sprintf("%s %s", sym.Type, sym.Identifier)
}

//----------------

type ConstantSymbol struct {
	AbstractSymbol
	Type *TypeSymbol
}

func NewConstantSymbol(scope Symbol, identifier string, typeSymbol *TypeSymbol) *ConstantSymbol {
	return &ConstantSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
		Type:           typeSymbol,
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
