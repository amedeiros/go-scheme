package parser

import (
	"strings"

	"github.com/amedeiros/go-scheme/lexer"
)

// Ast is our base interface for all other AST types
type Ast interface {
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

// ProcedureCall calls a procedure (lambda)
type ProcedureCall struct {
	Name      string
	Arguments []Ast
	Token     lexer.Token
}

// Inspect returns the external representation of the expression
func (procCall *ProcedureCall) Inspect() string {
	args := []string{}
	for _, arg := range procCall.Arguments {
		args = append(args, arg.Inspect())
	}

	return "(" + procCall.Name + " " + strings.Join(args, " ") + ")"
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
