package parser

import (
	"fmt"

	"strconv"

	"github.com/amedeiros/go-scheme/lexer"
)

type prefixParseFn func() Ast

// Parser reads tokens from the lexer build an AST for the evaluator
type Parser struct {
	lex            *lexer.Lexer
	currentToken   lexer.Token
	peekToken      lexer.Token
	errors         []string
	prefixParseFns map[lexer.TokenType]prefixParseFn
}

// NewParser creates a new parser from our lexer.
func NewParser(lex *lexer.Lexer) *Parser {
	p := &Parser{lex: lex}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.DIGIT, p.parseDigit)
	p.registerPrefix(lexer.LPAREN, p.parseCallExpression)
	return p
}

// ParseProgram is the entry point to generate the AST
func (p *Parser) ParseProgram() *Program {
	ast := &Program{}

	for p.currentToken.Type != lexer.EOF {
		statement := p.parseStatement()
		if statement != nil {
			ast.Expressions = append(ast.Expressions, statement)
		}

		p.nextToken()
	}

	return ast
}

func (p *Parser) parseStatement() Ast {
	switch p.currentToken.Type {
	case lexer.EOF:
		return nil
	default:
		return p.parseExpression()
	}
}

func (p *Parser) parseExpression() Ast {
	expression := p.prefixParseFns[p.currentToken.Type]

	if expression == nil {
		msg, _ := fmt.Printf("No parse function for literal %s and TokenType %d at %d:%d",
			p.currentToken.Literal, p.currentToken.Type, p.currentToken.Row, p.currentToken.Column)
		panic(msg)
	}

	return expression()
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) parseExpressionEnd(end lexer.TokenType) []Ast {
	ast := []Ast{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return ast
	}

	for !p.currentTokenIs(end) && !p.currentTokenIs(lexer.EOF) {
		ast = append(ast, p.parseExpression())
		p.nextToken()
	}

	if !p.currentTokenIs(end) {
		return nil
	}

	return ast
}

func (p *Parser) parseCallExpression() Ast {
	if p.peekTokenIs(lexer.RPAREN) {
		msg, _ := fmt.Printf("Unexpected ) at %d:%d", p.currentToken.Row, p.currentToken.Column)
		panic(msg)
	}

	if !p.peekTokenIs(lexer.IDENT) {
		msg, _ := fmt.Printf("Unexpected value %s at %d:%d", p.currentToken.Literal, p.currentToken.Row, p.currentToken.Column)
		panic(msg)
	}

	p.nextToken()
	callExp := &ProcedureCall{Token: p.currentToken, Name: p.currentToken.Literal}
	p.nextToken()

	if p.currentTokenIs(lexer.RPAREN) {
		return callExp
	}

	// args := []Ast{}

	// for !p.currentTokenIs(lexer.RPAREN) && !p.currentTokenIs(lexer.EOF) {
	// 	args = append(args, p.parseExpression())
	// 	p.nextToken()
	// }

	// if !p.currentTokenIs(lexer.RPAREN) {
	// 	msg, _ := fmt.Printf("Missing closing ) found %s instead at %d:%d ", p.currentToken.Literal, p.currentToken.Row, p.currentToken.Column)
	// 	panic(msg)
	// }

	// callExp.Arguments = args
	callExp.Arguments = p.parseExpressionEnd(lexer.RPAREN)
	return callExp
}

func (p *Parser) parseIdentifier() Ast {
	return &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseDigit() Ast {
	value, _ := strconv.Atoi(p.currentToken.Literal)
	return &IntegerLiteral{Token: p.currentToken, Value: value}
}

func (p *Parser) currentTokenIs(val lexer.TokenType) bool {
	return p.currentToken.Type == val
}

func (p *Parser) peekTokenIs(val lexer.TokenType) bool {
	return p.peekToken.Type == val
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
