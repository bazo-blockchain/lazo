package node

// AbstractVisitor holds a concrete visitor
// and implements the basic tree traversal logic using the visitor pattern.
// By doing so, concrete visitors should not re-implement the same traversal logic in their visit functions.
type AbstractVisitor struct {
	ConcreteVisitor Visitor
}

// VisitProgramNode traverses the contract node.
func (v *AbstractVisitor) VisitProgramNode(node *ProgramNode) {
	node.Contract.Accept(v.ConcreteVisitor)
}

// VisitContractNode traverses the variable and function nodes.
func (v *AbstractVisitor) VisitContractNode(node *ContractNode) {
	for _, variable := range node.Variables {
		variable.Accept(v.ConcreteVisitor)
	}
	for _, function := range node.Functions {
		function.Accept(v.ConcreteVisitor)
	}
}

// VisitFieldNode traverses the type node and the expression (if present).
func (v *AbstractVisitor) VisitFieldNode(node *FieldNode) {
	node.Type.Accept(v.ConcreteVisitor)
	if node.Expression != nil {
		node.Expression.Accept(v.ConcreteVisitor)
	}
}

// VisitFunctionNode traverses the return types, parameters and finally the statement block.
func (v *AbstractVisitor) VisitFunctionNode(node *FunctionNode) {
	for _, returnType := range node.ReturnTypes {
		returnType.Accept(v.ConcreteVisitor)
	}
	for _, paramType := range node.Parameters {
		paramType.Accept(v.ConcreteVisitor)
	}
	v.ConcreteVisitor.VisitStatementBlock(node.Body)
}

// VisitStatementBlock traverses the statement node.
func (v *AbstractVisitor) VisitStatementBlock(stmts []StatementNode) {
	for _, statement := range stmts {
		statement.Accept(v.ConcreteVisitor)
	}
}

// VisitVariableNode traverses the type node and the expression (if present).
func (v *AbstractVisitor) VisitVariableNode(node *VariableNode) {
	node.Type.Accept(v.ConcreteVisitor)
	if node.Expression != nil {
		node.Expression.Accept(v.ConcreteVisitor)
	}
}

// VisitTypeNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitTypeNode(node *TypeNode) {
	// Nothing to do here
}

// VisitIfStatementNode traverses the condition, then-block and finally else-block.
func (v *AbstractVisitor) VisitIfStatementNode(node *IfStatementNode) {
	node.Condition.Accept(v.ConcreteVisitor)
	v.ConcreteVisitor.VisitStatementBlock(node.Then)
	v.ConcreteVisitor.VisitStatementBlock(node.Else)
}

// VisitReturnStatementNode traverses all the available expressions.
func (v *AbstractVisitor) VisitReturnStatementNode(node *ReturnStatementNode) {
	for _, expr := range node.Expressions {
		expr.Accept(v.ConcreteVisitor)
	}
}

// VisitAssignmentStatementNode traverses the target designator and the value expression.
func (v *AbstractVisitor) VisitAssignmentStatementNode(node *AssignmentStatementNode) {
	node.Left.Accept(v.ConcreteVisitor)
	node.Right.Accept(v.ConcreteVisitor)
}

// VisitCallStatementNode traverses the function call expression
func (v *AbstractVisitor) VisitCallStatementNode(node *CallStatementNode) {
	node.Call.Accept(v.ConcreteVisitor)
}

// VisitBinaryExpressionNode traverses the left and right expressions.
func (v *AbstractVisitor) VisitBinaryExpressionNode(node *BinaryExpressionNode) {
	node.Left.Accept(v.ConcreteVisitor)
	node.Right.Accept(v.ConcreteVisitor)
}

// VisitUnaryExpressionNode traverses the expression.
func (v *AbstractVisitor) VisitUnaryExpressionNode(node *UnaryExpressionNode) {
	node.Expression.Accept(v.ConcreteVisitor)
}

// VisitFuncCallNode traverses the funcCall expression
func (v *AbstractVisitor) VisitFuncCallNode(node *FuncCallNode) {
	node.Designator.Accept(v.ConcreteVisitor)
	for _, expr := range node.Args {
		expr.Accept(v.ConcreteVisitor)
	}
}

// VisitDesignatorNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitDesignatorNode(node *DesignatorNode) {
	// Nothing to do here
}

// VisitIntegerLiteralNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitIntegerLiteralNode(node *IntegerLiteralNode) {
	// Nothing to do here
}

// VisitStringLiteralNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitStringLiteralNode(node *StringLiteralNode) {
	// Nothing to do here
}

// VisitCharacterLiteralNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitCharacterLiteralNode(node *CharacterLiteralNode) {
	// Nothing to do here
}

// VisitBoolLiteralNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitBoolLiteralNode(node *BoolLiteralNode) {
	// Nothing to do here
}

// VisitErrorNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitErrorNode(node *ErrorNode) {
	// Nothing to do here
}
