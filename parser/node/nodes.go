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
	return fmt.Sprintf("%s", n.Contract)
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
	Variables []*VariableNode
	Functions []*FunctionNode
}

func (n *ContractNode) String() string {
	return fmt.Sprintf("[%s] CONTRACT %s \n VARS: %s \n\n FUNCS: %s", n.Pos(), n.Name, n.Variables, n.Functions)
}

// Accept lets a visitor to traverse its node structure
func (n *ContractNode) Accept(v Visitor) {
	v.VisitContractNode(n)
}

// --------------------------
// Contract Body Parts
// --------------------------

// FunctionNode composes abstract node and holds a name, return types, parameters and statements.
type FunctionNode struct {
	AbstractNode
	Name        string
	ReturnTypes []*TypeNode
	Parameters  []*VariableNode
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
	return fmt.Sprintf("\n [%s] VARIABLE %s %s = %s", n.Pos(), n.Type.Identifier, n.Identifier, n.Expression)
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
	return fmt.Sprintf("TYPE %s", n.Identifier)
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
	return fmt.Sprintf("\n [%s] IF %s THEN %s ELSE %s", n.Pos(), n.Condition, n.Then, n.Else)
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
	return fmt.Sprintf("\n [%s] ASSIGN %s %s", n.Pos(), n.Left, n.Right)
}

// Accept lets a visitor to traverse its node structure
func (n *AssignmentStatementNode) Accept(v Visitor) {
	v.VisitAssignmentStatementNode(n)
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
	return fmt.Sprintf("(%s %s %s)", n.Left, token.SymbolLexeme[n.Operator], n.Right)
}

// Accept lets a visitor to traverse its node structure
func (n *BinaryExpressionNode) Accept(v Visitor) {
	v.VisitBinaryExpressionNode(n)
}

// --------------------------

// UnaryExpression composes abstract node and holds the type, identifier and expression
type UnaryExpression struct {
	AbstractNode
	Operator   token.Symbol
	Expression ExpressionNode
}

func (n *UnaryExpression) String() string {
	return fmt.Sprintf("(%s%s)", token.SymbolLexeme[n.Operator], n.Expression)
}

// Accept lets a visitor to traverse its node structure
func (n *UnaryExpression) Accept(v Visitor) {
	v.VisitUnaryExpressionNode(n)
}

// --------------------------

type DesignatorNode struct {
	AbstractNode
	Value string
}

func (n *DesignatorNode) String() string {
	return fmt.Sprintf("%s", n.Value)
}

// Accept lets a visitor to traverse its node structure
func (n *DesignatorNode) Accept(v Visitor) {
	v.VisitDesignatorNode(n)
}

// --------------------------
// Literal Nodes
// --------------------------

type IntegerLiteralNode struct {
	AbstractNode
	Value *big.Int
}

func (n *IntegerLiteralNode) String() string {
	return fmt.Sprintf("%d", n.Value)
}

// Accept lets a visitor to traverse its node structure
func (n *IntegerLiteralNode) Accept(v Visitor) {
	v.VisitIntegerLiteralNode(n)
}

// --------------------------

type StringLiteralNode struct {
	AbstractNode
	Value string
}

func (n *StringLiteralNode) String() string {
	return fmt.Sprintf("%s", n.Value)
}

// Accept lets a visitor to traverse its node structure
func (n *StringLiteralNode) Accept(v Visitor) {
	v.VisitStringLiteralNode(n)
}

// --------------------------

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

type ErrorNode struct {
	AbstractNode
	Message string
}

func (n *ErrorNode) String() string {
	return fmt.Sprintf("[%s] ERROR: %s", n.Pos(), n.Message)
}

// Accept lets a visitor to traverse its node structure
func (n *ErrorNode) Accept(v Visitor) {
	v.VisitErrorNode(n)
}
