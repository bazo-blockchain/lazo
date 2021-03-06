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

	if node.Constructor != nil {
		v.currentFunction = v.contractSymbol.Constructor
		node.Constructor.Accept(v.ConcreteVisitor)
		v.currentFunction = nil
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
	targetTypes := make([]symbol.TypeSymbol, len(node.Types))

	for i, t := range node.Types {
		targetTypes[i] = v.symbolTable.FindTypeByNode(t)
	}
	v.checkExpressionTypes(node.FuncCall, targetTypes...)
}

// VisitReturnStatementNode checks whether the return types and the values are of the same type
func (v *typeCheckVisitor) VisitReturnStatementNode(returnNode *node.ReturnStatementNode) {
	v.AbstractVisitor.VisitReturnStatementNode(returnNode)

	if v.contractSymbol.Constructor == v.currentFunction {
		v.reportError(returnNode, "return is not allowed in constructor")
		return
	}

	returnNodeExpressions := returnNode.Expressions
	returnSymbols := v.currentFunction.ReturnTypes

	// if funcCall is returned, check the return types of the called function
	if len(returnNodeExpressions) == 1 {
		if fc, ok := returnNodeExpressions[0].(*node.FuncCallNode); ok {
			v.checkExpressionTypes(fc, returnSymbols...)
			return
		}
	}

	if len(returnSymbols) != len(returnNodeExpressions) {
		v.reportError(returnNode,
			fmt.Sprintf("Expected %d return values, given %d", len(returnSymbols), len(returnNodeExpressions)))
		return
	}

	if len(returnSymbols) > 0 {
		for i, rtype := range returnSymbols {
			if returnNodeExpressions[i].String() == symbol.This {
				v.reportError(returnNode, "'this' cannot be returned")
				return
			}
			nodeType := v.symbolTable.GetTypeByExpression(returnNodeExpressions[i])
			if nodeType != rtype {
				v.reportError(returnNode, fmt.Sprintf("Return type mismatch: expected %s, given %s",
					rtype.Identifier(), getTypeString(nodeType)))
			}
		}
	} else if len(returnNodeExpressions) > 0 {
		v.reportError(returnNode, "void method should not return expression")
	}
}

// VisitAssignmentStatementNode checks whether the left and right part of the assignment are of the same type
func (v *typeCheckVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	v.AbstractVisitor.VisitAssignmentStatementNode(node)

	if node.Left.String() == symbol.This {
		v.reportError(node, "Assigning to 'this' is not allowed!")
		return
	}

	if node.Right.String() == symbol.This {
		v.reportError(node, "'this' cannot be assigned!")
		return
	}

	leftType := v.symbolTable.GetTypeByExpression(node.Left)
	rightType := v.symbolTable.GetTypeByExpression(node.Right)

	if leftType != rightType {
		v.reportError(node,
			fmt.Sprintf("assignment of %s is not compatible with target %s",
				getTypeString(rightType), getTypeString(leftType)))
	}
}

// VisitMultiAssignmentStatementNode checks whether the target designators matches with the function return types.
func (v *typeCheckVisitor) VisitMultiAssignmentStatementNode(node *node.MultiAssignmentStatementNode) {
	v.AbstractVisitor.VisitMultiAssignmentStatementNode(node)

	leftTypes := make([]symbol.TypeSymbol, len(node.Designators))
	for i, designator := range node.Designators {
		if designator.String() == symbol.This {
			v.reportError(node, "Assigning to 'this' is not allowed!")
		}
		leftTypes[i] = v.symbolTable.GetTypeByExpression(designator)
	}
	v.checkExpressionTypes(node.FuncCall, leftTypes...)
}

func (v *typeCheckVisitor) VisitShorthandAssignmentNode(node *node.ShorthandAssignmentStatementNode) {
	v.AbstractVisitor.VisitShorthandAssignmentNode(node)

	designatorType := v.symbolTable.GetTypeByExpression(node.Designator)

	// str += "hello"
	if v.isString(designatorType) && node.Operator == token.Plus {
		v.checkType(node.Expression, v.symbolTable.GlobalScope.StringType)
		return
	}

	// x += 1 or x++
	v.checkType(node.Designator, v.symbolTable.GlobalScope.IntType)
	v.checkType(node.Expression, v.symbolTable.GlobalScope.IntType)
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

func (v *typeCheckVisitor) VisitDeleteStatementNode(node *node.DeleteStatementNode) {
	v.AbstractVisitor.VisitDeleteStatementNode(node)

	designatorType := v.symbolTable.GetTypeByExpression(node.Element.Designator)
	if _, ok := designatorType.(*symbol.MapTypeSymbol); !ok {
		v.reportError(node, "delete requires map type")
	}
}

// Expressions
// -----------

func (v *typeCheckVisitor) VisitTernaryExpressionNode(node *node.TernaryExpressionNode) {
	v.AbstractVisitor.VisitTernaryExpressionNode(node)

	conditionType := v.symbolTable.GetTypeByExpression(node.Condition)
	if !v.isBool(conditionType) {
		v.reportError(node.Condition, "condition should be bool type")
	}

	trueExprType := v.symbolTable.GetTypeByExpression(node.Then)
	falseExprType := v.symbolTable.GetTypeByExpression(node.Else)
	if trueExprType != falseExprType {
		v.reportError(node, "ternary expression should return same type")
	} else {
		v.symbolTable.MapExpressionToType(node, trueExprType)
	}
}

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
		v.visitBinaryLogicalOperator(node, leftType, rightType)
	case token.BitwiseAnd, token.BitwiseOr, token.BitwiseXOr:
		v.visitBinaryBitwiseLogicalOperator(node, leftType, rightType)
	case token.Plus:
		v.visitBinaryPlusOperator(node, leftType, rightType)
	case token.Minus, token.Multiplication, token.Division, token.Modulo, token.Exponent:
		v.visitBinaryArithmeticOperator(node, leftType, rightType)
	case token.Equal, token.Unequal:
		v.visitBinaryEqualityComparisonOperator(node, leftType, rightType)
	case token.Less, token.LessEqual, token.GreaterEqual, token.Greater:
		v.visitBinaryRelationalComparisonOperator(node, leftType, rightType)
	case token.ShiftLeft, token.ShiftRight:
		v.visitBinaryShiftOperator(node, leftType, rightType)
	default:
		panic(fmt.Sprintf("Illegal binary operator %s", token.SymbolLexeme[node.Operator]))
	}
}

// visitBinaryLogicalOperator checks && and || operators.
func (v *typeCheckVisitor) visitBinaryLogicalOperator(node *node.BinaryExpressionNode,
	leftType symbol.TypeSymbol, rightType symbol.TypeSymbol) {
	if !v.isBool(leftType) || !v.isBool(rightType) {
		v.reportError(node, "Logic operators can only be applied to bool types")
	}
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
}

// visitBinaryBitwiseLogicalOperator checks &, | and ^ operators.
func (v *typeCheckVisitor) visitBinaryBitwiseLogicalOperator(node *node.BinaryExpressionNode,
	leftType symbol.TypeSymbol, rightType symbol.TypeSymbol) {
	if !v.isInt(leftType) || !v.isInt(rightType) {
		v.reportError(node, "Bitwise logic operators can only be applied to int types")
	}
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
}

// visitBinaryPlusOperator checks addition (1 + 1) and string concatenation ("hello" + "world").
func (v *typeCheckVisitor) visitBinaryPlusOperator(node *node.BinaryExpressionNode,
	leftType symbol.TypeSymbol, rightType symbol.TypeSymbol) {
	if v.isString(leftType) && v.isString(rightType) {
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.StringType)
		return
	}
	if !v.isInt(leftType) || !v.isInt(rightType) {
		v.reportError(node, "+ operator can only be applied to int/string types")
	}
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
}

// visitBinaryArithmeticOperator checks -, *, /, ** operators.
func (v *typeCheckVisitor) visitBinaryArithmeticOperator(node *node.BinaryExpressionNode,
	leftType symbol.TypeSymbol, rightType symbol.TypeSymbol) {
	if !v.isInt(leftType) || !v.isInt(rightType) {
		v.reportError(node, "Arithmetic operators can only be applied to int types")
	}
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
}

// visitBinaryEqualityComparisonOperator checks == and != operators.
func (v *typeCheckVisitor) visitBinaryEqualityComparisonOperator(node *node.BinaryExpressionNode,
	leftType symbol.TypeSymbol, rightType symbol.TypeSymbol) {
	if leftType != rightType {
		v.reportError(node, fmt.Sprintf("Equality comparison should have the same type, given %s and %s",
			leftType, rightType))
	}
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
}

// visitBinaryRelationalComparisonOperator checks comparison with <, <=, >, >= operators.
func (v *typeCheckVisitor) visitBinaryRelationalComparisonOperator(node *node.BinaryExpressionNode,
	leftType symbol.TypeSymbol, rightType symbol.TypeSymbol) {
	if leftType != rightType {
		v.reportError(node,
			fmt.Sprintf("Both sides of a compare operation need to have the same type, given %s and %s",
				leftType, rightType))
	} else if !(v.isInt(leftType) || v.isChar(leftType)) {
		v.reportError(node, fmt.Sprintf("Relational comparison is not supported for %s", leftType))
	}
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
}

// visitBinaryShiftOperator checks << and >> operators.
func (v *typeCheckVisitor) visitBinaryShiftOperator(node *node.BinaryExpressionNode,
	leftType symbol.TypeSymbol, rightType symbol.TypeSymbol) {
	if !v.isInt(leftType) || !v.isInt(rightType) {
		v.reportError(node, "Bitwise shift operators can only be applied to int types")
	}
	v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
}

// VisitUnaryExpressionNode checks that types of unary expressions are valid
// Expressions are +, -, !
func (v *typeCheckVisitor) VisitUnaryExpressionNode(node *node.UnaryExpressionNode) {
	v.AbstractVisitor.VisitUnaryExpressionNode(node)
	operand := node.Expression
	operandType := v.symbolTable.GetTypeByExpression(operand)

	switch node.Operator {
	case token.Plus, token.Minus:
		if !v.isInt(operandType) {
			v.reportError(node, "+ and - unary operators can only be applied to expressions of type int")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
	case token.Not:
		if !v.isBool(operandType) {
			v.reportError(node, "! unary operator can only be applied to expressions of type bool")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.BoolType)
	case token.BitwiseNot:
		if !v.isInt(operandType) {
			v.reportError(node, "~ unary operator can only be applied to int type")
		}
		v.symbolTable.MapExpressionToType(node, v.symbolTable.GlobalScope.IntType)
	default:
		panic(fmt.Sprintf("Illegal unary operator %s", token.SymbolLexeme[node.Operator]))
	}
}

func (v *typeCheckVisitor) VisitTypeCastNode(typeCastNode *node.TypeCastNode) {
	v.AbstractVisitor.VisitTypeCastNode(typeCastNode)

	castType := v.symbolTable.FindTypeByNode(typeCastNode.Type)
	exprType := v.symbolTable.GetTypeByExpression(typeCastNode.Expression)

	gs := v.symbolTable.GlobalScope
	if v.isString(castType) {
		if v.isAnyType(exprType, gs.IntType, gs.CharType, gs.BoolType, gs.StringType) {
			v.symbolTable.MapExpressionToType(typeCastNode, gs.StringType)
		} else {
			v.reportError(typeCastNode, fmt.Sprintf("String type cast is not supported for %s", exprType))
		}
		return
	}

	v.reportError(typeCastNode, fmt.Sprintf("Unsupported type cast to %s", castType))
}

// VisitFuncCallNode checks the types of passed arguments and declared return types.
func (v *typeCheckVisitor) VisitFuncCallNode(funcCallNode *node.FuncCallNode) {
	v.AbstractVisitor.VisitFuncCallNode(funcCallNode)
	funcSym, ok := v.symbolTable.GetDeclByDesignator(funcCallNode.Designator).(*symbol.FunctionSymbol)

	if !ok {
		v.reportError(funcCallNode, fmt.Sprintf("%s is not a function", funcCallNode.Designator))
		return
	}

	// Function with multiple return values are allowed only in multi-variable, multi-assignment and return statements.
	// Otherwise, the function call should have only one return type.
	// Void function has no type.
	if len(funcSym.ReturnTypes) == 1 {
		v.symbolTable.MapExpressionToType(funcCallNode, funcSym.ReturnTypes[0])
	}

	totalParams := len(funcSym.Parameters)
	totalArgs := len(funcCallNode.Args)
	if totalParams != totalArgs {
		v.reportError(funcCallNode, fmt.Sprintf("expected %d args, got %d", totalParams, totalArgs))
		return
	}

	// Check generic map key type
	if funcSym == v.symbolTable.GlobalScope.MapMemberFunctions[symbol.Contains] {
		targetMap := funcCallNode.Designator.(*node.MemberAccessNode).Designator
		targetMapType := v.symbolTable.GetTypeByExpression(targetMap).(*symbol.MapTypeSymbol)
		v.checkType(funcCallNode.Args[0], targetMapType.KeyType)
		return
	}

	for i, arg := range funcCallNode.Args {
		if arg.String() == symbol.This {
			v.reportError(funcCallNode, "'this' cannot be used as an argument")
			return
		}
		v.checkType(arg, funcSym.Parameters[i].Type)
	}
}

// VisitArrayLengthCreationNode checks that the lengths are of type int
func (v *typeCheckVisitor) VisitArrayLengthCreationNode(node *node.ArrayLengthCreationNode) {
	v.AbstractVisitor.VisitArrayLengthCreationNode(node)

	// int[] has element type 'int'
	typeSymbol := v.symbolTable.FindTypeByIdentifier(node.ElementType.String() + "[]")
	if typeSymbol == nil {
		typeSymbol = v.symbolTable.AddArrayType(node.ElementType)
	}
	if typeSymbol == nil {
		v.reportError(node, "Invalid array type")
	}

	v.symbolTable.MapExpressionToType(node, typeSymbol)
	for i, length := range node.Lengths {
		exprType := v.symbolTable.GetTypeByExpression(length)
		if exprType != v.symbolTable.GlobalScope.IntType {
			v.reportError(node.Lengths[i], "Only integer expressions are allowed as array length argument")
		}
	}
}

// VisitArrayValueCreationNode checks that each value of
func (v *typeCheckVisitor) VisitArrayValueCreationNode(node *node.ArrayValueCreationNode) {
	typeSymbol := v.symbolTable.FindTypeByNode(node.Type)
	if typeSymbol == nil {
		v.reportError(node, "Invalid array type")
		return
	}
	v.symbolTable.MapExpressionToType(node, typeSymbol)

	arrayTypeSymbol := typeSymbol.(*symbol.ArrayTypeSymbol)
	v.symbolTable.MapExpressionToType(node.Elements, arrayTypeSymbol)

	// Visit array values and check element type
	v.AbstractVisitor.VisitArrayValueCreationNode(node)
}

func (v *typeCheckVisitor) VisitArrayInitializationNode(arrayValueInitNode *node.ArrayInitializationNode) {
	arrayType := v.symbolTable.GetTypeByExpression(arrayValueInitNode).(*symbol.ArrayTypeSymbol)

	for _, element := range arrayValueInitNode.Values {
		if arrayInitValues, ok := element.(*node.ArrayInitializationNode); ok {
			// e.g. new int[][]{{1, 2}, {3, 4}} --> Value {1, 2} has array type int[]
			v.symbolTable.MapExpressionToType(arrayInitValues, arrayType.ElementType)
			arrayInitValues.Accept(v) // check value types recursively
		} else {
			element.Accept(v) // resolve expression type

			// e.g. new int[]{1, 2} --> Value 1 has basic type int
			v.checkType(element, arrayType.ElementType)
		}
	}
}

// VisitElementAccessNode checks that the expression is of type integer
func (v *typeCheckVisitor) VisitElementAccessNode(node *node.ElementAccessNode) {
	v.AbstractVisitor.VisitElementAccessNode(node)

	designatorType := v.symbolTable.GetTypeByExpression(node.Designator)
	if _, ok := designatorType.(*symbol.ArrayTypeSymbol); ok {
		if v.symbolTable.GetTypeByExpression(node.Expression) != v.symbolTable.GlobalScope.IntType {
			v.reportError(node, "Array index must be of type int")
		}
	} else if mapType, ok := designatorType.(*symbol.MapTypeSymbol); ok {
		v.checkType(node.Expression, mapType.KeyType)
	} else {
		panic("Unsupported element access designator")
	}
}

// VisitStructCreationNode maps the node to its struct declaration and checks field value types
func (v *typeCheckVisitor) VisitStructCreationNode(node *node.StructCreationNode) {
	v.AbstractVisitor.VisitStructCreationNode(node)

	structType, ok := v.symbolTable.GlobalScope.Structs[node.Name]
	if !ok {
		v.reportError(node, fmt.Sprintf("Struct %s is undefined", node.Name))
		return
	}
	v.symbolTable.MapExpressionToType(node, structType)

	if len(node.FieldValues) > len(structType.Fields) {
		v.reportError(node, fmt.Sprintf("Struct %s has only %d field(s), got %d value(s)",
			node.Name, len(node.FieldValues), len(structType.Fields)))
		return
	}

	for i, fieldValue := range node.FieldValues {
		exprType := v.symbolTable.GetTypeByExpression(fieldValue)
		expectedType := structType.Fields[i].Type
		if exprType != expectedType {
			v.reportError(fieldValue, fmt.Sprintf(typeErrorMsgTemplate, expectedType, exprType))
		}
	}
}

func (v *typeCheckVisitor) VisitStructNamedCreationNode(node *node.StructNamedCreationNode) {
	v.AbstractVisitor.VisitStructNamedCreationNode(node)

	structType, ok := v.symbolTable.GlobalScope.Structs[node.Name]
	if !ok {
		v.reportError(node, fmt.Sprintf("Struct %s is undefined", node.Name))
		return
	}
	v.symbolTable.MapExpressionToType(node, structType)

	if len(node.FieldValues) > len(structType.Fields) {
		v.reportError(node, fmt.Sprintf("Struct %s has only %d field(s), got %d value(s)",
			node.Name, len(node.FieldValues), len(structType.Fields)))
		return
	}

	for _, fieldValue := range node.FieldValues {
		exprType := v.symbolTable.GetTypeByExpression(fieldValue)
		fieldSymbol := structType.GetField(fieldValue.Name)
		if fieldSymbol == nil {
			v.reportError(fieldValue, fmt.Sprintf("Field %s not found", fieldValue.Name))
		} else if exprType != fieldSymbol.Type {
			v.reportError(fieldValue, fmt.Sprintf(typeErrorMsgTemplate, fieldSymbol.Type, exprType))
		}
	}
}

// VisitStructFieldAssignmentNode traverse the field initialization expression
func (v *typeCheckVisitor) VisitStructFieldAssignmentNode(node *node.StructFieldAssignmentNode) {
	v.AbstractVisitor.VisitStructFieldAssignmentNode(node)
	exprType := v.symbolTable.GetTypeByExpression(node.Expression)
	v.symbolTable.MapExpressionToType(node, exprType)
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

func (v *typeCheckVisitor) isInt(symbol symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.IntType
}

func (v *typeCheckVisitor) isBool(symbol symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.BoolType
}

func (v *typeCheckVisitor) isChar(symbol symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.CharType
}

func (v *typeCheckVisitor) isString(symbol symbol.TypeSymbol) bool {
	return symbol == v.symbolTable.GlobalScope.StringType
}

func (v *typeCheckVisitor) isAnyType(symbol symbol.TypeSymbol, expectedTypes ...symbol.TypeSymbol) bool {
	for _, t := range expectedTypes {
		if t == symbol {
			return true
		}
	}
	return false
}

func (v *typeCheckVisitor) checkType(expr node.ExpressionNode, expectedType symbol.TypeSymbol) {
	actualType := v.symbolTable.GetTypeByExpression(expr)
	if expectedType != actualType {
		v.reportError(expr, fmt.Sprintf(typeErrorMsgTemplate, expectedType, actualType))
	}
}

func (v *typeCheckVisitor) checkExpressionTypes(expr node.ExpressionNode, expectedTypes ...symbol.TypeSymbol) {
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
					returnType.Identifier(), expectedTypes[i].Identifier()))
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
		v.reportError(expr, fmt.Sprintf("Type mismatch: expected %s, given %s",
			expectedTypes[0].Identifier(), getTypeString(exprType)))
	}
}

const typeErrorMsgTemplate = "expected %s, got %s"

func (v *typeCheckVisitor) reportError(node node.Node, msg string) {
	v.Errors = append(v.Errors, fmt.Errorf("[%s] %s", node.Pos(), msg))
}

func getTypeString(t symbol.TypeSymbol) string {
	if t == nil {
		return "nil"
	}
	return t.Identifier()
}
