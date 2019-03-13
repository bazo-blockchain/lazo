package visitor

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type AbstractVisitor struct {
	concreteVisitor node.Visitor
}

func (v *AbstractVisitor) VisitProgramNode(node *node.ProgramNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitContractNode(node *node.ContractNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitFunctionNode(node *node.FunctionNode) {
	for _, statement := range node.Body {
		fmt.Println("base " + node.Name)
		statement.Accept(v.concreteVisitor)
	}
}

func (v *AbstractVisitor) VisitVariableNode(node *node.VariableNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitTypeNode(node *node.TypeNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitBinaryExpressionNode(node *node.BinaryExpressionNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitUnaryExpressionNode(node *node.UnaryExpression) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitDesignatorNode(node *node.DesignatorNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitStringLiteralNode(node *node.StringLiteralNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitCharacterLiteralNode(node *node.CharacterLiteralNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitBoolLiteralNode(node *node.BoolLiteralNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitErrorNode(node *node.ErrorNode) {
	// TODO Implement
}


