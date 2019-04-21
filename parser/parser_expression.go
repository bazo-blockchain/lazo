package parser

import (
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
)

// Expressions
// -------------------------

func (p *Parser) parseExpression() node.ExpressionNode {
	return p.parseOr()
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
	leftExpr := p.parseEquality()

	for p.isAnySymbol(token.And) {
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
	leftExpr := p.parseTerm()

	for p.isAnySymbol(token.Less, token.LessEqual, token.GreaterEqual, token.Greater) {
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

	for p.isAnySymbol(token.Addition, token.Subtraction) {
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
	if p.isAnySymbol(token.Addition, token.Subtraction, token.Not) {
		return p.parseUnaryExpression()
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
		return p.parseBoolean()
	}

	var error string
	if tok, ok := p.currentToken.(*token.ErrorToken); ok {
		error = tok.Msg
	} else {
		panic("Unsupported token type: " + p.currentToken.Literal())
	}

	return p.newErrorNode(error)
}

func (p *Parser) parseDesignator() node.DesignatorNode {
	abstractNode := p.newAbstractNode()
	var left node.DesignatorNode = &node.BasicDesignatorNode{
		AbstractNode: abstractNode,
		Value:        p.readIdentifier(),
	}

	for p.isSymbol(token.Period) || p.isSymbol(token.OpenBracket) {
		if p.isSymbol(token.Period) {
			p.nextToken()
			identifier := p.readIdentifier()
			left = &node.MemberAccessNode{
				AbstractNode: abstractNode,
				Designator:   left,
				Identifier:   identifier,
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

func (p *Parser) parseBoolean() node.ExpressionNode {
	tok, _ := p.currentToken.(*token.FixToken)

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
