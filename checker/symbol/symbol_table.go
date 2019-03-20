package symbol

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type SymbolTable struct {
	GlobalScope            *GlobalScope
	symbolToNode           map[Symbol]node.Node
	designatorDeclarations map[*node.DesignatorNode]Symbol
	expressionTypes        map[node.ExpressionNode]*TypeSymbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		GlobalScope:            newGlobalScope(),
		symbolToNode:           make(map[Symbol]node.Node),
		designatorDeclarations: make(map[*node.DesignatorNode]Symbol),
		expressionTypes:        make(map[node.ExpressionNode]*TypeSymbol),
	}
}

func (t *SymbolTable) FindTypeByIdentifier(identifier string) *TypeSymbol {
	for _, compilationType := range t.GlobalScope.Types {
		if compilationType.GetIdentifier() == identifier {
			return compilationType
		}
	}

	return nil
}

func (t *SymbolTable) FindTypeByNode(node *node.TypeNode) *TypeSymbol {
	return t.FindTypeByIdentifier(node.Identifier)
}

func (t *SymbolTable) Find(scope Symbol, identifier string) Symbol {
	if scope == nil {
		return nil
	}
	for _, declaration := range scope.AllDeclarations() {
		if declaration.GetIdentifier() == identifier {
			return declaration
		}
	}
	return t.Find(scope.GetScope(), identifier)
}

func (t *SymbolTable) MapSymbolToNode(symbol Symbol, node node.Node) {
	t.symbolToNode[symbol] = node
}

func (t *SymbolTable) GetNodeBySymbol(symbol Symbol) node.Node {
	return t.symbolToNode[symbol]
}

func (t *SymbolTable) MapDesignatorToDecl(designatorNode *node.DesignatorNode, symbol Symbol) {
	t.designatorDeclarations[designatorNode] = symbol
}

func (t *SymbolTable) GetDeclByDesignator(designatorNode *node.DesignatorNode) Symbol {
	return t.designatorDeclarations[designatorNode]
}

func (t *SymbolTable) MapExpressionToType(expressionNode node.ExpressionNode, symbol *TypeSymbol) {
	t.expressionTypes[expressionNode] = symbol
}

func (t *SymbolTable) GetTypeByExpression(expressionNode node.ExpressionNode) *TypeSymbol {
	return t.expressionTypes[expressionNode]
}

func (t *SymbolTable) String() string {
	return fmt.Sprintf("Global Scope: %s", t.GlobalScope)
}
