package parser

import (
	"fmt"

	"strconv"

	"github.com/amedeiros/go-scheme/lexer"
)

// Parser reads tokens from the lexer build an AST for the evaluator
type Parser struct {
	lex          *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
	errors       []string
}

// NewParser creates a new parser from our lexer.
func NewParser(lex *lexer.Lexer) *Parser {
	return &Parser{lex: lex, currentToken: lex.NextToken(), peekToken: lex.NextToken()}
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
		return p.parseExpression()
	case lexer.ADD:
		node := &Identifier{Value: p.currentToken.Literal, Token: p.currentToken}
		// p.nextToken()
		return node
	case lexer.DIGIT:
		value, _ := strconv.Atoi(p.currentToken.Literal)
		node := &IntegerLiteral{Value: value, Token: p.currentToken}
		// p.nextToken()
		return node
	case lexer.EOF:
		return nil
	default:
		msg, _ := fmt.Printf("Unkown node literal %s", p.currentToken.Literal)
		panic(msg)
	}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) parseExpression() Ast {
	cons := &Cons{Token: p.currentToken}

	if p.currentToken.Type != lexer.RPAREN {
		cons.Car = p.parseStatement()
		if p.currentToken.Type != lexer.RPAREN && p.currentToken.Type != lexer.EOF {
			cons.Cdr = p.parseStatement()
		}
	}

	if p.currentToken.Type != lexer.RPAREN {
		fmt.Println(p.currentToken.Literal)
		panic("Missing closing )")
	}

	p.nextToken() // consume the closing )

	return cons
}
