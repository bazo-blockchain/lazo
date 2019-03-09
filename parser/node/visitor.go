package node

type Visitor interface {
	VisitProgramNode(node *ProgramNode)
	VisitContractNode(node *ContractNode)
	VisitFunctionNode(node *FunctionNode)
	VisitVariableNode(node *VariableNode)
	VisitTypeNode(node *TypeNode)
	VisitIfStatementNode(node *IfStatementNode)
	VisitReturnStatementNode(node *ReturnStatementNode)
	VisitAssignmentStatementNode(node *AssignmentStatementNode)
	VisitBinaryExpressionNode(node *BinaryExpressionNode)
	VisitUnaryExpressionNode(node *UnaryExpression)
	VisitDesignatorNode(node *DesignatorNode)
	VisitIntegerLiteralNode(node *IntegerLiteralNode)
	VisitStringLiteralNode(node *StringLiteralNode)
	VisitCharacterLiteralNode(node *CharacterLiteralNode)
	VisitBoolLiteralNode(node *BoolLiteralNode)
	VisitErrorNode(node *ErrorNode)
}
