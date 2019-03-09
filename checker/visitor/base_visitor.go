package visitor

import "github.com/bazo-blockchain/lazo/parser/node"

type BaseVisitor struct {}

func (v *BaseVisitor) VisitProgramNode(node *node.ProgramNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitContractNode(node *node.ContractNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitFunctionNode(node *node.FunctionNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitVariableNode(node *node.VariableNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitTypeNode(node *node.TypeNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitIfStatementNode(node *node.IfStatementNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitBinaryExpressionNode(node *node.BinaryExpressionNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitUnaryExpressionNode(node *node.UnaryExpression) {
	// TODO Implement
}
func (v *BaseVisitor) VisitDesignatorNode(node *node.DesignatorNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitStringLiteralNode(node *node.StringLiteralNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitCharacterLiteralNode(node *node.CharacterLiteralNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitBoolLiteralNode(node *node.BoolLiteralNode) {
	// TODO Implement
}
func (v *BaseVisitor) VisitErrorNode(node *node.ErrorNode) {
	// TODO Implement
}


