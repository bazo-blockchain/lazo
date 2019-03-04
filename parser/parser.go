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
			// parse variable (or later other types of nodes)
			contract.Variables = append(contract.Variables, p.parseVariable())

		default:
			// error
		}
	}

	p.check(token.CloseBrace)
	return contract
}

func (p *Parser) parseVariable() *node.VariableNode {
	variable := &node.VariableNode{
		AbstractNode: p.newAbstractNode(),
		Type: p.parseType(),
		Identifier: p.readIdentifier(),
	}
	return variable
}

func (p * Parser) parseType() *node.TypeNode {
	// Later we need to distinguish between an array and a simple type

	typeNode := &node.TypeNode{
		AbstractNode: p.newAbstractNode(),
		Identifier: p.readIdentifier(),
	}
	return typeNode
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
