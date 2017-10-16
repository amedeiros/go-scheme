package parser

import (
	"github.com/amedeiros/go-scheme/lexer"
)

// Ast is our base interface for all other AST types
type Ast interface {
	GetToken() lexer.Token
	Inspect() string
}

// IntegerLiteral is an interger wrapped as an AST node
type IntegerLiteral struct {
	Token lexer.Token
	Value int
}

// Inspect returns the token literal.
func (intLiteral *IntegerLiteral) Inspect() string {
	return intLiteral.Token.Literal
}

// GetToken returns the token.
func (intLiteral *IntegerLiteral) GetToken() lexer.Token {
	return intLiteral.Token
}

// Cons is our cons cell representation
type Cons struct {
	Car   Ast
	Cdr   Ast
	Token lexer.Token // '('
}

// Inspect returns the token literal.
func (cons *Cons) Inspect() string {
	// return intLiteral.Token.Literal
	car := ""
	cdr := ""

	if cons.Car != nil {
		car = cons.Car.Inspect()
	}

	if cons.Cdr != nil {
		cdr = cons.Cdr.Inspect()
	}

	return "(" + car + " " + cdr + ")"
}

// GetToken returns the token.
func (cons *Cons) GetToken() lexer.Token {
	return cons.Token
}

// Identifier represents an identifier +, =, apples, oranges etc
type Identifier struct {
	Value string
	Token lexer.Token
}

// Inspect returns the token literal.
func (ident *Identifier) Inspect() string {
	return ident.Token.Literal
}

// GetToken returns the token.
func (ident *Identifier) GetToken() lexer.Token {
	return ident.Token
}
