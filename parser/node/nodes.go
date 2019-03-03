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
	Contract ContractNode
}

func (n *ProgramNode) String() string {
	return fmt.Sprintf("%s", n.Contract)
}

// --------------------------

type ContractNode struct {
	AbstractNode
	Identifier string
	Variables  []VariableNode
}

func (n *ContractNode) String() string {
	return fmt.Sprintf("[%s] CONTRACT %s", n.Pos(), n.Identifier)
}

// --------------------------
// Statement Nodes
// --------------------------

type VariableNode struct {
	StatementNode
	Type       string // todo create TypeNode
	Identifier string
}

func (n *VariableNode) String() string {
	return fmt.Sprintf("[%s] VARIABLE %s %s", n.Pos(), n.Type, n.Identifier)
}

// --------------------------
