package parser

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/pkg/errors"
)

type Parser struct {
	lex          *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []error
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		lex: lex,
	}

	// read two tokens at the beginning
	p.nextToken()
	p.nextTokenWhileNewLine()

	return p
}

func (p *Parser) ParseProgram() (*node.ProgramNode, []error) {
	program := &node.ProgramNode{}

	if p.isSymbol(token.Contract) {
		program.Contract = p.parseContract()
	}

	if !p.isEnd() {
		p.addError("Invalid token outside contract: " + p.currentToken.String())
	}
	return program, p.errors
}

func (p *Parser) parseContract() *node.ContractNode {
	contract := &node.ContractNode{
		AbstractNode: p.newAbstractNode(),
	}
	p.nextToken() // skip contract keyword

	contract.Name = p.readIdentifier()
	p.check(token.OpenBrace)
	p.checkAndSkipNewLines(token.NewLine) // force new line for contract body

	for !p.isEnd() && !p.isSymbol(token.CloseBrace) {
		p.parseContractBody(contract)
	}

	p.check(token.CloseBrace)
	return contract
}

func (p *Parser) parseContractBody(contract *node.ContractNode) {
	switch p.currentToken.Type() {
	case token.IDENTIFER:
		contract.Variables = append(contract.Variables, p.parseVariableStatement())
	case token.SYMBOL:
		ftok, _ := p.currentToken.(*token.FixToken)

		switch ftok.Value {
		case token.Function:
			contract.Functions = append(contract.Functions, p.parseFunction())
		default:
			// TODO Parse all types of fix tokens in a contract
			p.addError(fmt.Sprintf("Unsupported symbol %s in contract", ftok.Lexeme))
			p.nextToken()
		}
	default:
		p.addError("Unsupported contract part starting with" + p.currentToken.Literal())
		p.nextToken()
	}
}

func (p *Parser) parseFunction() *node.FunctionNode {
	p.nextToken() // skip function keyword

	function := &node.FunctionNode{
		AbstractNode: p.newAbstractNode(),
	}

	function.ReturnTypes = p.parseReturnTypes()
	function.Name = p.readIdentifier()
	function.Parameters = p.parseParameters()
	function.Body = p.parseStatementBlock()

	return function
}

func (p *Parser) parseReturnTypes() []*node.TypeNode {
	var returnTypes []*node.TypeNode

	if p.isSymbol(token.OpenParen) {
		p.nextToken() // skip '('
		returnTypes = append(returnTypes, p.parseType())
		for !p.isEnd() && p.isSymbol(token.Comma) {
			p.nextToken() // skip ','
			returnTypes = append(returnTypes, p.parseType())
		}
		p.check(token.CloseParen)
	} else {
		returnTypes = append(returnTypes, p.parseType())
	}

	return returnTypes
}

func (p *Parser) parseParameters() []*node.VariableNode {
	var parameters []*node.VariableNode

	p.check(token.OpenParen)
	isFirstParam := true
	for !p.isEnd() && !p.isSymbol(token.CloseParen) {
		if !isFirstParam {
			p.checkAndSkipNewLines(token.Comma)
		}
		parameters = append(parameters, p.parseVariable())
		isFirstParam = false
	}
	p.check(token.CloseParen)
	return parameters
}

// Statements
// -------------------------

func (p *Parser) parseStatementBlock() []node.StatementNode {
	p.check(token.OpenBrace)
	p.checkAndSkipNewLines(token.NewLine)

	var statements []node.StatementNode
	for !p.isEnd() && !p.isSymbol(token.CloseBrace) {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	p.checkAndSkipNewLines(token.CloseBrace)
	return statements
}

func (p *Parser) parseStatement() node.StatementNode {
	switch p.currentToken.Type() {
	case token.IDENTIFER:
		return p.parseStatementWithIdentifier()
	case token.SYMBOL:
		return p.parseStatementWithFixToken()
	default:
		p.addError("Unsupported statement starting with" + p.currentToken.Literal())
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseStatementWithIdentifier() node.StatementNode {
	if p.peekToken.Type() == token.IDENTIFER {
		return p.parseVariableStatement()
	}

	identifier := p.readIdentifier()
	if p.isSymbol(token.Assign) {
		return p.parseAssignmentStatement(identifier)
	}

	p.addError("%s not yet implemented" + p.currentToken.Literal())
	p.nextToken()
	return nil
}

func (p *Parser) parseStatementWithFixToken() node.StatementNode {
	ftok, _ := p.currentToken.(*token.FixToken)

	switch ftok.Value {
	case token.If:
		return p.parseIfStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		p.addError("Unsupported statement starting with" + ftok.Literal())
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseVariableStatement() *node.VariableNode {
	v := p.parseVariable()
	p.checkAndSkipNewLines(token.NewLine)
	return v
}

func (p *Parser) parseAssignmentStatement(identifier string) node.StatementNode {
	abstractNode := p.newAbstractNode()

	designator := &node.DesignatorNode{
		AbstractNode: abstractNode,
		Value:     	  identifier,
	}

	p.nextToken() // skip '=' sign

	expression := p.parseExpression()


	p.checkAndSkipNewLines(token.NewLine)
	return &node.AssignmentStatementNode{
		AbstractNode:	abstractNode,
		Left:			designator,
		Right:			expression,
	}
}

func (p *Parser) parseIfStatement() *node.IfStatementNode {
	abstractNode := p.newAbstractNode()

	p.nextToken() // skip 'if' keyword

	// Condition
	p.checkAndSkipNewLines(token.OpenParen)
	condition := p.parseExpression()
	p.checkAndSkipNewLines(token.CloseParen)

	// Then
	p.checkAndSkipNewLines(token.OpenBrace)
	then := &node.StatementBlockNode{
		AbstractNode: p.newAbstractNode(),
		Statements: p.parseStatementBlock(),
	}
	p.check(token.CloseBrace)

	alternative := &node.StatementBlockNode{}

	if p.isSymbol(token.Else) {
		p.nextToken() // skip 'else' keyword

		// Else
		p.checkAndSkipNewLines(token.OpenBrace)
		alternative.AbstractNode = p.newAbstractNode()
		alternative.Statements = p.parseStatementBlock()
		p.checkAndSkipNewLines(token.CloseBrace)
	}

	return &node.IfStatementNode{
		AbstractNode: abstractNode,
		Condition: condition,
		Then: then,
		Else: alternative,
	}
}

func (p *Parser) parseReturnStatement() *node.ReturnStatementNode {
	var returnValues []node.ExpressionNode
	abstractNode := p.newAbstractNode()

	p.nextToken() // skip 'return' keyword

	returnValues = append(returnValues, p.parseExpression())

	for !p.isEnd() && p.isSymbol(token.Comma){
		p.nextToken() // skip ','
		returnValues = append(returnValues, p.parseExpression())
	}

	return &node.ReturnStatementNode{
		AbstractNode: abstractNode,
		Expression: returnValues,
	}
}

func (p *Parser) parseVariable() *node.VariableNode {
	v := &node.VariableNode{
		AbstractNode: p.newAbstractNode(),
		Type:         p.parseType(),
		Identifier:   p.readIdentifier(),
	}
	if p.isSymbol(token.Assign) {
		p.nextToken()
		v.Expression = p.parseExpression()
	}

	return v
}

func (p *Parser) parseType() *node.TypeNode {
	// Later we need to distinguish between an array and a simple type
	return &node.TypeNode{
		AbstractNode: p.newAbstractNode(),
		Identifier:   p.readIdentifier(),
	}
}

// Expressions
// -------------------------

func (p *Parser) parseExpression() node.ExpressionNode {
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
	if p.isAnySymbol(token.Addition, token.Subtraction) {
		return p.parseUnaryExpression()
	}

	return p.parseOperand()
}

func (p *Parser) parseUnaryExpression() *node.UnaryExpression {
	return &node.UnaryExpression{
		AbstractNode: p.newAbstractNode(),
		Operator:     p.readSymbol(),
		Operand:      p.parseFactor(),
	}
}

func (p *Parser) parseOperand() node.ExpressionNode {
	switch p.currentToken.Type() {
	case token.IDENTIFER:
		return p.parseDesignator()
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
		error = "Invalid operand: " + p.currentToken.Literal()
	}

	return p.newErrorNode(error)
}

func (p *Parser) parseDesignator() *node.DesignatorNode {
	return &node.DesignatorNode{
		AbstractNode: p.newAbstractNode(),
		Value:        p.readIdentifier(),
	}
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

// Helper functions
// -----------------

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) nextTokenWhileNewLine() {
	p.nextToken()
	for p.isSymbol(token.NewLine) {
		p.nextToken()
	}
}

func (p *Parser) isType(expectedType token.TokenType) bool {
	return p.currentToken.Type() == expectedType
}

func (p *Parser) isSymbol(symbol token.Symbol) bool {
	tok, ok := p.currentToken.(*token.FixToken)
	return ok && tok.Value == symbol
}

func (p *Parser) isAnySymbol(expectedSymbols ...token.Symbol) bool {
	if tok, ok := p.currentToken.(*token.FixToken); ok {
		for _, s := range expectedSymbols {
			if tok.Value == s {
				return true
			}
		}
	}
	return false
}

func (p *Parser) check(symbol token.Symbol) {
	if !p.isSymbol(symbol) {
		p.addError(fmt.Sprintf("Symbol %s expected, but got %s", token.SymbolLexeme[symbol], p.currentToken.Literal()))
	}
	p.nextToken()
}

func (p *Parser) checkAndSkipNewLines(symbol token.Symbol) {
	p.check(symbol)
	p.skipNewLines()
}

func (p *Parser) skipNewLines() {
	for p.isSymbol(token.NewLine) {
		p.nextToken()
	}
}

func (p *Parser) readIdentifier() string {
	var identifier string

	if tok, ok := p.currentToken.(*token.IdentifierToken); ok {
		identifier = tok.Literal()
	} else {
		p.addError("Identifier expected")
		identifier = "ERROR"
	}

	p.nextToken()
	return identifier
}

func (p *Parser) readSymbol() token.Symbol {
	tok, ok := p.currentToken.(*token.FixToken)
	if ok {
		p.nextToken()
		return tok.Value
	}
	panic("Invalid operation")
}

func (p *Parser) newAbstractNode() node.AbstractNode {
	return node.AbstractNode{
		Position: p.currentToken.Pos(),
	}
}

func (p *Parser) newErrorNode(msg string) *node.ErrorNode {
	p.addError(msg)
	e := &node.ErrorNode{
		AbstractNode: p.newAbstractNode(),
		Message:      msg,
	}

	p.nextToken()
	return e
}

func (p *Parser) isEnd() bool {
	return p.lex.IsEnd
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors,
		errors.New(fmt.Sprintf("[%s] ERROR: %s", p.currentToken.Pos().String(), msg)))
}
