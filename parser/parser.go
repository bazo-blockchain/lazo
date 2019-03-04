package parser

import (
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
)

type Parser struct {
	lex          *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		lex: lex,
	}

	// read two tokens at the beginning
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *node.ProgramNode {
	program := &node.ProgramNode{}

	if p.is(token.Contract) {
		program.Contract = p.parseContract()
	}
	// todo error handling
	return program
}

func (p *Parser) parseContract() *node.ContractNode {
	contract := &node.ContractNode{
		AbstractNode: p.newAbstractNode(),
	}
	p.nextToken() // skip contract keyword

	contract.Identifier = p.readIdentifier()
	p.check(token.OpenBrace)

	// parse variables (later: extract to parseContractBody method with other types of nodes)
	contract.Variables = []*node.VariableNode{}
	for !p.isEnd() && !p.is(token.CloseBrace) {
		switch p.currentToken.(type) {
		case *token.IdentifierToken:
			// TODO Implement Assignment
			contract.Variables = append(contract.Variables, p.parseVariable())
		case *token.FixToken:
			// TODO Parse all types of fix tokens in a contract
			if p.is(token.Function) {
				contract.Functions = append(contract.Functions, p.parseFunction())
			}

		default:
			// error
		}
	}

	p.check(token.CloseBrace)
	return contract
}

func (p *Parser) parseExpression() *node.ExpressionNode {
	// TODO implement
	return nil
}

func (p *Parser) parseStatement() *node.StatementNode {
	// TODO implement
	return nil
}

func (p *Parser) parseFunction() *node.FunctionNode{
	// skip function keyword
	p.check(token.Function)

	function := &node.FunctionNode{
		AbstractNode: p.newAbstractNode(),
		ReturnTypes: []*node.TypeNode{},
		Parameters: []*node.VariableNode{},
		Body: []*node.StatementNode{},
	}

	function.ReturnTypes = p.parseReturnTypes()

	function.Identifier = p.readIdentifier()

	function.Parameters = p.parseParameters()

	function.Body = p.parseFunctionBody()

	return function
}

func (p *Parser) parseFunctionBody() []*node.StatementNode {
	p.check(token.OpenBrace)
	// TODO Implement
	p.check(token.CloseBrace)
	return nil
}

// TODO Refactor: Move for loop in parseParameters and parseReturnTypes to own function
func (p *Parser) parseParameters() []*node.VariableNode {
	var parameters []*node.VariableNode

	p.check(token.OpenParen)
	for !p.isEnd() && !p.is(token.CloseParen) {
		parameters = append(parameters, p.parseVariable())
		p.nextToken()
		if p.is(token.Comma) {
			p.nextToken()
		} else if p.is(token.CloseParen) {
			continue
		} else {
			// error
		}
	}
	p.check(token.CloseParen)
	return parameters
}

func (p *Parser) parseReturnTypes() []*node.TypeNode {
	var returnTypes []*node.TypeNode

	if p.is(token.OpenParen){
		p.nextToken()
		for !p.isEnd() && !p.is(token.CloseParen) {
			returnTypes = append(returnTypes, p.parseType())
			p.nextToken()
			if p.is(token.Comma) {
				p.nextToken()
			} else if p.is(token.CloseParen) {
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
	return &node.VariableNode{
		AbstractNode: p.newAbstractNode(),
		Type: p.parseType(),
		Identifier: p.readIdentifier(),
	}
}

func (p * Parser) parseType() *node.TypeNode {
	// Later we need to distinguish between an array and a simple type

	return &node.TypeNode{
		AbstractNode: p.newAbstractNode(),
		Identifier: p.readIdentifier(),
	}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

// helpers
// --------------

func (p *Parser) is(symbol token.Symbol) bool {
	tok, ok := p.currentToken.(*token.FixToken)
	return ok && tok.Value == symbol
}

func (p *Parser) check(symbol token.Symbol) {
	if !p.is(symbol) {
		// todo add to errors
	}
	p.nextToken()
}

func (p *Parser) readIdentifier() string {
	var identifier string

	if tok, ok := p.currentToken.(*token.IdentifierToken); ok {
		identifier = tok.Literal()
	} else {
		// todo add to errors
		identifier = "ERROR"
	}

	p.nextToken()
	return identifier
}

func (p *Parser) newAbstractNode() node.AbstractNode{
	return node.AbstractNode{
		Position: p.currentToken.Pos(),
	}
}

func (p *Parser) isEnd() bool {
	return p.lex.IsEnd
}
