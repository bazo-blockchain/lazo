package visitor

import "github.com/bazo-blockchain/lazo/parser/node"

type Visitor interface {
	VisitProgramNode(node *node.ProgramNode)
	VisitContractNode(node *node.ContractNode)
	VisitFunctionNode(node *node.FunctionNode)
	VisitVariableNode(node *node.VariableNode)
	VisitTypeNode(node *node.TypeNode)
	VisitIfStatementNode(node *node.IfStatementNode)
	VisitReturnStatementNode(node *node.ReturnStatementNode)
	VisitAssignmentStatementNode(node *node.AssignmentStatementNode)
	VisitBinaryExpressionNode(node *node.BinaryExpressionNode)
	VisitUnaryExpressionNode(node *node.UnaryExpression)
	VisitDesignatorNode(node *node.DesignatorNode)
	VisitIntegerLiteralNode(node *node.IntegerLiteralNode)
	VisitStringLiteralNode(node *node.StringLiteralNode)
	VisitCharacterLiteralNode(node *node.CharacterLiteralNode)
	VisitBoolLiteralNode(node *node.BoolLiteralNode)
	VisitErrorNode(node *node.ErrorNode)
}

type BaseVisitor struct {}

func (v *BaseVisitor) VisitProgramNode(node *node.ProgramNode) {}
func (v *BaseVisitor) VisitContractNode(node *node.ContractNode) {}
func (v *BaseVisitor) VisitFunctionNode(node *node.FunctionNode) {}
func (v *BaseVisitor) VisitVariableNode(node *node.VariableNode) {}
func (v *BaseVisitor) VisitTypeNode(node *node.TypeNode) {}
func (v *BaseVisitor) VisitIfStatementNode(node *node.IfStatementNode) {}
func (v *BaseVisitor) VisitReturnStatementNode(node *node.ReturnStatementNode) {}
func (v *BaseVisitor) VisitAssignmentStatementNode(node *node.AssignmentStatementNode) {}
func (v *BaseVisitor) VisitBinaryExpressionNode(node *node.BinaryExpressionNode) {}
func (v *BaseVisitor) VisitUnaryExpressionNode(node *node.UnaryExpression) {}
func (v *BaseVisitor) VisitDesignatorNode(node *node.DesignatorNode) {}
func (v *BaseVisitor) VisitIntegerLiteralNode(node *node.IntegerLiteralNode) {}
func (v *BaseVisitor) VisitStringLiteralNode(node *node.StringLiteralNode) {}
func (v *BaseVisitor) VisitCharacterLiteralNode(node *node.CharacterLiteralNode) {}
func (v *BaseVisitor) VisitBoolLiteralNode(node *node.BoolLiteralNode) {}
func (v *BaseVisitor) VisitErrorNode(node *node.ErrorNode) {}


