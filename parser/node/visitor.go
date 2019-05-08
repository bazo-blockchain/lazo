package node

// Visitor is the interface that defines the functions a concrete visitor should implement.
type Visitor interface {
	VisitProgramNode(node *ProgramNode)
	VisitContractNode(node *ContractNode)
	VisitFieldNode(node *FieldNode)
	VisitStructNode(node *StructNode)
	VisitStructFieldNode(node *StructFieldNode)
	VisitConstructorNode(node *ConstructorNode)
	VisitFunctionNode(node *FunctionNode)
	VisitParameterNode(node *ParameterNode)
	VisitStatementBlock(stmts []StatementNode)
	VisitVariableNode(node *VariableNode)
	VisitMultiVariableNode(node *MultiVariableNode)
	VisitBasicTypeNode(node *BasicTypeNode)
	VisitArrayTypeNode(node *ArrayTypeNode)
	VisitMapTypeNode(node *MapTypeNode)
	VisitIfStatementNode(node *IfStatementNode)
	VisitReturnStatementNode(node *ReturnStatementNode)
	VisitAssignmentStatementNode(node *AssignmentStatementNode)
	VisitMultiAssignmentStatementNode(node *MultiAssignmentStatementNode)
	VisitShorthandAssignmentNode(node *ShorthandAssignmentStatementNode)
	VisitCallStatementNode(node *CallStatementNode)
	VisitDeleteStatementNode(node *DeleteStatementNode)
	VisitBinaryExpressionNode(node *BinaryExpressionNode)
	VisitUnaryExpressionNode(node *UnaryExpressionNode)
	VisitBasicDesignatorNode(node *BasicDesignatorNode)
	VisitElementAccessNode(node *ElementAccessNode)
	VisitMemberAccessNode(node *MemberAccessNode)
	VisitArrayInitializationNode(node *ArrayInitializationNode)
	VisitFuncCallNode(node *FuncCallNode)
	VisitStructCreationNode(node *StructCreationNode)
	VisitStructNamedCreationNode(node *StructNamedCreationNode)
	VisitStructFieldAssignmentNode(node *StructFieldAssignmentNode)
	VisitArrayLengthCreationNode(node *ArrayLengthCreationNode)
	VisitArrayValueCreationNode(node *ArrayValueCreationNode)
	VisitIntegerLiteralNode(node *IntegerLiteralNode)
	VisitStringLiteralNode(node *StringLiteralNode)
	VisitCharacterLiteralNode(node *CharacterLiteralNode)
	VisitBoolLiteralNode(node *BoolLiteralNode)
	VisitErrorNode(node *ErrorNode)
}
