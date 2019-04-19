package designatorresolution

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type designatorResolutionVisitor struct {
	node.AbstractVisitor
	symbolTable           *symbol.SymbolTable
	contractSymbol        *symbol.ContractSymbol
	currentFunctionSymbol *symbol.FunctionSymbol
	currentStatement      node.StatementNode
	Errors                []error
}

func newDesignatorResolutionVisitor(symbolTable *symbol.SymbolTable, contractSymbol *symbol.ContractSymbol) *designatorResolutionVisitor {
	v := &designatorResolutionVisitor{
		symbolTable:    symbolTable,
		contractSymbol: contractSymbol,
	}
	v.ConcreteVisitor = v
	return v
}

// VisitContractNode visits all fields and functions of the contract. Stores the current function in the visitor.
func (v *designatorResolutionVisitor) VisitContractNode(node *node.ContractNode) {
	for _, variable := range node.Fields {
		variable.Accept(v.ConcreteVisitor)
	}

	if node.Constructor != nil {
		v.currentFunctionSymbol = v.contractSymbol.Constructor
		node.Constructor.Accept(v)
		v.currentFunctionSymbol = nil
	}

	for _, function := range v.contractSymbol.Functions {
		v.currentFunctionSymbol = function
		functionNode := v.symbolTable.GetNodeBySymbol(function)
		functionNode.Accept(v)
		v.currentFunctionSymbol = nil
	}
}

// VisitStatementBlock visits all the statements of the statement block
func (v *designatorResolutionVisitor) VisitStatementBlock(stmts []node.StatementNode) {
	for _, statement := range stmts {
		v.currentStatement = statement
		statement.Accept(v.ConcreteVisitor)
		v.currentStatement = nil
	}
}

// VisitDesignatorNode visits the designator node, maps the designator to its declaration and
// maps the expression to the type
func (v *designatorResolutionVisitor) VisitDesignatorNode(node *node.DesignatorNode) {
	var scope symbol.Symbol
	if v.currentFunctionSymbol == nil {
		scope = v.contractSymbol
	} else {
		scope = v.currentFunctionSymbol
	}
	sym := v.symbolTable.Find(scope, node.Value)
	if sym == nil || !isAllowedTarget(sym) {
		v.reportError(node, fmt.Sprintf("Designator %s is undefined", node.Value))
		return
	}

	if local, ok := sym.(*symbol.LocalVariableSymbol); ok {
		if !containsStatement(local.VisibleIn, v.currentStatement) {
			v.reportError(node, fmt.Sprintf("Local Variable %s is not visible", node.Value))
			return
		}
	}
	v.symbolTable.MapDesignatorToDecl(node, sym)
	symType, err := getType(sym)
	if err != nil {
		v.reportError(node, err.Error())
	} else {
		v.symbolTable.MapExpressionToType(node, symType)
	}
}

func (v *designatorResolutionVisitor) VisitElementAccessNode(node *node.ElementAccessNode) {
	v.AbstractVisitor.VisitElementAccessNode(node)
	typeSymbol := v.symbolTable.GetTypeByExpression(node.Designator)
	if array, ok := typeSymbol.(*symbol.ArrayTypeSymbol); ok {
		v.symbolTable.MapExpressionToType(node, &array.ElementType)
		v.symbolTable.MapDesignatorToDecl(node, &array.ElementType)
	} else {
		v.reportError(node, fmt.Sprintf("Designator %v does not refer to an array type", node))
		v.symbolTable.MapExpressionToType(node, nil)
	}

}

func (v *designatorResolutionVisitor) VisitMemberAccessNode(node *node.MemberAccessNode) {
	v.AbstractVisitor.VisitMemberAccessNode(node)
	typeSymbol := v.symbolTable.GetTypeByExpression(node)
	var target symbol.Symbol
	var targetType symbol.TypeSymbol
	switch typeSymbol.(type) {
	case *symbol.ArrayTypeSymbol:
		arrayLength := v.symbolTable.GlobalScope.ArrayLength
		if node.Identifier == arrayLength.Identifier() {
			target = arrayLength
			var err error
			targetType, err = getType(arrayLength)
			if err != nil {
				v.reportError(node, err.Error())
			}
		} else {
			v.reportError(node, fmt.Sprintf("Invalid member access %v on array %v", node.Identifier, node))
		}

	default:
		v.reportError(node, fmt.Sprintf("Designator %v does not refer to a class type", node))
	}
	v.symbolTable.MapDesignatorToDecl(node, target)
	v.symbolTable.MapExpressionToType(node, targetType)
}

func (v *designatorResolutionVisitor) reportError(node node.Node, msg string) {
	v.Errors = append(v.Errors, fmt.Errorf("[%s] %s", node.Pos(), msg))
}

func containsStatement(list []node.StatementNode, element node.StatementNode) bool {
	for _, listElement := range list {
		if listElement == element {
			return true
		}
	}
	return false
}

func getType(sym symbol.Symbol) (symbol.TypeSymbol, error) {
	switch sym.(type) {
	case *symbol.FieldSymbol:
		return sym.(*symbol.FieldSymbol).Type, nil
	case *symbol.ParameterSymbol:
		return sym.(*symbol.ParameterSymbol).Type, nil
	case *symbol.LocalVariableSymbol:
		return sym.(*symbol.LocalVariableSymbol).Type, nil
	case *symbol.FunctionSymbol:
		// FuncCall expression type will be resolved in type checker
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported designator target symbol %s", sym.Identifier())
	}
}

func isAllowedTarget(sym symbol.Symbol) bool {
	switch sym.(type) {
	case *symbol.ContractSymbol:
		return false
	default:
		return true
	}
}
