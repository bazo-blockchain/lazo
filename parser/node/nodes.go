package node

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"math/big"
)

type Node interface {
	Pos() token.Position
	String() string
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

// --------------------------

type TypeNode struct {
	AbstractNode
	Identifier string
}

func (n *TypeNode) String() string {
	return fmt.Sprintf("[%s] TYPE %s", n.Pos(), n.Identifier)
}

// --------------------------

type IfStatementNode struct {
	AbstractNode
	Condition ExpressionNode
	Then      *StatementBlockNode
	Else      *StatementBlockNode
}

func (n *IfStatementNode) String() string {
	return fmt.Sprintf("[%s] IF %s THEN %s ELSE %s", n.Pos(), n.Condition, n.Then, n.Else)
}

// --------------------------

type ReturnStatementNode struct {
	AbstractNode
	Expression ExpressionNode
}

func (n *ReturnStatementNode) String() string {
	return fmt.Sprintf("[%s] RETURNSTMT %s", n.Pos(), n.Expression)
}

// --------------------------

type StatementBlockNode struct {
	AbstractNode
	Statements []StatementNode
}

func (n *StatementBlockNode) String() string {
	return fmt.Sprintf("[%s] STMTBLOCK %s", n.Pos(), n.Statements)
}

// --------------------------
// Expression Nodes
// --------------------------

type BinaryExpressionNode struct {
	AbstractNode
	LeftExpr  ExpressionNode
	Operator  token.Symbol
	RightExpr ExpressionNode
}

func (n *BinaryExpressionNode) String() string {
	return fmt.Sprintf("[%s] EXPR (%s %s %s)", n.Position, n.LeftExpr, token.SymbolLexeme[n.Operator], n.RightExpr)
}

// --------------------------

type DesignatorNode struct {
	AbstractNode
	Value string
}

func (n *DesignatorNode) String() string {
	return fmt.Sprintf("[%s] DESIGNATOR %s", n.Pos(), n.Value)
}

// --------------------------
// Literal Nodes
// --------------------------

type IntegerLiteralNode struct {
	AbstractNode
	Value *big.Int
}

func (n *IntegerLiteralNode) String() string {
	return fmt.Sprintf("[%s] INT %d", n.Pos(), n.Value)
}

// --------------------------

type StringLiteralNode struct {
	AbstractNode
	Value string
}

func (n *StringLiteralNode) String() string {
	return fmt.Sprintf("[%s] STRING %s", n.Pos(), n.Value)
}

// --------------------------

type CharacterLiteralNode struct {
	AbstractNode
	Value rune
}

func (n *CharacterLiteralNode) String() string {
	return fmt.Sprintf("[%s] CHARACTER %c", n.Pos(), n.Value)
}

// --------------------------

type BoolLiteralNode struct {
	AbstractNode
	Value bool
}

func (n *BoolLiteralNode) String() string {
	return fmt.Sprintf("[%s] BOOL %t", n.Pos(), n.Value)
}

// --------------------------

type ErrorNode struct {
	AbstractNode
	Message string
}

func (n *ErrorNode) String() string {
	return fmt.Sprintf("[%s] ERROR: %s", n.Pos(), n.Message)
}
