package node

type AbstractVisitor struct {
	ConcreteVisitor Visitor
}

func (v *AbstractVisitor) VisitProgramNode(node *ProgramNode) {
	node.Contract.Accept(v.ConcreteVisitor)
}

func (v *AbstractVisitor) VisitContractNode(node *ContractNode) {
	for _, variable := range node.Variables {
		variable.Accept(v.ConcreteVisitor)
	}
	for _, function := range node.Functions {
		function.Accept(v.ConcreteVisitor)
	}
}

func (v *AbstractVisitor) VisitFunctionNode(node *FunctionNode) {
	for _, statement := range node.Body {
		statement.Accept(v.ConcreteVisitor)
	}
}

func (v *AbstractVisitor) VisitVariableNode(node *VariableNode) {
	node.Type.Accept(v.ConcreteVisitor)
	if node.Expression != nil {
		node.Expression.Accept(v.ConcreteVisitor)
	}
}

func (v *AbstractVisitor) VisitTypeNode(node *TypeNode) {
	// TODO Implement
}

func (v *AbstractVisitor) VisitIfStatementNode(node *IfStatementNode) {
	node.Condition.Accept(v.ConcreteVisitor)
	for _, statement := range node.Then {
		statement.Accept(v.ConcreteVisitor)
	}
	for _, statement := range node.Else {
		statement.Accept(v.ConcreteVisitor)
	}
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


