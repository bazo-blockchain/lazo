// Package node contains all the supported node types and their functions.
package node

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"math/big"
)

// Node is the interface that wraps the basic Node functions.
type Node interface {
	// Pos returns the position of the node in the source code.
	// It is also the position of the first token.
	Pos() token.Position
	// String returns a readable string representation of the node.
	String() string
	// Accept lets a visitor to traverse its node structure.
	Accept(v Visitor)
}

// AbstractNode contains node position, which all concrete nodes have.
type AbstractNode struct {
	Position token.Position
}

// Pos returns the node position
func (n *AbstractNode) Pos() token.Position {
	return n.Position
}

// StatementNode is the interface for statements, such as variable, assignment, if-statement etc.
type StatementNode interface {
	Node
}

// ExpressionNode is the interface for expressions, such as literal, identifier, binary expression etc.
type ExpressionNode interface {
	Node
}

// DesignatorNode is the interface for designators, such as identifier, member access and array access.
type DesignatorNode interface {
	Node
}

// TypeNode is the interface for types, such as array types or basic types
type TypeNode interface {
	Node
	Type() string
}

// Concrete Nodes
// -------------------------

// ProgramNode composes abstract node and holds contract.
type ProgramNode struct {
	AbstractNode
	Contract *ContractNode
}

func (n *ProgramNode) String() string {
	return getNodeString(n.Contract)
}

// Accept lets a visitor to traverse its node structure.
func (n *ProgramNode) Accept(v Visitor) {
	v.VisitProgramNode(n)
}

// --------------------------

// ContractNode composes abstract node and holds a name, state variables and functions.
type ContractNode struct {
	AbstractNode
	Name        string
	Fields      []*FieldNode
	Structs     []*StructNode
	Constructor *ConstructorNode
	Functions   []*FunctionNode
}

func (n *ContractNode) String() string {
	var strConstructor string
	if n.Constructor != (*ConstructorNode)(nil) {
		strConstructor = n.Constructor.String()
	}

	return fmt.Sprintf("[%s] CONTRACT %s \n FIELDS: %s \n\n STRUCTS: %s \n\n CONSTRUCTOR: %s \n\n FUNCS: %s",
		n.Pos(), n.Name, n.Fields, n.Structs, strConstructor, n.Functions)
}

// Accept lets a visitor to traverse its node structure
func (n *ContractNode) Accept(v Visitor) {
	v.VisitContractNode(n)
}

// --------------------------
// Contract Body Parts
// --------------------------

// FieldNode composes abstract node and holds the type, identifier and expression
type FieldNode struct {
	AbstractNode
	Type       TypeNode
	Identifier string
	Expression ExpressionNode
}

func (n *FieldNode) String() string {
	str := fmt.Sprintf("\n [%s] FIELD %s %s", n.Pos(), getNodeString(n.Type), n.Identifier)
	if n.Expression != nil {
		str += fmt.Sprintf(" = %s", getNodeString(n.Expression))
	}
	return str
}

// Accept lets a visitor to traverse its node structure
func (n *FieldNode) Accept(v Visitor) {
	v.VisitFieldNode(n)
}

// --------------------------

// StructNode composes abstract node
type StructNode struct {
	AbstractNode
	Name   string
	Fields []*StructFieldNode
}

func (n *StructNode) String() string {
	return fmt.Sprintf("\n [%s] STRUCT %s \n FIELDS: %s", n.Pos(), n.Name, n.Fields)
}

// Accept lets a visitor to traverse its node structure
func (n *StructNode) Accept(v Visitor) {
	v.VisitStructNode(n)
}

// --------------------------

// StructFieldNode composes abstract node and holds the type and identifier
type StructFieldNode struct {
	AbstractNode
	Type       TypeNode
	Identifier string
}

func (n *StructFieldNode) String() string {
	return fmt.Sprintf("\n [%s] FIELD %s %s", n.Pos(), getNodeString(n.Type), n.Identifier)
}

// Accept lets a visitor to traverse its node structure
func (n *StructFieldNode) Accept(v Visitor) {
	v.VisitStructFieldNode(n)
}

// --------------------------

// ConstructorNode composes abstract node and holds parameters and statements.
type ConstructorNode struct {
	AbstractNode
	Parameters []*ParameterNode
	Body       []StatementNode
}

func (n *ConstructorNode) String() string {
	return fmt.Sprintf("\n [%s] CONSTRUCTOR PARAMs %s, %s",
		n.Pos(), n.Parameters, n.Body)
}

// Accept lets a visitor to traverse its node structure
func (n *ConstructorNode) Accept(v Visitor) {
	v.VisitConstructorNode(n)
}

// --------------------------

// FunctionNode composes abstract node and holds a name, return types, parameters and statements.
type FunctionNode struct {
	AbstractNode
	Name        string
	ReturnTypes []TypeNode
	Parameters  []*ParameterNode
	Body        []StatementNode
}

func (n *FunctionNode) String() string {
	return fmt.Sprintf("\n [%s] FUNCTION %s, PARAMs %s, RTYPES %s %s",
		n.Pos(), n.Name, n.Parameters, n.ReturnTypes, n.Body)
}

// Accept lets a visitor to traverse its node structure
func (n *FunctionNode) Accept(v Visitor) {
	v.VisitFunctionNode(n)
}

// --------------------------

// ParameterNode composes abstract node and holds the type and identifier
type ParameterNode struct {
	AbstractNode
	Type       TypeNode
	Identifier string
}

func (n *ParameterNode) String() string {
	return fmt.Sprintf("\n [%s] PARAM %s %s", n.Pos(), getNodeString(n.Type), n.Identifier)
}

// Accept lets a visitor to traverse its node structure
func (n *ParameterNode) Accept(v Visitor) {
	v.VisitParameterNode(n)
}

// --------------------------
// Statement Nodes
// --------------------------

// VariableNode composes abstract node and holds the type, identifier and expression
type VariableNode struct {
	AbstractNode
	Type       TypeNode
	Identifier string
	Expression ExpressionNode
}

func (n *VariableNode) String() string {
	str := fmt.Sprintf("\n [%s] VAR %s %s", n.Pos(), getNodeString(n.Type), n.Identifier)
	if n.Expression != nil {
		str += fmt.Sprintf(" = %s", getNodeString(n.Expression))
	}
	return str
}

// Accept lets a visitor to traverse its node structure
func (n *VariableNode) Accept(v Visitor) {
	v.VisitVariableNode(n)
}

// --------------------------

// MultiVariableNode composes abstract node and holds multiple variables and a function call
type MultiVariableNode struct {
	AbstractNode
	Types       []TypeNode
	Identifiers []string
	FuncCall    *FuncCallNode
}

func (n *MultiVariableNode) String() string {
	str := fmt.Sprintf("\n [%s] VARS", n.Pos())
	for i, id := range n.Identifiers {
		str += fmt.Sprintf(" %s %s", n.Types[i], id)
	}
	str += fmt.Sprintf(" = %s", getNodeString(n.FuncCall))
	return str
}

// GetType returns the type of the given variable identifier
func (n *MultiVariableNode) GetType(id string) TypeNode {
	for i, varID := range n.Identifiers {
		if id == varID {
			return n.Types[i]
		}
	}
	return nil
}

// Accept lets a visitor to traverse its node structure
func (n *MultiVariableNode) Accept(v Visitor) {
	v.VisitMultiVariableNode(n)
}

// --------------------------

// BasicTypeNode composes abstract node and holds the type identifier.
type BasicTypeNode struct {
	AbstractNode
	Identifier string
}

func (n *BasicTypeNode) String() string {
	return n.Identifier
}

// Accept lets a visitor to traverse its node structure
func (n *BasicTypeNode) Accept(v Visitor) {
	v.VisitBasicTypeNode(n)
}

// Type returns the unique type representation
func (n *BasicTypeNode) Type() string {
	return n.Identifier
}

// --------------------------

// ArrayTypeNode composes abstract node and holds the type identifier.
type ArrayTypeNode struct {
	AbstractNode
	ElementType TypeNode
}

func (n *ArrayTypeNode) String() string {
	return fmt.Sprintf("%s[]", n.ElementType)
}

// Accept lets a visitor to traverse its node structure
func (n *ArrayTypeNode) Accept(v Visitor) {
	v.VisitArrayTypeNode(n)
}

// Type returns the unique type representation
func (n *ArrayTypeNode) Type() string {
	return n.String()
}

// --------------------------

// MapTypeNode composes abstract node and holds the types of key and value.
type MapTypeNode struct {
	AbstractNode
	KeyType   TypeNode
	ValueType TypeNode
}

func (n *MapTypeNode) String() string {
	return fmt.Sprintf("Map<%s,%s>", n.KeyType, n.ValueType)
}

// Accept lets a visitor to traverse its node structure
func (n *MapTypeNode) Accept(v Visitor) {
	v.VisitMapTypeNode(n)
}

// Type returns the unique type representation
func (n *MapTypeNode) Type() string {
	return n.String()
}

// --------------------------

// IfStatementNode composes abstract node and holds the condition, then and else statement block.
type IfStatementNode struct {
	AbstractNode
	Condition ExpressionNode
	Then      []StatementNode
	Else      []StatementNode
}

func (n *IfStatementNode) String() string {
	return fmt.Sprintf("\n [%s] IF %s THEN %s ELSE %s", n.Pos(), getNodeString(n.Condition), n.Then, n.Else)
}

// Accept lets a visitor to traverse its node structure
func (n *IfStatementNode) Accept(v Visitor) {
	v.VisitIfStatementNode(n)
}

// --------------------------

// ReturnStatementNode composes abstract node and holds the return expressions.
type ReturnStatementNode struct {
	AbstractNode
	Expressions []ExpressionNode
}

func (n *ReturnStatementNode) String() string {
	return fmt.Sprintf("\n [%s] RETURNSTMT %s", n.Pos(), n.Expressions)
}

// Accept lets a visitor to traverse its node structure
func (n *ReturnStatementNode) Accept(v Visitor) {
	v.VisitReturnStatementNode(n)
}

// --------------------------

// AssignmentStatementNode composes abstract node and holds the target designator and value expression.
type AssignmentStatementNode struct {
	AbstractNode
	Left  DesignatorNode
	Right ExpressionNode
}

func (n *AssignmentStatementNode) String() string {
	return fmt.Sprintf("\n [%s] ASSIGN %s %s", n.Pos(), getNodeString(n.Left), getNodeString(n.Right))
}

// Accept lets a visitor to traverse its node structure
func (n *AssignmentStatementNode) Accept(v Visitor) {
	v.VisitAssignmentStatementNode(n)
}

// --------------------------

// MultiAssignmentStatementNode composes abstract node and holds the target designators and a function call
type MultiAssignmentStatementNode struct {
	AbstractNode
	Designators []DesignatorNode
	FuncCall    *FuncCallNode
}

func (n *MultiAssignmentStatementNode) String() string {
	str := fmt.Sprintf("\n [%s] VARS", n.Pos())
	for _, id := range n.Designators {
		str += fmt.Sprintf(" %s", id)
	}
	str += fmt.Sprintf(" = %s", getNodeString(n.FuncCall))
	return str
}

// Accept lets a visitor to traverse its node structure
func (n *MultiAssignmentStatementNode) Accept(v Visitor) {
	v.VisitMultiAssignmentStatementNode(n)
}

// --------------------------

// ShorthandAssignmentStatementNode composes abstract node and holds the designator, operator and expression.
type ShorthandAssignmentStatementNode struct {
	AbstractNode
	Designator DesignatorNode
	Operator   token.Symbol
	Expression ExpressionNode
}

func (n *ShorthandAssignmentStatementNode) String() string {
	return fmt.Sprintf("\n %s %s=%s", n.Designator, token.SymbolLexeme[n.Operator], n.Expression)
}

// Accept lets a visitor to traverse its node structure
func (n *ShorthandAssignmentStatementNode) Accept(v Visitor) {
	v.VisitShorthandAssignmentNode(n)
}

// --------------------------

// CallStatementNode composes abstract node and holds the function call expression
type CallStatementNode struct {
	AbstractNode
	Call *FuncCallNode
}

func (n *CallStatementNode) String() string {
	return fmt.Sprintf("\n [%s] CALL %s", n.Pos(), getNodeString(n.Call))
}

// Accept lets a visitor to traverse its node structure
func (n *CallStatementNode) Accept(v Visitor) {
	v.VisitCallStatementNode(n)
}

// --------------------------

// DeleteStatementNode composes abstract node and holds the map element to be deleted
type DeleteStatementNode struct {
	AbstractNode
	Element *ElementAccessNode
}

func (n *DeleteStatementNode) String() string {
	return fmt.Sprintf("\n delete %s", n.Element)
}

// Accept lets a visitor to traverse its node structure
func (n *DeleteStatementNode) Accept(v Visitor) {
	v.VisitDeleteStatementNode(n)
}

// --------------------------
// Expression Nodes
// --------------------------

// TernaryExpressionNode composes abstract node and holds the binary operator and left & right expressions.
type TernaryExpressionNode struct {
	AbstractNode
	Condition ExpressionNode
	Then      ExpressionNode
	Else      ExpressionNode
}

func (n *TernaryExpressionNode) String() string {
	return fmt.Sprintf("%s ? %s : %s", n.Condition, n.Then, n.Else)
}

// Accept lets a visitor to traverse its node structure
func (n *TernaryExpressionNode) Accept(v Visitor) {
	v.VisitTernaryExpressionNode(n)
}

// --------------------------

// BinaryExpressionNode composes abstract node and holds the binary operator and left & right expressions.
type BinaryExpressionNode struct {
	AbstractNode
	Left     ExpressionNode
	Operator token.Symbol
	Right    ExpressionNode
}

func (n *BinaryExpressionNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left, token.SymbolLexeme[n.Operator], getNodeString(n.Right))
}

// Accept lets a visitor to traverse its node structure
func (n *BinaryExpressionNode) Accept(v Visitor) {
	v.VisitBinaryExpressionNode(n)
}

// --------------------------

// UnaryExpressionNode composes abstract node and holds the type, identifier and expression
type UnaryExpressionNode struct {
	AbstractNode
	Operator   token.Symbol
	Expression ExpressionNode
}

func (n *UnaryExpressionNode) String() string {
	return fmt.Sprintf("(%s%s)", token.SymbolLexeme[n.Operator], getNodeString(n.Expression))
}

// Accept lets a visitor to traverse its node structure
func (n *UnaryExpressionNode) Accept(v Visitor) {
	v.VisitUnaryExpressionNode(n)
}

// --------------------------

// TypeCastNode composes abstract node and holds the type and the designator
type TypeCastNode struct {
	AbstractNode
	Type       *BasicTypeNode
	Expression ExpressionNode
}

func (n *TypeCastNode) String() string {
	return fmt.Sprintf("(%s) %s", n.Type, n.Expression)
}

// Accept lets a visitor to traverse its node structure
func (n *TypeCastNode) Accept(v Visitor) {
	v.VisitTypeCastNode(n)
}

// --------------------------

// BasicDesignatorNode composes abstract node and holds the designator name.
type BasicDesignatorNode struct {
	AbstractNode
	Value string
}

func (n *BasicDesignatorNode) String() string {
	return n.Value
}

// Accept lets a visitor to traverse its node structure.
func (n *BasicDesignatorNode) Accept(v Visitor) {
	v.VisitBasicDesignatorNode(n)
}

// --------------------------

// ElementAccessNode composes abstract node and holds designator and expression
type ElementAccessNode struct {
	AbstractNode
	Designator DesignatorNode
	Expression ExpressionNode
}

func (n *ElementAccessNode) String() string {
	return fmt.Sprintf("%s[%s]", n.Designator, n.Expression)
}

// Accept lets a visitor traverse its node structure
func (n *ElementAccessNode) Accept(v Visitor) {
	v.VisitElementAccessNode(n)
}

// --------------------------

//MemberAccessNode composes abstract node and holds designator and identifier
type MemberAccessNode struct {
	AbstractNode
	Designator DesignatorNode
	Identifier string
}

func (n *MemberAccessNode) String() string {
	return fmt.Sprintf("%s.%s", n.Designator, n.Identifier)
}

// Accept lets a visitor traverse its node structure
func (n *MemberAccessNode) Accept(v Visitor) {
	v.VisitMemberAccessNode(n)
}

// --------------------------

// FuncCallNode composes abstract node and holds designator and arguments
type FuncCallNode struct {
	AbstractNode
	Designator DesignatorNode
	Args       []ExpressionNode
}

func (n *FuncCallNode) String() string {
	return fmt.Sprintf("%s(%s)", n.Designator, n.Args)
}

// Accept lets a visitor traverse its node structure
func (n *FuncCallNode) Accept(v Visitor) {
	v.VisitFuncCallNode(n)
}

// --------------------------

// StructCreationNode composes abstract node and holds the target struct and field arguments.
type StructCreationNode struct {
	AbstractNode
	Name        string
	FieldValues []ExpressionNode
}

func (n *StructCreationNode) String() string {
	return fmt.Sprintf("%s(%s)", n.Name, n.FieldValues)
}

// Accept lets a visitor traverse its node structure
func (n *StructCreationNode) Accept(v Visitor) {
	v.VisitStructCreationNode(n)
}

// --------------------------

// StructNamedCreationNode composes abstract node and holds the target struct and field arguments with field name.
type StructNamedCreationNode struct {
	AbstractNode
	Name        string
	FieldValues []*StructFieldAssignmentNode
}

func (n *StructNamedCreationNode) String() string {
	return fmt.Sprintf("%s(%s)", n.Name, n.FieldValues)
}

// Accept lets a visitor traverse its node structure
func (n *StructNamedCreationNode) Accept(v Visitor) {
	v.VisitStructNamedCreationNode(n)
}

// --------------------------

// StructFieldAssignmentNode composes abstract node and holds the target struct field name and expression.
type StructFieldAssignmentNode struct {
	AbstractNode
	Name       string
	Expression ExpressionNode
}

func (n *StructFieldAssignmentNode) String() string {
	return fmt.Sprintf("%s=%s", n.Name, n.Expression)
}

// Accept lets a visitor traverse its node structure
func (n *StructFieldAssignmentNode) Accept(v Visitor) {
	v.VisitStructFieldAssignmentNode(n)
}

// --------------------------

// ArrayLengthCreationNode composes abstract node and holds the target struct and field arguments.
type ArrayLengthCreationNode struct {
	AbstractNode
	ElementType TypeNode
	Lengths     []ExpressionNode
}

func (n *ArrayLengthCreationNode) String() string {
	line := fmt.Sprintf("%s[%s]", n.ElementType, n.Lengths[0])
	for i := 1; i < len(n.Lengths); i++ {
		line = line + fmt.Sprintf("[%s]", n.Lengths[i])
	}

	return line
}

// Accept lets a visitor traverse its node structure
func (n *ArrayLengthCreationNode) Accept(v Visitor) {
	v.VisitArrayLengthCreationNode(n)
}

// --------------------------

// ArrayValueCreationNode composes abstract node and holds the target struct and field arguments.
type ArrayValueCreationNode struct {
	AbstractNode
	Type     TypeNode
	Elements *ArrayInitializationNode
}

func (n *ArrayValueCreationNode) String() string {
	return fmt.Sprintf("%s{%s}", n.Type, n.Elements)
}

// Accept lets a visitor traverse its node structure
func (n *ArrayValueCreationNode) Accept(v Visitor) {
	v.VisitArrayValueCreationNode(n)
}

// --------------------------

// ArrayInitializationNode composes abstract node and holds the target struct and field arguments.
type ArrayInitializationNode struct {
	AbstractNode
	Values []ExpressionNode
}

func (n *ArrayInitializationNode) String() string {
	return fmt.Sprintf("[%s]", n.Values)
}

// Accept lets a visitor traverse its node structure
func (n *ArrayInitializationNode) Accept(v Visitor) {
	v.VisitArrayInitializationNode(n)
}

// --------------------------
// Literal Nodes
// --------------------------

// IntegerLiteralNode composes abstract node and holds the int value.
type IntegerLiteralNode struct {
	AbstractNode
	Value *big.Int
}

func (n *IntegerLiteralNode) String() string {
	return fmt.Sprintf("%d", n.Value)
}

// Accept lets a visitor to traverse its node structure.
func (n *IntegerLiteralNode) Accept(v Visitor) {
	v.VisitIntegerLiteralNode(n)
}

// --------------------------

// StringLiteralNode composes abstract node and holds string literal value.
type StringLiteralNode struct {
	AbstractNode
	Value string
}

func (n *StringLiteralNode) String() string {
	return fmt.Sprintf("%s", n.Value)
}

// Accept lets a visitor to traverse its node structure.
func (n *StringLiteralNode) Accept(v Visitor) {
	v.VisitStringLiteralNode(n)
}

// --------------------------

// CharacterLiteralNode composes abstract node and holds character literal value.
type CharacterLiteralNode struct {
	AbstractNode
	Value rune
}

func (n *CharacterLiteralNode) String() string {
	return fmt.Sprintf("%c", n.Value)
}

// Accept lets a visitor to traverse its node structure
func (n *CharacterLiteralNode) Accept(v Visitor) {
	v.VisitCharacterLiteralNode(n)
}

// --------------------------

// BoolLiteralNode composes abstract node and holds boolean literal value.
type BoolLiteralNode struct {
	AbstractNode
	Value bool
}

func (n *BoolLiteralNode) String() string {
	return fmt.Sprintf("%t", n.Value)
}

// Accept lets a visitor to traverse its node structure
func (n *BoolLiteralNode) Accept(v Visitor) {
	v.VisitBoolLiteralNode(n)
}

// --------------------------

// ErrorNode composes abstract node and holds the syntax error message.
type ErrorNode struct {
	AbstractNode
	Message string
}

func (n *ErrorNode) String() string {
	return fmt.Sprintf("[%s] ERROR: %s", n.Pos(), n.Message)
}

// Accept lets a visitor to traverse its node structure.
func (n *ErrorNode) Accept(v Visitor) {
	v.VisitErrorNode(n)
}

func getNodeString(node Node) string {
	if node == nil {
		return ""
	}
	return node.String()
}
