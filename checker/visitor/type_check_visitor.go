package visitor

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/lexer/token"
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

func (v *TypeCheckVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.AbstractVisitor.VisitAssignmentStatementNode(node)

	if _, ok := v.symbolTable.GetTarget(node.Left).(*symbol.FunctionSymbol); ok {
		fmt.Print("Error: Assignment to function is not allowed.")
	} else {
		leftType := v.symbolTable.FindTypeByExpressionNode(node.Left)
		rightType := v.symbolTable.FindTypeByExpressionNode(node.Right)

		if leftType != rightType {
			fmt.Printf("[%s] Error: %s of assignment is not compatible with target %s\n", node.Pos(), rightType, leftType)
		}

	}
}

func (v *TypeCheckVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	v.AbstractVisitor.VisitIfStatementNode(node)
	if !v.IsBool(v.symbolTable.FindTypeByExpressionNode(node.Condition)) {
		fmt.Printf("Error condition must return boolean.\n")
	}
}

func (v *TypeCheckVisitor) VisitBinaryExpressionNode(node *node.BinaryExpressionNode) {
	v.AbstractVisitor.VisitBinaryExpressionNode(node)
	left := node.Left
	right := node.Right
	leftType := v.symbolTable.FindTypeByExpressionNode(left)
	rightType := v.symbolTable.FindTypeByExpressionNode(right)
	switch node.Operator {
	case token.And, token.Or:
		if !v.IsBool(leftType) || !v.IsBool(rightType) {
			fmt.Print("&& and || can only be applied to expressions of type bool.\n")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
		break
	case token.Addition, token.Subtraction, token.Multiplication, token.Division, token.Modulo:
		if !v.IsInt(leftType) || !v.IsInt(rightType) {
			fmt.Print("Arithmetic operators can only be applied to expressions of type int.\n")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
		break
	case token.Less, token.LessEqual, token.Equal, token.Unequal, token.GreaterEqual, token.Greater:
		if leftType != rightType{
			fmt.Print("Both sides of a compare operation need to have compatible types.\n")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
	default:

	}
}

func (v *TypeCheckVisitor) VisitUnaryExpressionNode(node *node.UnaryExpression) {
	v.AbstractVisitor.VisitUnaryExpressionNode(node)
	operand := node.Expression
	operandType := v.symbolTable.FindTypeByExpressionNode(operand)
	switch node.Operator {
	case token.Addition, token.Subtraction:
		if !v.IsInt(operandType) {
			fmt.Print("+ and - unary operators can only be applied to expressions of type int.\n")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
		break
	case token.Not:
		if !v.IsBool(operandType) {
			fmt.Print("! unary operators can only be applied to expressions of type bool.\n")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
		break
	default:
		fmt.Print("Illegal unary operator found.\n")
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

func (v *TypeCheckVisitor) IsInt(symbol *symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.IntType
}

func (v *TypeCheckVisitor) IsBool(symbol *symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.BoolType
}