package symbol

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/parser/node"
)

// SymbolTable maps symbols to nodes, designators to declarations, expressions to types and contains the global scope
type SymbolTable struct {
	GlobalScope            *GlobalScope
	symbolToNode           map[Symbol]node.Node
	designatorDeclarations map[node.DesignatorNode]Symbol
	expressionTypes        map[node.ExpressionNode]TypeSymbol
}

// NewSymbolTable creates a new symbol table and initializes mappings
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		GlobalScope:            newGlobalScope(),
		symbolToNode:           make(map[Symbol]node.Node),
		designatorDeclarations: make(map[node.DesignatorNode]Symbol),
		expressionTypes:        make(map[node.ExpressionNode]TypeSymbol),
	}
}

// FindTypeByNode searches for a type symbol.
// If an array type or map type are not found, they will be added to the global types.
func (t *SymbolTable) FindTypeByNode(typeNode node.TypeNode) TypeSymbol {
	if typeSymbol := t.FindTypeByIdentifier(typeNode.Type()); typeSymbol != nil {
		return typeSymbol
	}

	switch typeNode.(type) {
	case *node.ArrayTypeNode:
		arrayType := typeNode.(*node.ArrayTypeNode)
		return t.AddArrayType(arrayType.ElementType)
	case *node.MapTypeNode:
		mapType := typeNode.(*node.MapTypeNode)
		return t.AddMapType(mapType)
	default:
		return nil
	}
}

// AddArrayType creates a new array type symbol and adds it to the global types
func (t *SymbolTable) AddArrayType(elementTypeNode node.TypeNode) TypeSymbol {
	elementType := t.FindTypeByNode(elementTypeNode)
	if elementType == nil {
		return nil
	}

	arrayType := NewArrayTypeSymbol(t.GlobalScope, elementType)
	t.GlobalScope.Types[arrayType.Identifier()] = arrayType
	return arrayType
}

// AddMapType creates a new map type and adds it to the global scope types
func (t *SymbolTable) AddMapType(mapTypeNode *node.MapTypeNode) TypeSymbol {
	keyType := t.FindTypeByNode(mapTypeNode.KeyType)
	valueType := t.FindTypeByNode(mapTypeNode.ValueType)

	// To create a map type, both types should be valid
	if keyType == nil || valueType == nil {
		return nil
	}

	mapType := NewMapTypeSymbol(t.GlobalScope, keyType, valueType)
	t.GlobalScope.Types[mapType.Identifier()] = mapType
	return mapType
}

// FindTypeByIdentifier searches for a type
// Returns the type or nil
func (t *SymbolTable) FindTypeByIdentifier(identifier string) TypeSymbol {
	if compilationType, ok := t.GlobalScope.Types[identifier]; ok {
		return compilationType
	}
	return nil
}

// Find recursively searches for a symbol within a specific scope
// Returns the symbol or nil
func (t *SymbolTable) Find(scope Symbol, identifier string) Symbol {
	if scope == nil {
		return nil
	}

	if identifier == This && scope == t.GlobalScope {
		return t.GlobalScope.Contract
	}

	for _, declaration := range scope.AllDeclarations() {

		if declaration.Identifier() == identifier {
			return declaration
		}
	}
	return t.Find(scope.Scope(), identifier)
}

// MapSymbolToNode maps a symbol to its node
func (t *SymbolTable) MapSymbolToNode(symbol Symbol, node node.Node) {
	t.symbolToNode[symbol] = node
}

// GetNodeBySymbol returns the node linked to the symbol
func (t *SymbolTable) GetNodeBySymbol(symbol Symbol) node.Node {
	return t.symbolToNode[symbol]
}

// MapDesignatorToDecl maps a designator to a declaration
func (t *SymbolTable) MapDesignatorToDecl(designatorNode node.DesignatorNode, symbol Symbol) {
	t.designatorDeclarations[designatorNode] = symbol
}

// GetDeclByDesignator returns the declaration for a designator
func (t *SymbolTable) GetDeclByDesignator(designatorNode node.DesignatorNode) Symbol {
	return t.designatorDeclarations[designatorNode]
}

// MapExpressionToType maps an expression to its type
func (t *SymbolTable) MapExpressionToType(expressionNode node.ExpressionNode, symbol TypeSymbol) {
	t.expressionTypes[expressionNode] = symbol
}

// GetTypeByExpression returns the type of the expression
func (t *SymbolTable) GetTypeByExpression(expressionNode node.ExpressionNode) TypeSymbol {
	return t.expressionTypes[expressionNode]
}

// String creates a string representation for the symbol table
func (t *SymbolTable) String() string {
	return fmt.Sprintf("Global Scope: %s", t.GlobalScope)
}
