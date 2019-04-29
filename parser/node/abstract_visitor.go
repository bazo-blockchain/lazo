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
	for _, variable := range node.Fields {
		variable.Accept(v.ConcreteVisitor)
	}

	if node.Constructor != nil {
		node.Constructor.Accept(v.ConcreteVisitor)
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

// VisitStructNode traverses the struct fields
func (v *AbstractVisitor) VisitStructNode(node *StructNode) {
	for _, field := range node.Fields {
		field.Accept(v.ConcreteVisitor)
	}
}

// VisitArrayNode traverses the type node
func (v *AbstractVisitor) VisitArrayNode(node *ArrayNode) {
	node.Type.Accept(v.ConcreteVisitor)
}

// VisitStructFieldNode traverses the type node
func (v *AbstractVisitor) VisitStructFieldNode(node *StructFieldNode) {
	node.Type.Accept(v.ConcreteVisitor)
}

// VisitConstructorNode traverses the parameters and the statement block
func (v *AbstractVisitor) VisitConstructorNode(node *ConstructorNode) {
	for _, paramType := range node.Parameters {
		paramType.Accept(v.ConcreteVisitor)
	}
	v.ConcreteVisitor.VisitStatementBlock(node.Body)
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

// VisitParameterNode traverses the type node
func (v *AbstractVisitor) VisitParameterNode(node *ParameterNode) {
	node.Type.Accept(v.ConcreteVisitor)
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

// VisitMultiVariableNode traverses multiple variables declarations and
func (v *AbstractVisitor) VisitMultiVariableNode(node *MultiVariableNode) {
	for _, t := range node.Types {
		t.Accept(v.ConcreteVisitor)
	}
	node.FuncCall.Accept(v.ConcreteVisitor)
}

// VisitBasicTypeNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitBasicTypeNode(node *BasicTypeNode) {
	// Nothing to do here
}

// VisitArrayTypeNode visits the type node
func (v *AbstractVisitor) VisitArrayTypeNode(node *ArrayTypeNode) {
	node.ElementType.Accept(v.ConcreteVisitor)
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

// VisitMultiAssignmentStatementNode traverses the target designators and the function call
func (v *AbstractVisitor) VisitMultiAssignmentStatementNode(node *MultiAssignmentStatementNode) {
	for _, designator := range node.Designators {
		designator.Accept(v.ConcreteVisitor)
	}
	node.FuncCall.Accept(v.ConcreteVisitor)
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

// VisitStructCreationNode traverses the field initialization expressions
func (v *AbstractVisitor) VisitStructCreationNode(node *StructCreationNode) {
	for _, expr := range node.FieldValues {
		expr.Accept(v.ConcreteVisitor)
	}
}

// VisitStructNamedCreationNode traverses the named field initialization expressions
func (v *AbstractVisitor) VisitStructNamedCreationNode(node *StructNamedCreationNode) {
	for _, namedField := range node.FieldValues {
		namedField.Accept(v.ConcreteVisitor)
	}
}

// VisitStructFieldAssignmentNode traverse the field initialization expression
func (v *AbstractVisitor) VisitStructFieldAssignmentNode(node *StructFieldAssignmentNode) {
	node.Expression.Accept(v.ConcreteVisitor)
}

// VisitArrayCreationNode traverses the element expression nodes
func (v *AbstractVisitor) VisitArrayCreationNode(node *ArrayCreationNode) {
	for _, elementValue := range node.ElementValues {
		elementValue.Accept(v.ConcreteVisitor)
	}
}

// VisitBasicDesignatorNode does nothing because it is the terminal node.
func (v *AbstractVisitor) VisitBasicDesignatorNode(node *BasicDesignatorNode) {
	// Nothing to do here
}

// VisitElementAccessNode visits the designator and the expression.
func (v *AbstractVisitor) VisitElementAccessNode(node *ElementAccessNode) {
	node.Designator.Accept(v.ConcreteVisitor)
	node.Expression.Accept(v.ConcreteVisitor)
}

// VisitMemberAccessNode visit the designator.
func (v *AbstractVisitor) VisitMemberAccessNode(node *MemberAccessNode) {
	node.Designator.Accept(v.ConcreteVisitor)
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
