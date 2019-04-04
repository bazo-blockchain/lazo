package symbol

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/parser/node"
)

// Symbol declares functions which symbols have to implement
type Symbol interface {
	GetScope() Symbol
	GetIdentifier() string
	AllDeclarations() []Symbol
	String() string
}

// AbstractSymbol is part of all symbols and contains the scope and identifier
type AbstractSymbol struct {
	Scope      Symbol
	Identifier string
}

// NewAbstractSymbol is a helper function that creates a new abstract symbol
func NewAbstractSymbol(scope Symbol, identifier string) AbstractSymbol {
	return AbstractSymbol{
		Scope:      scope,
		Identifier: identifier,
	}
}

// GetScope returns the scope of the symbol (will be deleted)
func (sym *AbstractSymbol) GetScope() Symbol {
	return sym.Scope
}

// GetIdentifier returns the identifier of the identifier (will be deleted)
func (sym *AbstractSymbol) GetIdentifier() string {
	return sym.Identifier
}

// AllDeclarations returns an empty symbol slice
func (sym *AbstractSymbol) AllDeclarations() []Symbol {
	return []Symbol{}
}

// String creates the string representation for the abstract symbol
func (sym *AbstractSymbol) String() string {
	return fmt.Sprintf("Abstract Symbol: %s", sym.Identifier)
}

// Concrete Symbols
//-----------------

// ContractSymbol contains fields and functions
type ContractSymbol struct {
	AbstractSymbol
	Fields    []*FieldSymbol
	Functions []*FunctionSymbol
}

// NewContractSymbol creates a new ContractSymbol
func NewContractSymbol(scope Symbol, identifier string) *ContractSymbol {
	return &ContractSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// AllDeclarations returns all field and function declarations
func (sym *ContractSymbol) AllDeclarations() []Symbol {
	var symbols []Symbol
	for _, s := range sym.Fields {
		symbols = append(symbols, s)
	}
	for _, s := range sym.Functions {
		symbols = append(symbols, s)
	}
	return symbols
}

// GetFieldIndex returns the index of the field
func (sym *ContractSymbol) GetFieldIndex(id string) int {
	for i, s := range sym.Fields {
		if s.GetIdentifier() == id {
			return i
		}
	}
	return -1
}

// String creates the string representation for the ContractSymbol
func (sym *ContractSymbol) String() string {
	return fmt.Sprintf("Contract: %s, \nFields: %s, \nFunctions %s", sym.Identifier, sym.Fields, sym.Functions)
}

//----------------

// FieldSymbol contains the type of the field
type FieldSymbol struct {
	AbstractSymbol
	Type *TypeSymbol
}

// NewFieldSymbol creates a new FieldSymbol
func NewFieldSymbol(scope Symbol, identifier string) *FieldSymbol {
	return &FieldSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// String creates the string representation of the FieldSymbol
func (sym *FieldSymbol) String() string {
	return fmt.Sprintf("%s %s", sym.Type, sym.Identifier)
}

//----------------

// FunctionSymbol contains the return types, parameters and local variables of a function
type FunctionSymbol struct {
	AbstractSymbol
	ReturnTypes    []*TypeSymbol
	Parameters     []*ParameterSymbol
	LocalVariables []*LocalVariableSymbol
}

// NewFunctionSymbol creates a new FunctionSymbol
func NewFunctionSymbol(scope Symbol, identifier string) *FunctionSymbol {
	return &FunctionSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// AllDeclarations returns all parameter and local variable declarations
func (sym *FunctionSymbol) AllDeclarations() []Symbol {
	var symbols []Symbol
	for _, s := range sym.Parameters {
		symbols = append(symbols, s)
	}
	for _, s := range sym.LocalVariables {
		symbols = append(symbols, s)
	}
	return symbols
}

// GetVarIndex returns the index of a variable
func (sym *FunctionSymbol) GetVarIndex(id string) int {
	for i, s := range sym.AllDeclarations() {
		if s.GetIdentifier() == id {
			return i
		}
	}
	return -1
}

// IsLocalVar checks whether the id is a local variable or not
func (sym *FunctionSymbol) IsLocalVar(id string) bool {
	for _, s := range sym.LocalVariables {
		if s.GetIdentifier() == id {
			return true
		}
	}
	return false
}

// String creates the string representation for the FunctionSymbol
func (sym *FunctionSymbol) String() string {
	return fmt.Sprintf("\n %s %s(%s): vars: %s", sym.ReturnTypes, sym.Identifier, sym.Parameters, sym.LocalVariables)
}

//----------------

// ParameterSymbol is an alias for VariableSymbol
type ParameterSymbol struct {
	VariableSymbol
}

// NewParameterSymbol creates a new ParameterSymbol
func NewParameterSymbol(scope Symbol, identifier string) *ParameterSymbol {
	return &ParameterSymbol{
		VariableSymbol: *NewVariableSymbol(scope, identifier),
	}
}

//----------------

// LocalVariableSymbol contains a variable symbol and stores information about its visibility
type LocalVariableSymbol struct {
	VariableSymbol
	VisibleIn []node.StatementNode
}

// NewLocalVariableSymbol creates a new LocalVariableSymbol
func NewLocalVariableSymbol(scope Symbol, identifier string) *LocalVariableSymbol {
	return &LocalVariableSymbol{
		VariableSymbol: *NewVariableSymbol(scope, identifier),
	}
}

//----------------

// VariableSymbol contains the variables type
type VariableSymbol struct {
	AbstractSymbol
	Type *TypeSymbol
}

// NewVariableSymbol creates a new VariableSymbol
func NewVariableSymbol(scope Symbol, identifier string) *VariableSymbol {
	return &VariableSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// String creates the string representation of the VariableSymbol
func (sym *VariableSymbol) String() string {
	return fmt.Sprintf("%s %s", sym.Type, sym.Identifier)
}

//----------------

// ConstantSymbol contains the type of the constant
type ConstantSymbol struct {
	AbstractSymbol
	Type *TypeSymbol
}

// NewConstantSymbol creates a new ConstantSymbol
func NewConstantSymbol(scope Symbol, identifier string, typeSymbol *TypeSymbol) *ConstantSymbol {
	return &ConstantSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
		Type:           typeSymbol,
	}
}

// String creates the string representation of the ConstantSymbol
func (sym *ConstantSymbol) String() string {
	return fmt.Sprintf("Constant %s", sym.Identifier)
}

//----------------

// TypeSymbol represents a type
type TypeSymbol struct {
	AbstractSymbol
}

// NewTypeSymbol creates a new TypeSymbol
func NewTypeSymbol(scope Symbol, identifier string) *TypeSymbol {
	return &TypeSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// String creates a new string representation for TypeSymbol
func (sym *TypeSymbol) String() string {
	return fmt.Sprintf("Type %s", sym.Identifier)
}
