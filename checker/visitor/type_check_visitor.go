package visitor

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
)

type TypeCheckVisitor struct {
	node.AbstractVisitor
	symbolTable     *symbol.SymbolTable
	contractSymbol  *symbol.ContractSymbol
	currentFunction *symbol.FunctionSymbol
	Errors          []error
}

func NewTypeCheckVisitor(symbolTable *symbol.SymbolTable, contractSymbol *symbol.ContractSymbol) *TypeCheckVisitor {
	v := &TypeCheckVisitor{
		symbolTable:    symbolTable,
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
		functionNode := v.symbolTable.GetNodeBySymbol(function)
		functionNode.Accept(v)
		v.currentFunction = nil
	}
}

func (v *TypeCheckVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(node)
	returnNodes := node.Expressions
	returnSymbols := v.currentFunction.ReturnTypes

	if len(returnSymbols) > 0 {
		if len(returnSymbols) != len(returnNodes) {
			v.reportError(node,
				fmt.Sprintf("Expected %d return values, given %d", len(returnSymbols), len(returnNodes)))
		} else {
			// Check types
			for i, rtype := range returnSymbols {
				nodeType := v.symbolTable.GetTypeByExpression(returnNodes[i])
				if nodeType != rtype {
					v.reportError(node, fmt.Sprintf("Return Types mismatch expected: %s given: %s",
						rtype.Identifier, nodeType.Identifier))
				}
			}
		}
	} else {
		if len(returnNodes) > 0 {
			v.reportError(node, "void method should not return expression")
		}
	}
}

func (v *TypeCheckVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.AbstractVisitor.VisitAssignmentStatementNode(node)

	if _, ok := v.symbolTable.GetDeclByDesignator(node.Left).(*symbol.FunctionSymbol); ok {
		v.reportError(node, "Assignment to function is not allowed")
	} else {
		leftType := v.symbolTable.GetTypeByExpression(node.Left)
		rightType := v.symbolTable.GetTypeByExpression(node.Right)

		if leftType.GetIdentifier() != rightType.GetIdentifier() {
			v.reportError(node,
				fmt.Sprintf("%s of assignment is not compatible with target %s", rightType, leftType))
		}
	}
}

func (v *TypeCheckVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	v.AbstractVisitor.VisitIfStatementNode(node)
	if !v.IsBool(v.symbolTable.GetTypeByExpression(node.Condition)) {
		v.reportError(node, "condition must return boolean")
	}
}

func (v *TypeCheckVisitor) VisitBinaryExpressionNode(node *node.BinaryExpressionNode) {
	v.AbstractVisitor.VisitBinaryExpressionNode(node)
	left := node.Left
	right := node.Right
	leftType := v.symbolTable.GetTypeByExpression(left)
	rightType := v.symbolTable.GetTypeByExpression(right)
	switch node.Operator {
	case token.And, token.Or:
		if !v.IsBool(leftType) || !v.IsBool(rightType) {
			v.reportError(node, "&& and || can only be applied to expressions of type bool")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
		break
	case token.Addition, token.Subtraction, token.Multiplication, token.Division, token.Modulo:
		if !v.IsInt(leftType) || !v.IsInt(rightType) {
			v.reportError(node, "Arithmetic operators can only be applied to expressions of type int")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
		break
	case token.Less, token.LessEqual, token.Equal, token.Unequal, token.GreaterEqual, token.Greater:
		if leftType != rightType {
			v.reportError(node, "Both sides of a compare operation need to have compatible types")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
	default:
		v.reportError(node, "Illegal binary operator found")
	}
}

func (v *TypeCheckVisitor) VisitUnaryExpressionNode(node *node.UnaryExpression) {
	v.AbstractVisitor.VisitUnaryExpressionNode(node)
	operand := node.Expression
	operandType := v.symbolTable.GetTypeByExpression(operand)
	switch node.Operator {
	case token.Addition, token.Subtraction:
		if !v.IsInt(operandType) {
			v.reportError(node, "+ and - unary operators can only be applied to expressions of type int")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
		break
	case token.Not:
		if !v.IsBool(operandType) {
			v.reportError(node, "! unary operators can only be applied to expressions of type bool")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
		break
	default:
		v.reportError(node, "Illegal unary operator found")
	}
}

func (v *TypeCheckVisitor) VisitTypeNode(node *node.TypeNode) {
	// To be done as soon as own types are introduced
}

func (v *TypeCheckVisitor) VisitVariableNode(node *node.VariableNode) {
	v.AbstractVisitor.VisitVariableNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)
	expType := v.symbolTable.GetTypeByExpression(node.Expression)
	if expType != nil && targetType != expType {
		v.reportError(node, "Type mismatch")
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

func (v *TypeCheckVisitor) reportError(node node.Node, msg string) {
	v.Errors = append(v.Errors, errors.New(
		fmt.Sprintf("[%s] %s", node.Pos(), msg)))
}
