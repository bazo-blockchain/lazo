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

// Concrete Nodes
// -------------------------

// ProgramNode composes abstract node and holds contract.
type ProgramNode struct {
	AbstractNode
	Contract *ContractNode
}

func (n *ProgramNode) String() string {
	return fmt.Sprintf("%s", getNodeString(n.Contract))
}

// Accept lets a visitor to traverse its node structure.
func (n *ProgramNode) Accept(v Visitor) {
	v.VisitProgramNode(n)
}

// --------------------------

// ContractNode composes abstract node and holds a name, state variables and functions.
type ContractNode struct {
	AbstractNode
	Name      string
	Fields    []*FieldNode
	Functions []*FunctionNode
}

func (n *ContractNode) String() string {
	return fmt.Sprintf("[%s] CONTRACT %s \n FIELDS: %s \n\n FUNCS: %s", n.Pos(), n.Name, n.Fields, n.Functions)
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
	Type       *TypeNode
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

// FunctionNode composes abstract node and holds a name, return types, parameters and statements.
type FunctionNode struct {
	AbstractNode
	Name        string
	ReturnTypes []*TypeNode
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
	Type       *TypeNode
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
	Type       *TypeNode
	Identifier string
	Expression ExpressionNode
}

func (n *VariableNode) String() string {
	var str string
	if n.Expression == nil {
		str = fmt.Sprintf("\n [%s] VARIABLE %s %s", n.Pos(), getNodeString(n.Type), n.Identifier)
	} else {
		str = fmt.Sprintf("\n [%s] VARIABLE %s %s = %s", n.Pos(), getNodeString(n.Type), n.Identifier, getNodeString(n.Expression))
	}
	return str
}

// Accept lets a visitor to traverse its node structure
func (n *VariableNode) Accept(v Visitor) {
	v.VisitVariableNode(n)
}

// --------------------------

// TypeNode composes abstract node and holds the type identifier.
type TypeNode struct {
	AbstractNode
	Identifier string
}

func (n *TypeNode) String() string {
	return fmt.Sprintf("%s", n.Identifier)
}

// Accept lets a visitor to traverse its node structure
func (n *TypeNode) Accept(v Visitor) {
	v.VisitTypeNode(n)
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
	Left  *DesignatorNode
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
// Expression Nodes
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

// DesignatorNode composes abstract node and holds the designator name.
type DesignatorNode struct {
	AbstractNode
	Value string
}

func (n *DesignatorNode) String() string {
	return fmt.Sprintf("%s", n.Value)
}

// Accept lets a visitor to traverse its node structure.
func (n *DesignatorNode) Accept(v Visitor) {
	v.VisitDesignatorNode(n)
}

// --------------------------

// FuncCallNode compose abstract node and holds designator and arguments
type FuncCallNode struct {
	AbstractNode
	Designator *DesignatorNode
	Args       []ExpressionNode
}

func (n *FuncCallNode) String() string {
	return fmt.Sprintf("%s(%s)", n.Designator, n.Args)
}

// Accept lets a visitor to traverse its node structure
func (n *FuncCallNode) Accept(v Visitor) {
	v.VisitFuncCallNode(n)
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
