package node

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/visitor"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"math/big"
)

type Node interface {
	Pos() token.Position
	String() string
	Accept(v visitor.Visitor)
}

type AbstractNode struct {
	Position token.Position
}

func (n *AbstractNode) Pos() token.Position {
	return n.Position
}

type StatementNode interface {
	Node
}

type ExpressionNode interface {
	Node
}

// Concrete Nodes
// -------------------------

type ProgramNode struct {
	AbstractNode
	Contract *ContractNode
}

func (n *ProgramNode) String() string {
	return fmt.Sprintf("%s", n.Contract)
}

func (n *ProgramNode) Accept(v visitor.Visitor) {
	v.VisitProgramNode(n)
}

// --------------------------

type ContractNode struct {
	AbstractNode
	Name      string
	Variables []*VariableNode
	Functions []*FunctionNode
}

func (n *ContractNode) String() string {
	return fmt.Sprintf("[%s] CONTRACT %s \n VARS: %s \n\n FUNCS: %s", n.Pos(), n.Name, n.Variables, n.Functions)
}

func (n *ContractNode) Accept(v visitor.Visitor) {
	v.VisitContractNode(n)
}

// --------------------------
// Contract Body Parts
// --------------------------

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

func (n *FunctionNode) Accept(v visitor.Visitor) {
	v.VisitFunctionNode(n)
}

// --------------------------
// Statement Nodes
// --------------------------

type VariableNode struct {
	AbstractNode
	Type       *TypeNode
	Identifier string
	Expression ExpressionNode
}

func (n *VariableNode) String() string {
	return fmt.Sprintf("\n [%s] VARIABLE %s %s = %s", n.Pos(), n.Type.Identifier, n.Identifier, n.Expression)
}

func (n *VariableNode) Accept(v visitor.Visitor) {
	v.VisitVariableNode(n)
}

// --------------------------

type TypeNode struct {
	AbstractNode
	Identifier string
}

func (n *TypeNode) String() string {
	return fmt.Sprintf("TYPE %s", n.Identifier)
}

func (n *TypeNode) Accept(v visitor.Visitor) {
	v.VisitTypeNode(n)
}

// --------------------------

type IfStatementNode struct {
	AbstractNode
	Condition ExpressionNode
	Then      []StatementNode
	Else      []StatementNode
}

func (n *IfStatementNode) String() string {
	return fmt.Sprintf("\n [%s] IF %s THEN %s ELSE %s", n.Pos(), n.Condition, n.Then, n.Else)
}

func (n *IfStatementNode) Accept(v visitor.Visitor) {
	v.VisitIfStatementNode(n)
}

// --------------------------

type ReturnStatementNode struct {
	AbstractNode
	Expression []ExpressionNode
}

func (n *ReturnStatementNode) String() string {
	return fmt.Sprintf("\n [%s] RETURNSTMT %s", n.Pos(), n.Expression)
}

func (n *ReturnStatementNode) Accept(v visitor.Visitor) {
	v.VisitAssignmentStatementNode(n)
}

// --------------------------

type AssignmentStatementNode struct {
	AbstractNode
	Left *DesignatorNode
	Right ExpressionNode
}

func (n *AssignmentStatementNode) String() string {
	return fmt.Sprintf("\n [%s] ASSIGN %s %s", n.Pos(), n.Left, n.Right)
}

func (n *AssignmentStatementNode) Accept(v visitor.Visitor) {
	v.VisitAssignmentStatementNode(n)
}

// --------------------------
// Expression Nodes
// --------------------------

type BinaryExpressionNode struct {
	AbstractNode
	Left     ExpressionNode
	Operator token.Symbol
	Right    ExpressionNode
}

func (n *BinaryExpressionNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left, token.SymbolLexeme[n.Operator], n.Right)
}

func (n *BinaryExpressionNode) Accept(v visitor.Visitor) {
	v.VisitBinaryExpressionNode(n)
}

// --------------------------

type UnaryExpression struct {
	AbstractNode
	Operator token.Symbol
	Operand  ExpressionNode
}

func (n *UnaryExpression) String() string {
	return fmt.Sprintf("EXPR (%s %s)", token.SymbolLexeme[n.Operator], n.Operand)
}

func (n *UnaryExpression) Accept(v visitor.Visitor) {
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

func (n *DesignatorNode) Accept(v visitor.Visitor) {
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

func (n *IntegerLiteralNode) Accept(v visitor.Visitor) {
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

func (n *StringLiteralNode) Accept(v visitor.Visitor) {
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

func (n *CharacterLiteralNode) Accept(v visitor.Visitor) {
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

func (n *BoolLiteralNode) Accept(v visitor.Visitor) {
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

func (n *ErrorNode) Accept(v visitor.Visitor) {
	v.VisitErrorNode(n)
}