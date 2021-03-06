// Package symbol contains all the supported symbol types and their functions.
package symbol

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/parser/node"
)

// This is a constant for the 'this' keyword
const This = "this"

// Contains is an identifier for built-in map member function.
const Contains = "contains"

// Symbol declares functions which symbols have to implement
type Symbol interface {
	Scope() Symbol
	Identifier() string
	AllDeclarations() []Symbol
	String() string
}

// AbstractSymbol is part of all symbols and contains the scope and identifier
type AbstractSymbol struct {
	Parent Symbol
	ID     string
}

// NewAbstractSymbol is a helper function that creates a new abstract symbol
func NewAbstractSymbol(scope Symbol, identifier string) AbstractSymbol {
	return AbstractSymbol{
		Parent: scope,
		ID:     identifier,
	}
}

// Scope returns the scope of the symbol (will be deleted)
func (sym *AbstractSymbol) Scope() Symbol {
	return sym.Parent
}

// Identifier returns the identifier of the identifier (will be deleted)
func (sym *AbstractSymbol) Identifier() string {
	return sym.ID
}

// AllDeclarations returns an empty symbol slice
func (sym *AbstractSymbol) AllDeclarations() []Symbol {
	return []Symbol{}
}

// String creates the string representation for the abstract symbol
func (sym *AbstractSymbol) String() string {
	return fmt.Sprintf("Abstract Symbol: %s", sym.ID)
}

// Concrete Symbols
//-----------------

// ContractSymbol contains fields and functions
type ContractSymbol struct {
	AbstractSymbol
	Fields      []*FieldSymbol
	Constructor *FunctionSymbol
	Functions   []*FunctionSymbol
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
		if s.Identifier() == id {
			return i
		}
	}
	return -1
}

// String creates the string representation
func (sym *ContractSymbol) String() string {
	return fmt.Sprintf("Contract: %s, \nFields: %s, \nConstructor %s \nFunctions %s",
		sym.ID, sym.Fields, sym.Constructor, sym.Functions)
}

//----------------

// FieldSymbol contains the type of the field
type FieldSymbol struct {
	AbstractSymbol
	Type TypeSymbol
}

// NewFieldSymbol creates a new FieldSymbol
func NewFieldSymbol(scope Symbol, identifier string) *FieldSymbol {
	return &FieldSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// String creates the string representation
func (sym *FieldSymbol) String() string {
	return fmt.Sprintf("%s %s", sym.Type, sym.ID)
}

//----------------

// FunctionSymbol contains the return types, parameters and local variables of a function
type FunctionSymbol struct {
	AbstractSymbol
	ReturnTypes    []TypeSymbol
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
		if s.Identifier() == id {
			return i
		}
	}
	return -1
}

// IsLocalVar checks whether the id is a local variable or not
func (sym *FunctionSymbol) IsLocalVar(id string) bool {
	for _, s := range sym.LocalVariables {
		if s.Identifier() == id {
			return true
		}
	}
	return false
}

// String creates the string representation
func (sym *FunctionSymbol) String() string {
	return fmt.Sprintf("\n %s %s(%s): vars: %s", sym.ReturnTypes, sym.ID, sym.Parameters, sym.LocalVariables)
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
	Type TypeSymbol
}

// NewVariableSymbol creates a new VariableSymbol
func NewVariableSymbol(scope Symbol, identifier string) *VariableSymbol {
	return &VariableSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// String creates the string representation
func (sym *VariableSymbol) String() string {
	return fmt.Sprintf("%s %s", sym.Type, sym.ID)
}

//----------------

// ConstantSymbol contains the type of the constant
type ConstantSymbol struct {
	AbstractSymbol
	Type TypeSymbol
}

// NewConstantSymbol creates a new ConstantSymbol
func NewConstantSymbol(scope Symbol, identifier string, typeSymbol TypeSymbol) *ConstantSymbol {
	return &ConstantSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
		Type:           typeSymbol,
	}
}

// String creates the string representation
func (sym *ConstantSymbol) String() string {
	return fmt.Sprintf("Constant %s", sym.ID)
}

//----------------

// TypeSymbol declares functions which all type symbols should implement
type TypeSymbol interface {
	Symbol
}

//----------------

// BasicTypeSymbol represents a type
type BasicTypeSymbol struct {
	AbstractSymbol
}

// NewBasicTypeSymbol creates a new BasicTypeSymbol
func NewBasicTypeSymbol(scope Symbol, identifier string) *BasicTypeSymbol {
	return &BasicTypeSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// String creates a new string representation
func (sym *BasicTypeSymbol) String() string {
	return fmt.Sprintf("Type %s", sym.ID)
}

//----------------

// ArrayTypeSymbol represents a array type
type ArrayTypeSymbol struct {
	AbstractSymbol
	ElementType TypeSymbol
}

// NewArrayTypeSymbol creates a new ArrayTypeSymbol
func NewArrayTypeSymbol(scope Symbol, elementType TypeSymbol) *ArrayTypeSymbol {
	return &ArrayTypeSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, elementType.Identifier()+"[]"),
		ElementType:    elementType,
	}
}

// String creates a new string representation
func (sym *ArrayTypeSymbol) String() string {
	return fmt.Sprintf("Array of %s", sym.ElementType)
}

//----------------

// StructTypeSymbol represents a struct type
type StructTypeSymbol struct {
	AbstractSymbol
	Fields []*FieldSymbol
}

// NewStructTypeSymbol creates a new StructTypeSymbol
func NewStructTypeSymbol(scope Symbol, identifier string) *StructTypeSymbol {
	return &StructTypeSymbol{
		AbstractSymbol: NewAbstractSymbol(scope, identifier),
	}
}

// AllDeclarations returns all field declarations
func (sym *StructTypeSymbol) AllDeclarations() []Symbol {
	symbols := make([]Symbol, len(sym.Fields))
	for i, s := range sym.Fields {
		symbols[i] = s
	}
	return symbols
}

// String creates the string representation
func (sym *StructTypeSymbol) String() string {
	return fmt.Sprintf("Struct: %s, \nFields: %s", sym.Identifier(), sym.Fields)
}

// GetField returns the field symbol by identifier
func (sym *StructTypeSymbol) GetField(identifier string) *FieldSymbol {
	for _, f := range sym.Fields {
		if f.Identifier() == identifier {
			return f
		}
	}
	return nil
}

// GetFieldIndex returns the index of the field
func (sym *StructTypeSymbol) GetFieldIndex(identifier string) int {
	for i, s := range sym.Fields {
		if s.Identifier() == identifier {
			return i
		}
	}
	return -1
}

//----------------

// MapTypeSymbol represents a map type consists of key and value type
type MapTypeSymbol struct {
	AbstractSymbol
	KeyType   TypeSymbol
	ValueType TypeSymbol
}

// NewMapTypeSymbol creates a new ArrayTypeSymbol
func NewMapTypeSymbol(scope Symbol, keyType TypeSymbol, valueType TypeSymbol) *MapTypeSymbol {
	return &MapTypeSymbol{
		AbstractSymbol: NewAbstractSymbol(scope,
			fmt.Sprintf("Map<%s,%s>", keyType.Identifier(), valueType.Identifier())),
		KeyType:   keyType,
		ValueType: valueType,
	}
}

// String creates a new string representation
func (sym *MapTypeSymbol) String() string {
	return sym.Identifier()
}
