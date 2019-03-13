package symbol

import "github.com/bazo-blockchain/lazo/parser/node"

type SymbolTable struct {
	Compilation *CompilationUnit
	symbolToNode map[Symbol]node.Node
	targetFixup map[*node.DesignatorNode]Symbol
	typeFixup map[node.ExpressionNode]*TypeSymbol
}

func (t *SymbolTable) FindTypeByIdentifier(identifier string) *TypeSymbol {
	for _, compilationType := range t.Compilation.Types {
		if compilationType.String() == identifier {
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
		if declaration.String() == identifier {
			return declaration
		}
	}

	return t.Find(scope.GetScope(), identifier)
}

func (t *SymbolTable) LinkDeclaration(node node.Node, symbol Symbol) {
	t.symbolToNode[symbol] = node
}

func (t *SymbolTable) GetDeclaration(symbol Symbol) node.Node {
	return t.symbolToNode[symbol]
}

func (t *SymbolTable) FixTarget(designatorNode *node.DesignatorNode, symbol Symbol) {
	t.targetFixup[designatorNode] = symbol
}

func (t *SymbolTable) GetTarget(designatorNode *node.DesignatorNode) Symbol{
	return t.targetFixup[designatorNode]
}

func (t *SymbolTable) FixType(expressionNode node.ExpressionNode, symbol *TypeSymbol) {
	t.typeFixup[expressionNode] = symbol
}

func (t *SymbolTable) FindTypeByExpressionNode(expressionNode node.ExpressionNode) *TypeSymbol {
	return t.typeFixup[expressionNode]
}