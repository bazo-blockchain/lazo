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

func (v *DesignatorResolutionVisitor) VisitDesignatorNode(node *node.DesignatorNode) {
	var scope symbol.Symbol
	if v.currentFunctionSymbol == nil {
		scope = v.contractSymbol
	} else {
		scope = v.currentFunctionSymbol
	}
	sym := v.symbolTable.Find(scope, node.Value)
	if sym == nil {
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

func (v *DesignatorResolutionVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.currentStatement = node
	v.AbstractVisitor.VisitAssignmentStatementNode(node)
}

func (v *DesignatorResolutionVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	v.currentStatement = node
	v.AbstractVisitor.VisitIfStatementNode(node)
}

func (v *DesignatorResolutionVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.currentStatement = node
	v.AbstractVisitor.VisitReturnStatementNode(node)
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
	if field, ok := sym.(*symbol.FieldSymbol); ok {
		return field.Type
	} else if variable, ok := sym.(*symbol.VariableSymbol); ok {
		return variable.Type
	} else if localVariable, ok := sym.(*symbol.LocalVariableSymbol); ok {
		return localVariable.Type
	} else if constant, ok := sym.(*symbol.ConstantSymbol); ok { // TODO Remove: not needed
		return constant.Type
	} else {
		return nil
	}
}
