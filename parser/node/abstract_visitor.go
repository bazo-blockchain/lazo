package node

import (
	"fmt"
)

type AbstractVisitor struct {
	ConcreteVisitor Visitor
}

func (v *AbstractVisitor) VisitProgramNode(node *ProgramNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitContractNode(node *ContractNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitFunctionNode(node *FunctionNode) {
	for _, statement := range node.Body {
		fmt.Println("base " + node.Name)
		statement.Accept(v.ConcreteVisitor)
	}
}

func (v *AbstractVisitor) VisitVariableNode(node *VariableNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitTypeNode(node *TypeNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitIfStatementNode(node *IfStatementNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitReturnStatementNode(node *ReturnStatementNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitAssignmentStatementNode(node *AssignmentStatementNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitBinaryExpressionNode(node *BinaryExpressionNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitUnaryExpressionNode(node *UnaryExpression) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitDesignatorNode(node *DesignatorNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitIntegerLiteralNode(node *IntegerLiteralNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitStringLiteralNode(node *StringLiteralNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitCharacterLiteralNode(node *CharacterLiteralNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitBoolLiteralNode(node *BoolLiteralNode) {
	// TODO Implement
}
func (v *AbstractVisitor) VisitErrorNode(node *ErrorNode) {
	// TODO Implement
}


