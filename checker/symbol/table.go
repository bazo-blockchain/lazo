package symbol

import "github.com/bazo-blockchain/lazo/parser/node"

type Table struct {
	Compilation *CompilationUnit
	symbolToNode map[Symbol]node.Node
}

func (t *Table) FindTypeByIdentifier(identifier string) *TypeSymbol {
	for _, compilationType := range t.Compilation.Types {
		if compilationType.String() == identifier {
			return compilationType
		}
	}

	return nil
}

func (t *Table) FindTypeByNode(node *node.TypeNode) *TypeSymbol {
	return t.FindTypeByIdentifier(node.Identifier)
}

func (t *Table) Find(scope Symbol, identifier string) Symbol {
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

func (t *Table) LinkDeclaration(node node.Node, symbol Symbol) {
	t.symbolToNode[symbol] = node
}

func (t *Table) GetDeclaration(symbol Symbol) node.Node {
	return t.symbolToNode[symbol]
}