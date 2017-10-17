package parser

import (
	"fmt"

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
	p := &Parser{lex: lex, currentToken: lex.NextToken(), peekToken: lex.NextToken()}

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	return p
}

// ParseProgram is the entry point to generate the AST
func (p *Parser) ParseProgram() []Ast {
	ast := []Ast{}

	for p.currentToken.Type != lexer.EOF {
		statement := p.parseStatement()
		if statement != nil {
			ast = append(ast, statement)
		}

		p.nextToken()
	}

	return ast
}

func (p *Parser) parseStatement() Ast {
	switch p.currentToken.Type {
	case lexer.LPAREN:
		p.nextToken()
		return p.parseCallExpression()
	// case lexer.IDENT:
	// 	return &Identifier{Value: p.currentToken.Literal, Token: p.currentToken}
	// case lexer.DIGIT:
	// 	value, _ := strconv.Atoi(p.currentToken.Literal)
	// 	return &IntegerLiteral{Value: value, Token: p.currentToken}
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

func (p *Parser) parseCallExpression() Ast {
	fmt.Print(p.currentToken.Literal)

	if p.currentTokenIs(lexer.RPAREN) {
		msg, _ := fmt.Printf("Unexpected ) at %d:%d", p.currentToken.Row, p.currentToken.Column)
		panic(msg)
	}

	fmt.Print(p.currentToken.Literal)

	if !p.currentTokenIs(lexer.IDENT) {
		msg, _ := fmt.Printf("Unexpected value %s at %d:%d", p.currentToken.Literal, p.currentToken.Row, p.currentToken.Column)
		panic(msg)
	}

	callExp := &ProcedureCall{Token: p.currentToken, Name: p.currentToken.Literal}

	p.nextToken()
	if p.currentTokenIs(lexer.RPAREN) {
		return callExp
	}

	args := []Ast{}
	for !p.currentTokenIs(lexer.RPAREN) && !p.currentTokenIs(lexer.EOF) {
		fmt.Println(p.currentToken.Literal)
		args = append(args, p.parseStatement())
	}

	if !p.currentTokenIs(lexer.RPAREN) {
		panic("Missing closing )")
	}
	p.nextToken()

	callExp.Arguments = args
	return callExp
}

func (p *Parser) parseIdentifier() Ast {
	return &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) currentTokenIs(val lexer.TokenType) bool {
	return p.currentToken.Type == val
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
