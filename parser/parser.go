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
	p.checkAndSkipNewLines(token.NewLine)

	for !p.isEnd() && !p.isSymbol(token.CloseBrace) {
		p.parseContractBody(contract)
	}

	p.check(token.CloseBrace)
	return contract
}

func (p *Parser) parseContractBody(contract *node.ContractNode) {
	switch p.currentToken.Type() {
	case token.IDENTIFER:
		contract.Variables = append(contract.Variables, p.parseVariable())
	case token.SYMBOL:
		ftok, _ :=  p.currentToken.(*token.FixToken)

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

func (p *Parser) parseExpression() node.ExpressionNode {
	// TODO implement
	return nil
}

func (p *Parser) parseFunction() *node.FunctionNode {
	// skip function keyword
	p.check(token.Function)

	function := &node.FunctionNode{
		AbstractNode: p.newAbstractNode(),
		ReturnTypes:  []*node.TypeNode{},
		Parameters:   []*node.VariableNode{},
		Body:         []node.StatementNode{},
	}

	function.ReturnTypes = p.parseReturnTypes()

	function.Identifier = p.readIdentifier()

	function.Parameters = p.parseParameters()

	function.Body = p.parseFunctionBody()

	return function
}

// TODO Refactor: Move for loop in parseParameters and parseReturnTypes to own function
func (p *Parser) parseParameters() []*node.VariableNode {
	var parameters []*node.VariableNode

	p.check(token.OpenParen)
	for !p.isEnd() && !p.isSymbol(token.CloseParen) {
		parameters = append(parameters, p.parseVariable())
		p.nextToken()
		if p.isSymbol(token.Comma) {
			p.nextToken()
		} else if p.isSymbol(token.CloseParen) {
			continue
		} else {
			// error
		}
	}
	p.check(token.CloseParen)
	return parameters
}

func (p *Parser) parseFunctionBody() []node.StatementNode {
	p.check(token.OpenBrace)
	// TODO Implement
	p.check(token.CloseBrace)
	return nil
}

func (p *Parser) parseStatement() node.StatementNode {

	if p.isSymbol(token.If) {
		return p.parseIfStatement()
	} else if p.isSymbol(token.Return) {
		return p.parseReturnStatement()
	} else if p.isType(token.IDENTIFER) {
		identifier := p.readIdentifier()
		return p.parseStatementWithIdentifier(identifier)
	} else {
		// error
		return nil
	}
}

func (p *Parser) parseIfStatement() *node.IfStatementNode {
	return nil
}

func (p *Parser) parseReturnStatement() *node.ReturnStatementNode {
	return nil
}

func (p *Parser) parseStatementWithIdentifier(identifier string) node.StatementNode {
	return nil
}

func (p *Parser) parseReturnTypes() []*node.TypeNode {
	var returnTypes []*node.TypeNode

	if p.isSymbol(token.OpenParen) {
		p.nextToken()
		for !p.isEnd() && !p.isSymbol(token.CloseParen) {
			returnTypes = append(returnTypes, p.parseType())
			p.nextToken()
			if p.isSymbol(token.Comma) {
				p.nextToken()
			} else if p.isSymbol(token.CloseParen) {
				continue
			} else {
				// error
			}
		}
		p.check(token.CloseParen)
	} else {
		returnTypes = append(returnTypes, p.parseType())
	}

	return returnTypes
}

func (p *Parser) parseVariable() *node.VariableNode {
	// TODO Implement Assignment
	v := &node.VariableNode{
		AbstractNode: p.newAbstractNode(),
		Type:         p.parseType(),
		Identifier:   p.readIdentifier(),
	}
	p.checkAndSkipNewLines(token.NewLine)
	return v
}

func (p *Parser) parseType() *node.TypeNode {
	// Later we need to distinguish between an array and a simple type
	return &node.TypeNode{
		AbstractNode: p.newAbstractNode(),
		Identifier:   p.readIdentifier(),
	}
}

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

// helpers
// --------------

func (p *Parser) isType(expectedType token.TokenType) bool {
	return p.currentToken.Type() == expectedType
}

func (p *Parser) isSymbol(symbol token.Symbol) bool {
	tok, ok := p.currentToken.(*token.FixToken)
	return ok && tok.Value == symbol
}

func (p *Parser) check(symbol token.Symbol) {
	if !p.isSymbol(symbol) {
		p.addError(fmt.Sprintf("FixToken symbol[%d] expected, but got %s", symbol, p.currentToken.Literal()))
	}
	p.nextToken()
}

func (p *Parser) checkAndSkipNewLines(symbol token.Symbol) {
	p.check(symbol)
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

func (p *Parser) newAbstractNode() node.AbstractNode {
	return node.AbstractNode{
		Position: p.currentToken.Pos(),
	}
}

func (p *Parser) isEnd() bool {
	return p.lex.IsEnd
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors,
		errors.New(fmt.Sprintf("[%s] ERROR: %s", p.currentToken.Pos().String(), msg)))
}
