package node

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer/token"
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
	Identifier string
	Variables  []*VariableNode
	Functions []*FunctionNode
}

func (n *ContractNode) String() string {
	return fmt.Sprintf("[%s] CONTRACT %s %s", n.Pos(), n.Identifier, n.Variables)
}

// --------------------------
// Statement Nodes
// --------------------------

type VariableNode struct {
	AbstractNode
	Type       *TypeNode
	Identifier string
}

func (n *VariableNode) String() string {
	return fmt.Sprintf("\n [%s] VARIABLE %s %s", n.Pos(), n.Type, n.Identifier)
}

type TypeNode struct {
	AbstractNode
	Identifier string

}

func (n *TypeNode) String() string {
	return fmt.Sprintf("[%s] TYPE %s", n.Pos(), n.Identifier)
}

// --------------------------

// --------------------------
// Expression Nodes
// --------------------------

type FunctionNode struct {
	AbstractNode
	Identifier string
	// TODO Add further members
}

func (n *FunctionNode) String() string {
	// TODO Implement
	return fmt.Sprintf("[%s] Function %s", n.Pos(), n.Identifier)
}

// --------------------------
