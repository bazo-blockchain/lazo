package parser

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
)

// Expressions
// -------------------------

func (p *Parser) parseExpression() node.ExpressionNode {
	return p.parseTernaryExpression()
}

func (p *Parser) parseTernaryExpression() node.ExpressionNode {
	expr := p.parseOr()

	if p.isSymbol(token.QuestionMark) {
		p.nextToken()

		ternary := &node.TernaryExpressionNode{
			AbstractNode: p.newAbstractNodeWithPos(expr.Pos()),
			Condition:    expr,
			Then:         p.parseOr(),
		}
		p.check(token.Colon)
		ternary.Else = p.parseOr()
		return ternary
	}
	return expr
}

func (p *Parser) parseOr() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseAnd()

	for p.isAnySymbol(token.Or) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseAnd(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseAnd() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseBitwiseOr()

	for p.isAnySymbol(token.And) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseBitwiseOr(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseBitwiseOr() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseBitwiseXOr()

	for p.isAnySymbol(token.BitwiseOr) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseBitwiseXOr(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseBitwiseXOr() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseBitwiseAnd()

	for p.isAnySymbol(token.BitwiseXOr) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseBitwiseAnd(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseBitwiseAnd() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseEquality()

	for p.isAnySymbol(token.BitwiseAnd) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseEquality(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseEquality() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseRelationalComparison()

	for p.isAnySymbol(token.Equal, token.Unequal) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseRelationalComparison(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseRelationalComparison() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseBitwiseShift()

	for p.isAnySymbol(token.Less, token.LessEqual, token.GreaterEqual, token.Greater) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseBitwiseShift(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseBitwiseShift() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseTerm()

	for p.isAnySymbol(token.ShiftLeft, token.ShiftRight) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseTerm(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseTerm() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseFactor()

	for p.isAnySymbol(token.Plus, token.Minus) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseFactor(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseFactor() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseExponent()

	for p.isAnySymbol(token.Multiplication, token.Division, token.Modulo) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseExponent(),
		}
		leftExpr = binExpr
	}
	return leftExpr
}

func (p *Parser) parseExponent() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	leftExpr := p.parseExpressionRest()

	if p.isSymbol(token.Exponent) {
		binExpr := &node.BinaryExpressionNode{
			AbstractNode: abstractNode,
			Left:         leftExpr,
			Operator:     p.readSymbol(),
			Right:        p.parseExponent(), // recursive because of right-to-left associativity
		}
		return binExpr
	}
	return leftExpr
}

func (p *Parser) parseExpressionRest() node.ExpressionNode {
	if p.isAnySymbol(token.Plus, token.Minus, token.Not, token.BitwiseNot) {
		return p.parseUnaryExpression()
	}

	if p.isSymbol(token.OpenParen) {
		abstractNode := p.newAbstractNode()
		p.nextToken()
		expr := p.parseExpression()
		p.check(token.CloseParen)

		switch p.currentToken.Type() {
		case token.IDENTIFER, token.CHARACTER, token.INTEGER:
			// (String) x.y.z, (String) 'c', (String) 5
			return p.parseTypeCast(abstractNode, expr)
		case token.SYMBOL:
			// (String) true
			if p.isAnySymbol(token.True, token.False) {
				return p.parseTypeCast(abstractNode, expr)
			}
		}
		return expr
	}

	return p.parseOperand()
}

func (p *Parser) parseUnaryExpression() *node.UnaryExpressionNode {
	return &node.UnaryExpressionNode{
		AbstractNode: p.newAbstractNode(),
		Operator:     p.readSymbol(),
		Expression:   p.parseFactor(),
	}
}

func (p *Parser) parseTypeCast(abstractNode node.AbstractNode, expr node.ExpressionNode) *node.TypeCastNode {
	typeCast := &node.TypeCastNode{
		AbstractNode: abstractNode,
		Expression:   p.parseOperand(),
	}

	if basicDesignator, ok := expr.(*node.BasicDesignatorNode); ok {
		typeCast.Type = &node.BasicTypeNode{
			AbstractNode: p.newAbstractNodeWithPos(expr.Pos()),
			Identifier:   basicDesignator.Value,
		}
	} else {
		p.addError(fmt.Sprintf("Invalid type %s", expr))
	}
	return typeCast
}

func (p *Parser) parseOperand() node.ExpressionNode {
	switch p.currentToken.Type() {
	case token.IDENTIFER:
		designator := p.parseDesignator()
		if p.isSymbol(token.OpenParen) {
			return p.parseFuncCall(designator)
		}
		return designator
	case token.INTEGER:
		return p.parseInteger()
	case token.CHARACTER:
		return p.parseCharacter()
	case token.STRING:
		return p.parseString()
	case token.SYMBOL:
		return p.parseOperandSymbol()
	}

	var error string
	if tok, ok := p.currentToken.(*token.ErrorToken); ok {
		error = tok.Msg
	} else {
		panic("Unsupported token type: " + p.currentToken.Literal())
	}

	return p.newErrorNode(error)
}

func (p *Parser) parseOperandSymbol() node.ExpressionNode {
	tok, ok := p.currentToken.(*token.FixToken)

	if !ok {
		panic("Invalid operation")
	}

	switch tok.Value {
	case token.True, token.False:
		return p.parseBoolean(tok)
	case token.New:
		return p.parseCreation()
	default:
		return p.newErrorNode("Unsupported expression symbol " + p.currentToken.Literal())
	}
}

func (p *Parser) parseDesignator() node.DesignatorNode {
	return p.parseDesignatorWithIdentifier(p.newAbstractNode(), p.readIdentifier())
}

func (p *Parser) parseDesignatorWithIdentifier(abstractNode node.AbstractNode, identifier string) node.DesignatorNode {
	var left node.DesignatorNode = &node.BasicDesignatorNode{
		AbstractNode: abstractNode,
		Value:        identifier,
	}

	for p.isSymbol(token.Period) || p.isSymbol(token.OpenBracket) {
		if p.isSymbol(token.Period) {
			p.nextToken()
			memberIdentifier := p.readIdentifier()
			left = &node.MemberAccessNode{
				AbstractNode: abstractNode,
				Designator:   left,
				Identifier:   memberIdentifier,
			}
		} else {
			p.check(token.OpenBracket)
			exp := p.parseExpression()
			p.check(token.CloseBracket)
			left = &node.ElementAccessNode{
				AbstractNode: abstractNode,
				Designator:   left,
				Expression:   exp,
			}
		}
	}
	return left
}

func (p *Parser) parseFuncCall(designator node.DesignatorNode) *node.FuncCallNode {
	funcCall := &node.FuncCallNode{
		AbstractNode: p.newAbstractNodeWithPos(designator.Pos()),
		Designator:   designator,
	}
	p.check(token.OpenParen)

	isFirstArg := true
	for !p.isEnd() && !p.isSymbol(token.CloseParen) {
		if !isFirstArg {
			p.check(token.Comma)
		}
		funcCall.Args = append(funcCall.Args, p.parseExpression())
		isFirstArg = false
	}
	p.check(token.CloseParen)
	return funcCall
}

func (p *Parser) parseCreation() node.ExpressionNode {
	abstractNode := p.newAbstractNode()
	p.nextToken() // skip 'new' keyword

	identifier := p.readIdentifier()
	if p.isSymbol(token.OpenParen) {
		return p.parseStructCreation(abstractNode, identifier)
	} else if p.isSymbol(token.OpenBracket) {
		return p.parseArrayCreation(abstractNode, identifier)
	}

	return p.newErrorNode(fmt.Sprintf("Unsupported creation type with %s", p.currentToken.Literal()))
}

func (p *Parser) parseArrayCreation(abstractNode node.AbstractNode, identifier string) node.ExpressionNode {
	var arrayType node.TypeNode = &node.BasicTypeNode{
		AbstractNode: abstractNode,
		Identifier:   identifier,
	}

	// Initialization using values: new int[][]{{1, 2}, {3, 4}}
	if p.peekIsSymbol(token.CloseBracket) {
		arrayType = p.parseArrayType(arrayType)
		return &node.ArrayValueCreationNode{
			AbstractNode: abstractNode,
			Type:         arrayType,
			Elements:     p.parseArrayInitialization(),
		}
	}

	// Initialization using Length: new int[2][3]
	p.check(token.OpenBracket)
	var expressions []node.ExpressionNode
	expression := p.parseExpression() // Read length expression
	expressions = append(expressions, expression)
	p.check(token.CloseBracket)

	for !p.isEnd() && p.isSymbol(token.OpenBracket) {
		p.nextToken()
		expressions = append(expressions, p.parseExpression())
		p.check(token.CloseBracket)

		arrayType = &node.ArrayTypeNode{
			AbstractNode: p.newAbstractNodeWithPos(arrayType.Pos()),
			ElementType:  arrayType,
		}
	}

	return &node.ArrayLengthCreationNode{
		AbstractNode: abstractNode,
		ElementType:  arrayType,
		Lengths:      expressions,
	}
}

func (p *Parser) parseArrayInitialization() *node.ArrayInitializationNode {
	abstractNode := p.newAbstractNode()
	p.checkAndSkipNewLines(token.OpenBrace) // skip '{'
	if !p.isEnd() && !p.isSymbol(token.OpenBrace) {
		var expressions []node.ExpressionNode
		expressions = append(expressions, p.parseExpression())
		for !p.isEnd() && !p.isSymbol(token.CloseBrace) {
			p.checkAndSkipNewLines(token.Comma)
			expressions = append(expressions, p.parseExpression())
		}
		p.check(token.CloseBrace)

		return &node.ArrayInitializationNode{
			AbstractNode: abstractNode,
			Values:       expressions,
		}
	}

	arrayInitialization := &node.ArrayInitializationNode{
		AbstractNode: abstractNode,
	}

	if !p.isEnd() {
		var expressions []node.ExpressionNode

		expressions = append(expressions, p.parseArrayInitialization())

		for !p.isEnd() && !p.isSymbol(token.CloseBrace) {
			p.checkAndSkipNewLines(token.Comma)
			expressions = append(expressions, p.parseArrayInitialization())
		}

		p.check(token.CloseBrace)

		arrayInitialization.Values = expressions
	} else {
		p.addError("Invalid array initialization")
	}

	return arrayInitialization
}

func (p *Parser) parseStructCreation(abstractNode node.AbstractNode, identifier string) node.ExpressionNode {
	p.nextTokenWhileNewLine() // skip '('

	if ftok, ok := p.peekToken.(*token.FixToken); ok && ftok.Value == token.Assign {
		return p.parseStructNamedCreation(abstractNode, identifier)
	}

	structCreation := &node.StructCreationNode{
		AbstractNode: abstractNode,
		Name:         identifier,
	}

	isFirstArg := true
	for !p.isEnd() && !p.isSymbol(token.CloseParen) {
		if !isFirstArg {
			p.check(token.Comma)
		}
		structCreation.FieldValues = append(structCreation.FieldValues, p.parseExpression())
		isFirstArg = false
	}
	p.check(token.CloseParen)
	return structCreation
}

func (p *Parser) parseStructNamedCreation(abstractNode node.AbstractNode, identifier string) *node.StructNamedCreationNode {
	structCreation := &node.StructNamedCreationNode{
		AbstractNode: abstractNode,
		Name:         identifier,
	}

	isFirstArg := true
	for !p.isEnd() && !p.isSymbol(token.CloseParen) {
		if !isFirstArg {
			p.checkAndSkipNewLines(token.Comma)
		}

		field := &node.StructFieldAssignmentNode{
			AbstractNode: p.newAbstractNode(),
			Name:         p.readIdentifier(),
		}
		p.check(token.Assign)
		field.Expression = p.parseExpression()
		structCreation.FieldValues = append(structCreation.FieldValues, field)

		p.skipNewLines()
		isFirstArg = false
	}

	p.check(token.CloseParen)
	return structCreation
}

func (p *Parser) parseInteger() *node.IntegerLiteralNode {
	tok, _ := p.currentToken.(*token.IntegerToken)

	i := &node.IntegerLiteralNode{
		AbstractNode: p.newAbstractNode(),
		Value:        tok.Value,
	}
	p.nextToken()
	return i
}

func (p *Parser) parseCharacter() *node.CharacterLiteralNode {
	tok, _ := p.currentToken.(*token.CharacterToken)

	c := &node.CharacterLiteralNode{
		AbstractNode: p.newAbstractNode(),
		Value:        tok.Value,
	}
	p.nextToken()
	return c
}

func (p *Parser) parseString() *node.StringLiteralNode {
	tok, _ := p.currentToken.(*token.StringToken)

	s := &node.StringLiteralNode{
		AbstractNode: p.newAbstractNode(),
		Value:        tok.Literal(),
	}
	p.nextToken()
	return s
}

func (p *Parser) parseBoolean(tok *token.FixToken) node.ExpressionNode {
	if value, ok := token.BooleanConstants[tok.Value]; ok {
		b := &node.BoolLiteralNode{
			AbstractNode: p.newAbstractNode(),
			Value:        value,
		}
		p.nextToken()
		return b
	}

	return p.newErrorNode("Invalid boolean value " + tok.Literal())
}
