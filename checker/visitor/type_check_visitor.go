package visitor

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type TypeCheckVisitor struct {
	node.AbstractVisitor
	symbolTable *symbol.SymbolTable
	contractSymbol *symbol.ContractSymbol
	currentFunction *symbol.FunctionSymbol
}

func NewTypeCheckVisitor(symbolTable *symbol.SymbolTable, contractSymbol *symbol.ContractSymbol) *TypeCheckVisitor {
	v := &TypeCheckVisitor{
		symbolTable: symbolTable,
		contractSymbol: contractSymbol,
	}
	v.ConcreteVisitor = v
	return v
}

func (v *TypeCheckVisitor) VisitContractNode(node *node.ContractNode) {
	for _, field := range v.contractSymbol.Fields {
		fieldNode := v.symbolTable.GetNodeBySymbol(field)
		fieldNode.Accept(v)
	}

	for _, function := range v.contractSymbol.Functions {
		v.currentFunction = function
		functionNode :=  v.symbolTable.GetNodeBySymbol(function)
		functionNode.Accept(v)
		v.currentFunction = nil
	}
}

func (v *TypeCheckVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(node)
	returnNodes := node.Expressions
	returnSymbols := v.currentFunction.ReturnTypes

	if len(returnSymbols) > 0 {
		if len(returnSymbols)!= len(returnNodes) {
			fmt.Printf("Error: Expected %d return values, given %d\n", len(returnSymbols), len(returnNodes))
		} else {
			// Check types
			for i, rtype := range returnSymbols {
				nodeType := v.symbolTable.FindTypeByExpressionNode(returnNodes[i])
				if nodeType != rtype {
					fmt.Printf("Error: Return Types mismatch expected: %s given: %s\n", rtype.Identifier, nodeType.Identifier)
				}
			}
		}
	} else {
		if len(returnNodes) > 0 {
			fmt.Printf("Error: void method should not return expression\n")
		}
	}
}

func (v *TypeCheckVisitor) VisitTypeNode(node *node.TypeNode) {
	// To be done as soon as own types are introduced
}

func (v *TypeCheckVisitor) VisitVariableNode(node *node.VariableNode) {
	v.AbstractVisitor.VisitVariableNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)
	expType := v.symbolTable.FindTypeByExpressionNode(node.Expression)
	if expType != nil && targetType != expType {
		fmt.Printf("[%s]Error Type mismatch\n", node.Pos())
	}
}

func (v *TypeCheckVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
}

func (v *TypeCheckVisitor) VisitBoolLiteralNode(node *node.BoolLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
}

func (v *TypeCheckVisitor) VisitStringLiteralNode(node *node.StringLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.StringType)
}

func (v *TypeCheckVisitor) VisitCharacterLiteralNode(node *node.CharacterLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.CharType)
}