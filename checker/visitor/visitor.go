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
