package typecheck

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type typeCheckVisitor struct {
	node.AbstractVisitor
	symbolTable     *symbol.SymbolTable
	contractSymbol  *symbol.ContractSymbol
	currentFunction *symbol.FunctionSymbol
	Errors          []error
}

func newTypeCheckVisitor(symbolTable *symbol.SymbolTable, contractSymbol *symbol.ContractSymbol) *typeCheckVisitor {
	v := &typeCheckVisitor{
		symbolTable:    symbolTable,
		contractSymbol: contractSymbol,
	}
	v.ConcreteVisitor = v
	return v
}

// VisitContractNode visits the fields and functions of the contract
func (v *typeCheckVisitor) VisitContractNode(node *node.ContractNode) {
	for _, variable := range node.Fields {
		variable.Accept(v.ConcreteVisitor)
	}

	for _, function := range v.contractSymbol.Functions {
		v.currentFunction = function
		functionNode := v.symbolTable.GetNodeBySymbol(function)
		functionNode.Accept(v)
		v.currentFunction = nil
	}
}

// VisitFieldNode checks whether the variable type and value are of the same type
func (v *typeCheckVisitor) VisitFieldNode(node *node.FieldNode) {
	v.AbstractVisitor.VisitFieldNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)

	if node.Expression != nil {
		v.checkExpressionTypes(node.Expression, targetType)
	}
}

// Statements
// ----------

// VisitVariableNode checks whether the variable type and value are of the same type
func (v *typeCheckVisitor) VisitVariableNode(node *node.VariableNode) {
	v.AbstractVisitor.VisitVariableNode(node)
	targetType := v.symbolTable.FindTypeByNode(node.Type)

	if node.Expression != nil {
		v.checkExpressionTypes(node.Expression, targetType)
	}
}

// VisitMultiVariableNode checks whether the variable types matches with the function return types
func (v *typeCheckVisitor) VisitMultiVariableNode(node *node.MultiVariableNode) {
	v.AbstractVisitor.VisitMultiVariableNode(node)
	targetTypes := make([]*symbol.TypeSymbol, len(node.Types))

	for i, t := range node.Types {
		targetTypes[i] = v.symbolTable.FindTypeByNode(t)
	}
	v.checkExpressionTypes(node.FuncCall, targetTypes...)
}

// VisitReturnStatementNode checks whether the return types and the values are of the same type
func (v *typeCheckVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(node)
	returnNodes := node.Expressions
	returnSymbols := v.currentFunction.ReturnTypes

	// TODO: Allow funcCall with multiple return values

	if len(returnSymbols) > 0 {
		if len(returnSymbols) != len(returnNodes) {
			v.reportError(node,
				fmt.Sprintf("Expected %d return values, given %d", len(returnSymbols), len(returnNodes)))
		} else {
			for i, rtype := range returnSymbols {
				nodeType := v.symbolTable.GetTypeByExpression(returnNodes[i])
				if nodeType != rtype {
					v.reportError(node, fmt.Sprintf("Return Type mismatch: expected %s, given %s",
						rtype.ID, nodeType.ID))
				}
			}
		}
	} else if len(returnNodes) > 0 {
		v.reportError(node, "void method should not return expression")
	}
}

// VisitAssignmentStatementNode checks whether the left and right part of the assignment are of the same type
func (v *typeCheckVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.AbstractVisitor.VisitAssignmentStatementNode(node)

	leftType := v.symbolTable.GetTypeByExpression(node.Left)
	rightType := v.symbolTable.GetTypeByExpression(node.Right)

	if leftType != rightType {
		v.reportError(node,
			fmt.Sprintf("assignment of %s is not compatible with target %s", rightType, leftType))
	}
}

// VisitIfStatementNode checks whether the condition is a boolean expression
func (v *typeCheckVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	v.AbstractVisitor.VisitIfStatementNode(node)
	if !v.isBool(v.symbolTable.GetTypeByExpression(node.Condition)) {
		v.reportError(node, "condition must return boolean")
	}
}

func (v *typeCheckVisitor) VisitCallStatementNode(node *node.CallStatementNode) {
	v.AbstractVisitor.VisitCallStatementNode(node)
	if v.symbolTable.GetTypeByExpression(node.Call) != nil {
		v.reportError(node, "function call as statement should be void")
	}
}

// Expressions
// -----------

// VisitBinaryExpressionNode checks if the types for different binary expressions match
// Expressions are &&, ||, +, -, *, /, %, **, ==, !=, >, >=, <= and <
func (v *typeCheckVisitor) VisitBinaryExpressionNode(node *node.BinaryExpressionNode) {
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
			v.reportError(node, "Arithmetic operators can only be applied to int types")
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
func (v *typeCheckVisitor) VisitUnaryExpressionNode(node *node.UnaryExpressionNode) {
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

func (v *typeCheckVisitor) VisitFuncCallNode(funcCallNode *node.FuncCallNode) {
	v.AbstractVisitor.VisitFuncCallNode(funcCallNode)
	funcSym, ok := v.symbolTable.GetDeclByDesignator(funcCallNode.Designator).(*symbol.FunctionSymbol)

	if !ok {
		v.reportError(funcCallNode, fmt.Sprintf("%s is not a function", funcCallNode.Designator))
		return
	}

	totalParams := len(funcSym.Parameters)
	totalArgs := len(funcCallNode.Args)
	if len(funcSym.ReturnTypes) > 0 {
		if totalParams != totalArgs {
			v.reportError(funcCallNode, fmt.Sprintf("expected %d args, got %d", totalParams, totalArgs))
		} else {
			for i, arg := range funcCallNode.Args {
				v.checkType(arg, funcSym.Parameters[i].Type)
			}
		}
	}

	// Function with multiple return values are allowed only in multi-variable, multi-assignment and return statements.
	// Otherwise, the function call should have only one return type.
	// Void function has no type.
	if len(funcSym.ReturnTypes) == 1 {
		v.symbolTable.MapExpressionToType(funcCallNode, funcSym.ReturnTypes[0])
	}
}

// VisitTypeNode currently does nothing
func (v *typeCheckVisitor) VisitTypeNode(node *node.TypeNode) {
	// To be done as soon as own types are introduced
}

// VisitIntegerLiteralNode maps the integer literal node to its type
func (v *typeCheckVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
}

// VisitBoolLiteralNode maps the bool literal node to its type
func (v *typeCheckVisitor) VisitBoolLiteralNode(node *node.BoolLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
}

// VisitStringLiteralNode maps the string literal to its type
func (v *typeCheckVisitor) VisitStringLiteralNode(node *node.StringLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.StringType)
}

// VisitCharacterLiteralNode maps the character literal to its type
func (v *typeCheckVisitor) VisitCharacterLiteralNode(node *node.CharacterLiteralNode) {
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.CharType)
}

// Helper Functions
// ----------------

func (v *typeCheckVisitor) isInt(symbol *symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.IntType
}

func (v *typeCheckVisitor) isBool(symbol *symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.BoolType
}

func (v *typeCheckVisitor) isChar(symbol *symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.CharType
}

func (v *typeCheckVisitor) checkType(expr node.ExpressionNode, expectedType *symbol.TypeSymbol) {
	actualType := v.symbolTable.GetTypeByExpression(expr)
	if expectedType != actualType {
		v.reportError(expr, fmt.Sprintf("expected %s, got %s", expectedType, actualType))
	}
}

func (v *typeCheckVisitor) checkExpressionTypes(expr node.ExpressionNode, expectedTypes ...*symbol.TypeSymbol) {
	// Only function call are allowed to have multiple types
	if fc, ok := expr.(*node.FuncCallNode); ok {
		calledFuncSym, ok := v.symbolTable.GetDeclByDesignator(fc.Designator).(*symbol.FunctionSymbol)

		if !ok {
			v.reportError(fc, fmt.Sprintf("%s is not a function", fc.Designator))
			return
		}

		if len(calledFuncSym.ReturnTypes) != len(expectedTypes) {
			v.reportError(expr,
				fmt.Sprintf("expected %d return value(s), but function returns %d",
					len(expectedTypes), len(calledFuncSym.ReturnTypes)))
			return
		}

		for i, returnType := range calledFuncSym.ReturnTypes {
			if expectedTypes[i] != returnType {
				v.reportError(expr, fmt.Sprintf("Return type mismatch: expected %s, given %s",
					returnType.ID, expectedTypes[i].ID))
			}
		}
		return
	}

	if len(expectedTypes) > 1 {
		v.reportError(expr, "only single type is allowed")
		return
	}

	exprType := v.symbolTable.GetTypeByExpression(expr)
	if exprType != expectedTypes[0] {
		v.reportError(expr, fmt.Sprintf("Type mismatch: expected %s, given %s", expectedTypes[0], exprType))
	}
}

func (v *typeCheckVisitor) reportError(node node.Node, msg string) {
	v.Errors = append(v.Errors, fmt.Errorf("[%s] %s", node.Pos(), msg))
}
