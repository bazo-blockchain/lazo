// Package parser performs syntactic analysis and creates nodes.
// It takes the token stream from lexer, recognizes the nodes and outputs an abstract syntax tree (AST).
package parser

import (
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser/node"
	"math/big"
)

// Parser is a LL(k=2) parser, which means "Left-to-right, Leftmost derivation" top-down parser.
// It holds 2 lookahead tokens (current and peek token) from the given lexer to parse the input.
// It also collects all the syntactic errors.
type Parser struct {
	lex          *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []error
}

// New creates a new Parser struct with the given lexer.
// Since it is a LL(2) parser, it reads the next two tokens and initializes current and peek tokens.
// It returns the created parser struct.
func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		lex: lex,
	}

	p.nextToken()
	p.nextTokenWhileNewLine()

	return p
}

// ParseProgram reads token by token from lexer and creates a ProgramNode (aka. abstract syntax tree).
//
// The syntax tree consists of nodes. Every node stands for a construct of the source code.
// Abstract in this context means that not every detail is captured in the syntax tree.
// For example, already recognized keywords (e.g. 'contract', 'if' etc.) and fix symbols (e.g. comma, parentheses etc.)
// are skipped, since they are not relevant for further steps.
//
// It returns the parsed ProgramNode/syntax tree and syntactic errors
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

	p.checkAndSkipNewLines(token.CloseBrace)
	return contract
}

func (p *Parser) parseContractBody(contract *node.ContractNode) {
	switch p.currentToken.Type() {
	case token.IDENTIFER:
		contract.Fields = append(contract.Fields, p.parseField())
	case token.SYMBOL:
		ftok, _ := p.currentToken.(*token.FixToken)

		switch ftok.Value {
		case token.Function:
			contract.Functions = append(contract.Functions, p.parseFunction())
		case token.Constructor:
			if contract.Constructor == nil {
				contract.Constructor = p.parseConstructor()
			} else {
				p.addError(fmt.Sprintf("Only one constructor is allowed"))
				p.nextToken()
			}
		case token.Struct:
			contract.Structs = append(contract.Structs, p.parseStruct())
		case token.Map:
			contract.Fields = append(contract.Fields, p.parseField())
		default:
			p.addError(fmt.Sprintf("Unsupported symbol %s in contract", ftok.Lexeme))
			p.nextToken()
		}
	default:
		p.addError("Unsupported contract part: " + p.currentToken.Literal())
		p.nextToken()
	}
}

func (p *Parser) parseField() *node.FieldNode {
	v := &node.FieldNode{
		AbstractNode: p.newAbstractNode(),
		Type:         p.parseType(),
		Identifier:   p.readIdentifier(),
	}

	if p.isSymbol(token.Assign) {
		p.nextToken()
		v.Expression = p.parseExpression()
	}
	p.checkAndSkipNewLines(token.NewLine)
	return v
}

func (p *Parser) parseStruct() *node.StructNode {
	s := &node.StructNode{
		AbstractNode: p.newAbstractNode(),
	}
	p.nextToken() // skip 'struct' keyword

	s.Name = p.readIdentifier()
	p.check(token.OpenBrace)
	p.checkAndSkipNewLines(token.NewLine)

	for !p.isEnd() && !p.isSymbol(token.CloseBrace) {
		f := &node.StructFieldNode{
			AbstractNode: p.newAbstractNode(),
			Type:         p.parseType(),
			Identifier:   p.readIdentifier(),
		}
		p.checkAndSkipNewLines(token.NewLine)
		s.Fields = append(s.Fields, f)
	}

	p.check(token.CloseBrace)
	p.checkAndSkipNewLines(token.NewLine)

	return s
}

func (p *Parser) parseConstructor() *node.ConstructorNode {
	constructor := &node.ConstructorNode{
		AbstractNode: p.newAbstractNode(),
	}
	p.nextToken() // skip constructor keyword

	constructor.Parameters = p.parseParameters()
	constructor.Body = p.parseStatementBlock()

	return constructor
}

func (p *Parser) parseFunction() *node.FunctionNode {
	function := &node.FunctionNode{
		AbstractNode: p.newAbstractNode(),
	}
	p.nextToken() // skip function keyword

	function.ReturnTypes = p.parseReturnTypes()
	function.Name = p.readIdentifier()
	function.Parameters = p.parseParameters()
	function.Body = p.parseStatementBlock()

	return function
}

func (p *Parser) parseReturnTypes() []node.TypeNode {
	var returnTypes []node.TypeNode

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

func (p *Parser) parseParameters() []*node.ParameterNode {
	var parameters []*node.ParameterNode

	p.check(token.OpenParen)
	isFirstParam := true
	for !p.isEnd() && !p.isSymbol(token.CloseParen) {
		if !isFirstParam {
			p.checkAndSkipNewLines(token.Comma)
		}
		param := &node.ParameterNode{
			AbstractNode: p.newAbstractNode(),
			Type:         p.parseType(),
			Identifier:   p.readIdentifier(),
		}
		parameters = append(parameters, param)
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
	p.check(token.CloseBrace)
	if !p.isSymbol(token.Else) {
		p.checkAndSkipNewLines(token.NewLine)
	}
	return statements
}

func (p *Parser) parseStatement() node.StatementNode {
	switch p.currentToken.Type() {
	case token.IDENTIFER:
		return p.parseStatementWithIdentifier()
	case token.SYMBOL:
		return p.parseStatementWithFixToken()
	default:
		p.addError("Unsupported statement starting with " + p.currentToken.Literal())
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseStatementWithIdentifier() node.StatementNode {
	abstractNode := p.newAbstractNode()
	identifier := p.readIdentifier()

	if p.currentToken.Type() == token.IDENTIFER || p.isSymbol(token.OpenBracket) && p.peekIsSymbol(token.CloseBracket) {
		return p.parseVariableStatementWithIdentifier(abstractNode, identifier)
	}
	designator := p.parseDesignatorWithIdentifier(abstractNode, identifier)
	if p.isType(token.SYMBOL) {
		tok := p.currentToken.(*token.FixToken)
		switch tok.Value {
		case token.Assign:
			return p.parseAssignmentStatement(designator)
		case token.Comma:
			return p.parseMultiAssignmentStatement(designator)
		case token.OpenParen:
			return p.parseCallStatement(designator)
		default:
			return p.parseShorthandAssignmentStatement(designator, tok.Value)
		}
	}

	if p.isType(token.IDENTIFER) {
		p.addError("Invalid Array declaration")
		return nil
	}

	p.addError("not yet implemented " + p.currentToken.Literal())
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
	case token.Map:
		return p.parseVariableStatement()
	case token.Delete:
		return p.parseDeleteStatement()
	default:
		p.addError("Unsupported statement starting with " + ftok.Literal())
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseVariableStatement() node.StatementNode {
	varType := p.parseType()
	return p.parseVariableStatementWithType(varType)
}

func (p *Parser) parseVariableStatementWithIdentifier(abstractNode node.AbstractNode, identifier string) node.StatementNode {
	varType := p.parseTypeWithIdentifier(abstractNode, identifier)
	return p.parseVariableStatementWithType(varType)
}

func (p *Parser) parseVariableStatementWithType(varType node.TypeNode) node.StatementNode {
	id := p.readIdentifier()

	if p.isSymbol(token.Comma) {
		return p.parseMultiVariableStatement(varType, id)
	}

	v := &node.VariableNode{
		AbstractNode: p.newAbstractNodeWithPos(varType.Pos()),
		Type:         varType,
		Identifier:   id,
	}
	if p.isSymbol(token.Assign) {
		p.nextToken()
		v.Expression = p.parseExpression()
	}
	p.checkAndSkipNewLines(token.NewLine)
	return v
}

func (p *Parser) parseMultiVariableStatement(varType node.TypeNode, id string) *node.MultiVariableNode {
	types := []node.TypeNode{varType}
	ids := []string{id}

	if !p.isEnd() && p.isSymbol(token.Comma) {
		p.nextToken()
		types = append(types, p.parseType())
		ids = append(ids, p.readIdentifier())
	}
	p.check(token.Assign)
	mv := &node.MultiVariableNode{
		AbstractNode: p.newAbstractNodeWithPos(varType.Pos()),
		Types:        types,
		Identifiers:  ids,
		FuncCall:     p.parseFuncCall(p.parseDesignator()),
	}
	p.checkAndSkipNewLines(token.NewLine)
	return mv
}

func (p *Parser) parseAssignmentStatement(left node.DesignatorNode) *node.AssignmentStatementNode {
	p.nextToken() // skip '=' sign

	expression := p.parseExpression()
	p.checkAndSkipNewLines(token.NewLine)

	return &node.AssignmentStatementNode{
		AbstractNode: p.newAbstractNodeWithPos(left.Pos()),
		Left:         left,
		Right:        expression,
	}
}

func (p *Parser) parseMultiAssignmentStatement(designator node.DesignatorNode) *node.MultiAssignmentStatementNode {
	designators := []node.DesignatorNode{designator}

	if !p.isEnd() && p.isSymbol(token.Comma) {
		p.nextToken()
		designators = append(designators, p.parseDesignator())
	}
	p.check(token.Assign)
	ma := &node.MultiAssignmentStatementNode{
		AbstractNode: p.newAbstractNodeWithPos(designator.Pos()),
		Designators:  designators,
		FuncCall:     p.parseFuncCall(p.parseDesignator()),
	}
	p.checkAndSkipNewLines(token.NewLine)
	return ma
}

var allowedShorthandOperators = []token.Symbol{
	token.Plus,
	token.Minus,
	token.Multiplication,
	token.Division,
	token.Modulo,
	token.Exponent,
}

func (p *Parser) parseShorthandAssignmentStatement(designator node.DesignatorNode, operator token.Symbol) node.StatementNode {
	if !containsSymbol(allowedShorthandOperators, operator) {
		p.addError(fmt.Sprintf("Unsupported symbol %s", token.SymbolLexeme[operator]))
		return nil
	}

	assignment := &node.ShorthandAssignmentStatementNode{
		AbstractNode: p.newAbstractNodeWithPos(designator.Pos()),
		Designator:   designator,
		Operator:     operator,
	}

	p.nextToken()
	if !p.isType(token.SYMBOL) {
		p.addError(fmt.Sprintf("Symbol token expected"))
		return nil
	}

	// postfix operator (x++ or x--)
	if p.isAnySymbol(token.Plus, token.Minus) && p.isSymbol(operator) {
		assignment.Expression = &node.IntegerLiteralNode{
			AbstractNode: p.newAbstractNode(),
			Value:        big.NewInt(1),
		}
		p.nextToken()
	} else {
		// x += 2
		p.check(token.Assign)
		assignment.Expression = p.parseExpression()
	}

	p.checkAndSkipNewLines(token.NewLine)
	return assignment
}

func (p *Parser) parseIfStatement() *node.IfStatementNode {
	abstractNode := p.newAbstractNode()
	p.nextToken() // skip 'if' keyword

	// Condition
	p.check(token.OpenParen)
	condition := p.parseExpression()
	p.check(token.CloseParen)

	// Then
	then := p.parseStatementBlock()

	var alternative []node.StatementNode

	if p.isSymbol(token.Else) {
		p.nextToken() // skip 'else' keyword

		// Else
		alternative = p.parseStatementBlock()
	}

	return &node.IfStatementNode{
		AbstractNode: abstractNode,
		Condition:    condition,
		Then:         then,
		Else:         alternative,
	}
}

func (p *Parser) parseReturnStatement() *node.ReturnStatementNode {
	var returnValues []node.ExpressionNode
	abstractNode := p.newAbstractNode()

	p.nextToken() // skip 'return' keyword

	if !p.isSymbol(token.NewLine) {
		returnValues = append(returnValues, p.parseExpression())

		for !p.isEnd() && p.isSymbol(token.Comma) {
			p.nextToken() // skip ','
			returnValues = append(returnValues, p.parseExpression())
		}
	}

	p.checkAndSkipNewLines(token.NewLine)

	return &node.ReturnStatementNode{
		AbstractNode: abstractNode,
		Expressions:  returnValues,
	}
}

func (p *Parser) parseCallStatement(designator node.DesignatorNode) *node.CallStatementNode {
	fc := p.parseFuncCall(designator)
	p.checkAndSkipNewLines(token.NewLine)

	return &node.CallStatementNode{
		AbstractNode: fc.AbstractNode,
		Call:         fc,
	}
}

func (p *Parser) parseDeleteStatement() *node.DeleteStatementNode {
	d := &node.DeleteStatementNode{
		AbstractNode: p.newAbstractNode(),
	}

	p.check(token.Delete)
	designator := p.parseDesignator()
	if elementAccess, ok := designator.(*node.ElementAccessNode); ok {
		d.Element = elementAccess
	} else {
		p.addError("delete requires element access expression")
	}
	p.checkAndSkipNewLines(token.NewLine)
	return d
}

// ----------------------------------------- End of Statements

func (p *Parser) parseType() node.TypeNode {
	if p.isType(token.IDENTIFER) {
		return p.parseTypeWithIdentifier(p.newAbstractNode(), p.readIdentifier())
	}

	if p.isSymbol(token.Map) {
		return p.parseMapType()
	}
	return p.newErrorNode("Invalid type")
}

func (p *Parser) parseTypeWithIdentifier(abstractNode node.AbstractNode, identifier string) node.TypeNode {
	typeNode := &node.BasicTypeNode{
		AbstractNode: abstractNode,
		Identifier:   identifier,
	}

	if p.isSymbol(token.OpenBracket) {
		return p.parseArrayType(typeNode)
	}

	return typeNode
}

func (p *Parser) parseArrayType(arrayType node.TypeNode) node.TypeNode {
	p.nextToken() // Skip '[' symbol
	p.check(token.CloseBracket)

	arrayTypeNode := &node.ArrayTypeNode{
		AbstractNode: p.newAbstractNodeWithPos(arrayType.Pos()),
		ElementType:  arrayType,
	}
	if p.isSymbol(token.OpenBracket) {
		return p.parseArrayType(arrayTypeNode)
	}
	return arrayTypeNode
}

func (p *Parser) parseMapType() node.TypeNode {
	mapType := &node.MapTypeNode{
		AbstractNode: p.newAbstractNode(),
	}
	p.check(token.Map)
	p.check(token.Less)
	mapType.KeyType = p.parseType()
	p.check(token.Comma)
	mapType.ValueType = p.parseType()
	p.check(token.Greater)

	return mapType
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
		return containsSymbol(expectedSymbols, tok.Value)
	}
	return false
}

func containsSymbol(symbols []token.Symbol, symbol token.Symbol) bool {
	for _, s := range symbols {
		if symbol == s {
			return true
		}
	}
	return false
}

func (p *Parser) peekIsSymbol(symbol token.Symbol) bool {
	tok, ok := p.peekToken.(*token.FixToken)
	return ok && tok.Value == symbol
}

func (p *Parser) check(symbol token.Symbol) {
	if !p.isSymbol(symbol) {
		var lexeme string
		if ftok, ok := p.currentToken.(*token.FixToken); ok {
			lexeme = token.SymbolLexeme[ftok.Value]
		} else {
			lexeme = p.currentToken.Literal()
		}
		p.addError(fmt.Sprintf("Symbol %s expected, but got %s", token.SymbolLexeme[symbol], lexeme))
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

func (p *Parser) newAbstractNodeWithPos(pos token.Position) node.AbstractNode {
	return node.AbstractNode{
		Position: pos,
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
	return p.isSymbol(token.EOF)
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors,
		fmt.Errorf("[%s] ERROR: %s", p.currentToken.Pos().String(), msg))
}
