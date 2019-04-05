package node

// Visitor is the interface that defines the functions a concrete visitor should implement.
type Visitor interface {
	VisitProgramNode(node *ProgramNode)
	VisitContractNode(node *ContractNode)
	VisitFunctionNode(node *FunctionNode)
	VisitStatementBlock(stmts []StatementNode)
	VisitVariableNode(node *VariableNode)
	VisitTypeNode(node *TypeNode)
	VisitIfStatementNode(node *IfStatementNode)
	VisitReturnStatementNode(node *ReturnStatementNode)
	VisitAssignmentStatementNode(node *AssignmentStatementNode)
	VisitBinaryExpressionNode(node *BinaryExpressionNode)
	VisitUnaryExpressionNode(node *UnaryExpressionNode)
	VisitDesignatorNode(node *DesignatorNode)
	VisitIntegerLiteralNode(node *IntegerLiteralNode)
	VisitStringLiteralNode(node *StringLiteralNode)
	VisitCharacterLiteralNode(node *CharacterLiteralNode)
	VisitBoolLiteralNode(node *BoolLiteralNode)
	VisitErrorNode(node *ErrorNode)
}
