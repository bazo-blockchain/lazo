package node

// Visitor is the interface that defines the functions a concrete visitor should implement.
type Visitor interface {
	VisitProgramNode(node *ProgramNode)
	VisitContractNode(node *ContractNode)
	VisitFieldNode(node *FieldNode)
	VisitFunctionNode(node *FunctionNode)
	VisitParameterNode(node *ParameterNode)
	VisitStatementBlock(stmts []StatementNode)
	VisitVariableNode(node *VariableNode)
	VisitMultiVariableNode(node *MultiVariableNode)
	VisitTypeNode(node *TypeNode)
	VisitIfStatementNode(node *IfStatementNode)
	VisitReturnStatementNode(node *ReturnStatementNode)
	VisitAssignmentStatementNode(node *AssignmentStatementNode)
	VisitMultiAssignmentStatementNode(node *MultiAssignmentStatementNode)
	VisitCallStatementNode(node *CallStatementNode)
	VisitBinaryExpressionNode(node *BinaryExpressionNode)
	VisitUnaryExpressionNode(node *UnaryExpressionNode)
	VisitDesignatorNode(node *DesignatorNode)
	VisitFuncCallNode(node *FuncCallNode)
	VisitIntegerLiteralNode(node *IntegerLiteralNode)
	VisitStringLiteralNode(node *StringLiteralNode)
	VisitCharacterLiteralNode(node *CharacterLiteralNode)
	VisitBoolLiteralNode(node *BoolLiteralNode)
	VisitErrorNode(node *ErrorNode)
}
