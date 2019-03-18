package visitor

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
)

type DesignatorResolutionVisitor struct {
	node.AbstractVisitor
	symbolTable           *symbol.SymbolTable
	contractSymbol        *symbol.ContractSymbol
	currentFunctionSymbol *symbol.FunctionSymbol
	currentStatement      node.StatementNode
	Errors                []error
}

func NewDesignatorResolutionVisitor(symbolTable *symbol.SymbolTable, contractSymbol *symbol.ContractSymbol) *DesignatorResolutionVisitor {
	v := &DesignatorResolutionVisitor{
		symbolTable:    symbolTable,
		contractSymbol: contractSymbol,
	}
	v.ConcreteVisitor = v
	return v
}

func (v *DesignatorResolutionVisitor) VisitContractNode(node *node.ContractNode) {
	for _, variable := range node.Variables {
		variable.Accept(v.ConcreteVisitor)
	}

	for _, function := range v.contractSymbol.Functions {
		v.currentFunctionSymbol = function
		functionNode := v.symbolTable.GetNodeBySymbol(function)
		functionNode.Accept(v)
		v.currentFunctionSymbol = nil
	}
}

func (v *DesignatorResolutionVisitor) VisitStatementBlock(stmts []node.StatementNode){
	for _, statement := range stmts {
		v.currentStatement = statement
		statement.Accept(v.ConcreteVisitor)
		v.currentStatement = nil
	}
}

func (v *DesignatorResolutionVisitor) VisitDesignatorNode(node *node.DesignatorNode) {
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
	} else if local, ok := sym.(*symbol.LocalVariableSymbol); ok {
		if !containsStatement(local.VisibleIn, v.currentStatement) {
			v.reportError(node, fmt.Sprintf("Local Variable %s is not visible", node.Value))
			return
		}
	}
	v.symbolTable.MapDesignatorToDecl(node, sym)
	v.symbolTable.MapExpressionToType(node, getType(sym))
}

func (v *DesignatorResolutionVisitor) reportError(node node.Node, msg string) {
	v.Errors = append(v.Errors, errors.New(
		fmt.Sprintf("[%s] %s", node.Pos(), msg)))
}

func containsStatement(list []node.StatementNode, element node.StatementNode) bool {
	for _, listElement := range list {
		if listElement == element {
			return true
		}
	}
	return false
}

func getType(sym symbol.Symbol) *symbol.TypeSymbol {
	switch sym.(type){
	case *symbol.FieldSymbol:
		return sym.(*symbol.FieldSymbol).Type
	case *symbol.ParameterSymbol:
		return sym.(*symbol.ParameterSymbol).Type
	case *symbol.LocalVariableSymbol:
		return sym.(*symbol.LocalVariableSymbol).Type
	default:
		panic(fmt.Sprintf("Unsupported designator target symbol %s", sym.GetIdentifier()))
	}
}

func isAllowedTarget(sym symbol.Symbol) bool {
	switch sym.(type) {
	case *symbol.ContractSymbol, *symbol.FunctionSymbol:
		return false
	default:
		return true
	}
}
