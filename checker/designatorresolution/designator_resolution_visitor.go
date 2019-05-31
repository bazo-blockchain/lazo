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

// VisitBasicDesignatorNode visits the designator node, maps the designator to its declaration and
// maps the expression to the type
func (v *designatorResolutionVisitor) VisitBasicDesignatorNode(node *node.BasicDesignatorNode) {
	var scope symbol.Symbol
	if v.currentFunctionSymbol == nil {
		scope = v.contractSymbol
	} else {
		scope = v.currentFunctionSymbol
	}
	sym := v.symbolTable.Find(scope, node.Value)
	if sym == nil || !isAllowedTarget(sym) && node.Value != symbol.This {
		v.reportError(node, fmt.Sprintf("Designator %s is undefined", node.Value))
		return
	}

	if local, ok := sym.(*symbol.LocalVariableSymbol); ok && !containsStatement(local.VisibleIn, v.currentStatement) {
		v.reportError(node, fmt.Sprintf("Local Variable %s is not visible", node.Value))
		return
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
	if arrayType, ok := typeSymbol.(*symbol.ArrayTypeSymbol); ok {
		v.symbolTable.MapExpressionToType(node, arrayType.ElementType)
		v.symbolTable.MapDesignatorToDecl(node, arrayType)
	} else if mapType, ok := typeSymbol.(*symbol.MapTypeSymbol); ok {
		v.symbolTable.MapExpressionToType(node, mapType.ValueType)
		v.symbolTable.MapDesignatorToDecl(node, mapType)
	} else {
		v.reportError(node, fmt.Sprintf("Designator %v does not refer to an array/map type", node))
	}
}

func (v *designatorResolutionVisitor) VisitMemberAccessNode(node *node.MemberAccessNode) {
	v.AbstractVisitor.VisitMemberAccessNode(node)
	designatorType := v.symbolTable.GetTypeByExpression(node.Designator)

	if node.Identifier == symbol.This {
		v.reportError(node, "Invalid member designator 'this'")
		return
	}

	switch designatorType.(type) {
	case *symbol.ArrayTypeSymbol:
		v.visitArrayMemberAccess(node)
	case *symbol.StructTypeSymbol:
		v.visitStructMemberAccess(node, designatorType.(*symbol.StructTypeSymbol))
	case *symbol.MapTypeSymbol:
		v.visitMapMemberAccess(node)
	case *symbol.ContractSymbol:
		v.visitContractMemberAccess(node, designatorType.(*symbol.ContractSymbol))
	default:
		v.reportError(node, fmt.Sprintf("Designator %v does not refer to a composite type", node))
	}
}

func (v *designatorResolutionVisitor) visitArrayMemberAccess(node *node.MemberAccessNode) {
	arrayLength := v.symbolTable.GlobalScope.ArrayLengthField
	if node.Identifier != arrayLength.Identifier() {
		v.reportError(node, fmt.Sprintf("Invalid member access %v on array %v", node.Identifier, node))
		return
	}

	targetType, err := getType(arrayLength)
	if err != nil {
		v.reportError(node, err.Error())
	}

	v.symbolTable.MapDesignatorToDecl(node, arrayLength)
	v.symbolTable.MapExpressionToType(node, targetType)
}

func (v *designatorResolutionVisitor) visitStructMemberAccess(node *node.MemberAccessNode, structType *symbol.StructTypeSymbol) {
	fieldSymbol := structType.GetField(node.Identifier)

	if fieldSymbol == (*symbol.FieldSymbol)(nil) {
		v.reportError(node, fmt.Sprintf("Member %s does not exist on struct %s",
			node.Identifier, structType.Identifier()))
		return
	}

	targetType, err := getType(fieldSymbol)
	if err != nil {
		v.reportError(node, err.Error())
	}

	v.symbolTable.MapDesignatorToDecl(node, fieldSymbol)
	v.symbolTable.MapExpressionToType(node, targetType)
}

func (v *designatorResolutionVisitor) visitMapMemberAccess(node *node.MemberAccessNode) {
	if node.Identifier != symbol.Contains {
		v.reportError(node, fmt.Sprintf("Invalid member access %v on map %v", node.Identifier, node))
		return
	}
	target := v.symbolTable.GlobalScope.MapMemberFunctions[symbol.Contains]
	v.symbolTable.MapDesignatorToDecl(node, target)
}

func (v *designatorResolutionVisitor) visitContractMemberAccess(node *node.MemberAccessNode, contractType *symbol.ContractSymbol) {
	targetIndex := contractType.GetFieldIndex(node.Identifier)
	if targetIndex < 0 {
		v.reportError(node, fmt.Sprintf("Member %s does not exist on contract %v", node.Identifier, contractType.Identifier()))
		return
	}

	target := contractType.Fields[targetIndex]
	targetType, err := getType(target)

	if err != nil {
		v.reportError(node, err.Error())
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
	case *symbol.ContractSymbol:
		return sym.(*symbol.ContractSymbol), nil
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
