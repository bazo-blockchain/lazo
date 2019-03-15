package symbol

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type SymbolTable struct {
	GlobalScope       *CompilationUnit
	symbolToNode      map[Symbol]node.Node
	designatorSymbols map[*node.DesignatorNode]Symbol
	exprTypes         map[node.ExpressionNode]*TypeSymbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		GlobalScope:       newCompilationUnit(),
		symbolToNode:      make(map[Symbol]node.Node),
		designatorSymbols: make(map[*node.DesignatorNode]Symbol),
		exprTypes:         make(map[node.ExpressionNode]*TypeSymbol),
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

func (t *SymbolTable) MapDesignatorToType(designatorNode *node.DesignatorNode, symbol Symbol) {
	t.designatorSymbols[designatorNode] = symbol
}

func (t *SymbolTable) FindTypeByDesignatorNode(designatorNode *node.DesignatorNode) Symbol {
	return t.designatorSymbols[designatorNode]
}

func (t *SymbolTable) MapExpressionToType(expressionNode node.ExpressionNode, symbol *TypeSymbol) {
	t.exprTypes[expressionNode] = symbol
}

func (t *SymbolTable) FindTypeByExpressionNode(expressionNode node.ExpressionNode) *TypeSymbol {
	return t.exprTypes[expressionNode]
}

func (t *SymbolTable) String() string{
	return fmt.Sprintf("Global Scope: %s", t.GlobalScope)
}
