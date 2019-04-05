package visitor

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
)

// TypeCheckVisitor contains the symbol table, contract symbol, current function and errors. It traverses the abstract
// syntax tree and checks if types match
type TypeCheckVisitor struct {
	node.AbstractVisitor
	symbolTable     *symbol.SymbolTable
	contractSymbol  *symbol.ContractSymbol
	currentFunction *symbol.FunctionSymbol
	Errors          []error
}

// NewTypeCheckVisitor creates a new TypeCheckVisitor
func NewTypeCheckVisitor(symbolTable *symbol.SymbolTable, contractSymbol *symbol.ContractSymbol) *TypeCheckVisitor {
	v := &TypeCheckVisitor{
		symbolTable:    symbolTable,
		contractSymbol: contractSymbol,
	}
	v.ConcreteVisitor = v
	return v
}

// VisitContractNode visits the fields and functions of the contract
func (v *TypeCheckVisitor) VisitContractNode(node *node.ContractNode) {
	for _, variable := range node.Variables {
		variable.Accept(v.ConcreteVisitor)
	}

	for _, function := range v.contractSymbol.Functions {
		v.currentFunction = function
		functionNode := v.symbolTable.GetNodeBySymbol(function)
		functionNode.Accept(v)
		v.currentFunction = nil
	}
}

// Statements
// ----------

// VisitVariableNode checks whether the variable type and value are of the same type
func (v *TypeCheckVisitor) VisitVariableNode(node *node.VariableNode) {
	v.AbstractVisitor.VisitVariableNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)
	expType := v.symbolTable.GetTypeByExpression(node.Expression)

	if expType != nil && targetType != expType {
		v.reportError(node, fmt.Sprintf("Type mismatch: expected %s, given %s", targetType, expType))
	}
}

// VisitReturnStatementNode checks whether the return types and the values are of the same type
func (v *TypeCheckVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(node)
	returnNodes := node.Expressions
	returnSymbols := v.currentFunction.ReturnTypes

	if len(returnSymbols) > 0 {
		if len(returnSymbols) != len(returnNodes) {
			v.reportError(node,
				fmt.Sprintf("Expected %d return values, given %d", len(returnSymbols), len(returnNodes)))
		} else {
			for i, rtype := range returnSymbols {
				nodeType := v.symbolTable.GetTypeByExpression(returnNodes[i])
				if nodeType != rtype {
					v.reportError(node, fmt.Sprintf("Return Type mismatch: expected %s, given %s",
						rtype.Identifier, nodeType.Identifier))
				}
			}
		}
	} else if len(returnNodes) > 0 {
		v.reportError(node, "void method should not return expression")
	}
}

// VisitAssignmentStatementNode checks whether the left and right part of the assignment are of the same type
func (v *TypeCheckVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.AbstractVisitor.VisitAssignmentStatementNode(node)

	leftType := v.symbolTable.GetTypeByExpression(node.Left)
	rightType := v.symbolTable.GetTypeByExpression(node.Right)

	if leftType != rightType {
		v.reportError(node,
			fmt.Sprintf("assignment of %s is not compatible with target %s", rightType, leftType))
	}
}

// VisitIfStatementNode checks whether the condition is a boolean expression
func (v *TypeCheckVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	v.AbstractVisitor.VisitIfStatementNode(node)
	if !v.isBool(v.symbolTable.GetTypeByExpression(node.Condition)) {
		v.reportError(node, "condition must return boolean")
	}
}

// Expressions
// -----------

// VisitBinaryExpressionNode checks if the types for different binary expressions match
// Expressions are &&, ||, +, -, *, /, %, **, ==, !=, >, >=, <= and <
func (v *TypeCheckVisitor) VisitBinaryExpressionNode(node *node.BinaryExpressionNode) {
	v.AbstractVisitor.VisitBinaryExpressionNode(node)
	left := node.Left
	right := node.Right
	leftType := v.symbolTable.GetTypeByExpression(left)
	rightType := v.symbolTable.GetTypeByExpression(right)

	switch node.Operator {
	case token.And, token.Or:
		if !v.isBool(leftType) || !v.isBool(rightType) {
			v.reportError(node, "&& and || can only be applied to expressions of type bool")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
	case token.Addition, token.Subtraction, token.Multiplication, token.Division, token.Modulo, token.Exponent:
		if !v.isInt(leftType) || !v.isInt(rightType) {
			v.reportError(node, "Arithmetic operators can only be applied to expressions of type int")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
	case token.Equal, token.Unequal:
		if leftType != rightType {
			v.reportError(node, fmt.Sprintf("Equality comparison should have the same type, given %s and %s",
				leftType, rightType))
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
	case token.Less, token.LessEqual, token.GreaterEqual, token.Greater:
		if leftType != rightType {
			v.reportError(node,
				fmt.Sprintf("Both sides of a compare operation need to have the same type, given %s and %s",
					leftType, rightType))
		} else if !(v.isInt(leftType) || v.isChar(leftType)) {
			v.reportError(node, fmt.Sprintf("Relational comparison is not supported for %s", leftType))
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
	default:
		panic(fmt.Sprintf("Illegal binary operator %s", token.SymbolLexeme[node.Operator]))
	}
}

// VisitUnaryExpressionNode checks that types of unary expressions are valid
// Expressions are +, -, !
func (v *TypeCheckVisitor) VisitUnaryExpressionNode(node *node.UnaryExpressionNode) {
	v.AbstractVisitor.VisitUnaryExpressionNode(node)
	operand := node.Expression
	operandType := v.symbolTable.GetTypeByExpression(operand)

	switch node.Operator {
	case token.Addition, token.Subtraction:
		if !v.isInt(operandType) {
			v.reportError(node, "+ and - unary operators can only be applied to expressions of type int")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
		break
	case token.Not:
		if !v.isBool(operandType) {
			v.reportError(node, "! unary operators can only be applied to expressions of type bool")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
		break
	default:
		panic(fmt.Sprintf("Illegal unary operator %s", token.SymbolLexeme[node.Operator]))
	}
}

// VisitTypeNode currently does nothing
func (v *TypeCheckVisitor) VisitTypeNode(node *node.TypeNode) {
	// To be done as soon as own types are introduced
}

// VisitIntegerLiteralNode maps the integer literal node to its type
func (v *TypeCheckVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
}

// VisitBoolLiteralNode maps the bool literal node to its type
func (v *TypeCheckVisitor) VisitBoolLiteralNode(node *node.BoolLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
}

// VisitStringLiteralNode maps the string literal to its type
func (v *TypeCheckVisitor) VisitStringLiteralNode(node *node.StringLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.StringType)
}

// VisitCharacterLiteralNode maps the character literal to its type
func (v *TypeCheckVisitor) VisitCharacterLiteralNode(node *node.CharacterLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.CharType)
}

func (v *TypeCheckVisitor) isInt(symbol *symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.IntType
}

func (v *TypeCheckVisitor) isBool(symbol *symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.BoolType
}

func (v *TypeCheckVisitor) isChar(symbol *symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.CharType
}

func (v *TypeCheckVisitor) reportError(node node.Node, msg string) {
	v.Errors = append(v.Errors, errors.New(
		fmt.Sprintf("[%s] %s", node.Pos(), msg)))
}
