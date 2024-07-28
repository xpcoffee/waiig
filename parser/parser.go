package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// NOTE: the order encodes operator precedence!
const (
	_ int = iota // start with iota to give constants incrementing values
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	INDEX       // [1, 2, 3][5]
	CALL        // myfunction(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(left ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixParseFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixParseFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixParseFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParseFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixParseFn(token.TRUE, p.parseBooleanExpression)
	p.registerPrefixParseFn(token.FALSE, p.parseBooleanExpression)
	p.registerPrefixParseFn(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixParseFn(token.IF, p.parseIfExpression)
	p.registerPrefixParseFn(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefixParseFn(token.STRING, p.parseStringLiteral)
	p.registerPrefixParseFn(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixParseFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixParseFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParseFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixParseFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixParseFn(token.GT, p.parseInfixExpression)
	p.registerInfixParseFn(token.LT, p.parseInfixExpression)
	p.registerInfixParseFn(token.EQ, p.parseInfixExpression)
	p.registerInfixParseFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixParseFn(token.LPAREN, p.parseFunctionCall)
	p.registerInfixParseFn(token.LBRACKET, p.parseIndexingExpression)

	// initialize peek & cur
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefixParseFn(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfixParseFn(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return nil
	}

	p.expectPeek(token.ASSIGN)

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	p.expectPeek(token.SEMICOLON)
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	// expression
	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	// look for OPTIONAL semicolon, and advance past it if we find one (parseExpression won't do this)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parsePrefix := p.prefixParseFns[p.curToken.Type]
	if parsePrefix == nil {
		p.noPrefixParseError(p.curToken.Type)
		return nil
	}
	leftExp := parsePrefix()

	for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.COMMA) && precedence < p.peekPrecedence() {
		parseInfix := p.infixParseFns[p.peekToken.Type]
		if parseInfix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = parseInfix(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function found for %s", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixExp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	prefixExp.Right = p.parseExpression(PREFIX)

	return prefixExp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infixExpression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()

	p.nextToken()
	infixExpression.Right = p.parseExpression(precedence)

	return infixExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	stmt := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	stmt.Value = value
	return stmt
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.BooleanExpression{Token: p.curToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	exp := &ast.FunctionLiteralExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	exp.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Body = p.parseBlockStatement()
	return exp
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}

	for !p.currTokenIs(token.RPAREN) {
		idnt, ok := p.parseIdentifier().(*ast.Identifier)
		if !ok {
			return nil
		}
		parameters = append(parameters, idnt)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
		p.nextToken()
	}

	return parameters
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionCall(expr ast.Expression) ast.Expression {
	exp := &ast.FunctionCallExpression{Token: p.curToken, Function: expr}

	p.nextToken()
	exp.Parameters = p.parseFunctionCallParameters()
	return exp
}

func (p *Parser) parseFunctionCallParameters() []ast.Expression {
	parameters := []ast.Expression{}

	for !p.currTokenIs(token.RPAREN) {
		parameters = append(parameters, p.parseExpression(LOWEST))

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
		p.nextToken()
	}

	return parameters
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	exp := &ast.ArrayLiteral{Token: p.curToken}
	elements := []ast.Expression{}
	p.nextToken()

	for !p.currTokenIs(token.RBRACKET) {
		elements = append(elements, p.parseExpression(LOWEST))

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
		p.nextToken()
	}
	exp.Elements = elements

	return exp
}

func (p *Parser) parseIndexingExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexingExpression{Token: p.curToken, Target: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	p.expectPeek(token.RBRACKET)

	return exp
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("unexpected next token expected=%s got=%s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}
