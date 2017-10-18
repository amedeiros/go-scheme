package parser

import (
	"fmt"

	"strconv"

	"github.com/amedeiros/go-scheme/lexer"
)

type prefixParseFn func() Ast
type reservedParseFn func() Ast

// Parser reads tokens from the lexer build an AST for the evaluator
type Parser struct {
	lex              *lexer.Lexer
	currentToken     lexer.Token
	peekToken        lexer.Token
	errors           []string
	prefixParseFns   map[lexer.TokenType]prefixParseFn
	reservedParseFns map[string]reservedParseFn
}

// NewParser creates a new parser from our lexer.
func NewParser(lex *lexer.Lexer) *Parser {
	p := &Parser{lex: lex}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.DIGIT, p.parseDigit)
	p.registerPrefix(lexer.LPAREN, p.parseList)
	p.registerPrefix(lexer.STRING, p.parseString)

	p.reservedParseFns = make(map[string]reservedParseFn)
	p.registerReservedProc("LAMBDA", p.parseLambda)
	return p
}

// ParseProgram is the entry point to generate the AST
func (p *Parser) ParseProgram() *Cons {
	// ast := &Program{}
	cons := &Cons{}
	cons.Car = p.parseStatement()
	p.nextToken()

	for p.currentToken.Type != lexer.EOF {
		statement := p.parseStatement()
		if statement != nil {
			cons.Cdr = append(cons.Cdr, statement)
		}

		p.nextToken()
	}

	return cons
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
	p.nextToken()
	var ast []Ast

	for !p.currentTokenIs(end) && !p.currentTokenIs(lexer.EOF) {
		p.nextToken()
		ast = append(ast, p.parseExpression())
		p.nextToken()
	}

	if !p.currentTokenIs(end) {
		panic(fmt.Sprintf("Expected TokenType %d found literal %s instead.", end, p.currentToken))
	}

	p.nextToken()

	return ast
}

func (p *Parser) parseList() Ast {
	p.nextToken() // Consume (
	fmt.Println(fmt.Sprintf("BEFORE CAR: %#v", p.currentToken))
	if p.currentTokenIs(lexer.RPAREN) {
		msg, _ := fmt.Printf("Unexpected ) at %d:%d", p.currentToken.Row, p.currentToken.Column)
		panic(msg)
	}

	cons := &Cons{}
	cons.Car = p.parseExpression()
	p.nextToken()
	fmt.Println(fmt.Sprintf("AFTER CAR: %#v", p.currentToken))
	fmt.Println(fmt.Sprintf("CAR %#v", cons.Car))

	if p.currentTokenIs(lexer.RPAREN) {
		p.nextToken()
		return cons
	}

	cdr := []Ast{}

	fmt.Println(fmt.Sprintf("BEFORE CDR: %#v", p.currentToken))
	for !p.currentTokenIs(lexer.RPAREN) && !p.currentTokenIs(lexer.EOF) {
		cdr = append(cdr, p.parseExpression())
		p.nextToken()
	}

	if !p.currentTokenIs(lexer.RPAREN) {
		panic(fmt.Sprintf("Expecting ) at %d:%d", p.currentToken.Row, p.currentToken.Column))
	}

	p.nextToken()

	// cons.Cdr = p.parseExpressionEnd(lexer.RPAREN)
	cons.Cdr = cdr
	fmt.Println(fmt.Sprintf("AFTER CDR: %#v", p.currentToken))
	fmt.Println(fmt.Sprintf("CDR %#v", cons.Cdr))
	// p.nextToken()

	return cons
}

func (p *Parser) parseCallExpression() Ast {
	if p.peekTokenIs(lexer.RPAREN) {
		msg, _ := fmt.Printf("Unexpected ) at %d:%d", p.currentToken.Row, p.currentToken.Column)
		panic(msg)
	}

	// p.nextToken()

	// switch p.currentToken.Type {
	// case lexer.IDENT:
	// 	fmt.Println(">>> IDENT")
	// 	// node := p.currentToken.(lexer.Identifier)

	// 	if reservedProc, ok := p.reservedParseFns[p.currentToken.Literal]; ok {
	// 		fmt.Println(">>> WTF?")
	// 		return reservedProc()
	// 	}

	// 	callExp := &ProcedureCall{Token: node.Token, Function: node}
	// 	// p.nextToken()

	// 	// callExp.Arguments = p.parseExpressionEnd(lexer.RPAREN)
	// 	// return callExp
	// case lexer.LPAREN:
	// 	fmt.Println("???? LPAREN")
	// 	// return &ProcedureCall{Token: node.Token, Function: node}
	// default:
	// 	msg, _ := fmt.Printf("Unexpected value %s at %d:%d", p.currentToken.Literal, p.currentToken.Row, p.currentToken.Column)
	// 	panic(msg)
	// }

	return nil
}

func (p *Parser) parseLambda() Ast {
	lambdaLiteral := &LambdaLiteral{Token: p.currentToken}
	p.nextToken()

	if !p.currentTokenIs(lexer.LPAREN) {
		panic(fmt.Sprintf("Expected ( found %s instead at %d:%d", p.currentToken.Literal, p.currentToken.Row, p.currentToken.Column))
	}

	if p.currentTokenIs(lexer.RPAREN) {
		// no params parse the body and return
		p.nextToken()
		lambdaLiteral.Body = p.parseExpressionEnd(lexer.RPAREN) // &Program{Expressions: p.parseExpressionEnd(lexer.RPAREN)}
		return lambdaLiteral
	}

	// Parse params and body
	lambdaLiteral.Paramemeters = p.parseIdentifierList(lexer.RPAREN)
	p.nextToken()
	lambdaLiteral.Body = p.parseExpressionEnd(lexer.RPAREN)

	return lambdaLiteral
}

func (p *Parser) parseIdentifierList(end lexer.TokenType) []*Identifier {
	l := p.parseExpressionEnd(end)
	i := []*Identifier{}

	for _, val := range l {
		i = append(i, val.(*Identifier))
	}

	return i
}

func (p *Parser) parseIdentifier() Ast {
	return &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseString() Ast {
	return &String{Value: p.currentToken.Literal, Token: p.currentToken}
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

func (p *Parser) registerReservedProc(name string, fn reservedParseFn) {
	p.reservedParseFns[name] = fn
}
