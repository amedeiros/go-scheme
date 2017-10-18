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
	Arguments []Ast
	Function  Ast // Lambda or Identifier
	Token     lexer.Token
}

// Inspect returns the external representation of the expression
func (procCall *ProcedureCall) Inspect() string {
	args := []string{}
	for _, arg := range procCall.Arguments {
		args = append(args, arg.Inspect())
	}

	arguments := ""

	if len(args) > 0 {
		arguments = strings.Join(args, " ")
	}

	funcString := procCall.Function.Inspect()

	return "(" + funcString + " " + arguments + ")"
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

// String node
type String struct {
	Value string
	Token lexer.Token
}

// Inspect the string
func (str *String) Inspect() string {
	return str.Value
}

// LambdaLiteral node
type LambdaLiteral struct {
	Token        lexer.Token
	Paramemeters []*Identifier
	Body         []Ast
}

// Inspect a function
func (funcLit *LambdaLiteral) Inspect() string {
	params := []string{}
	body := []string{}

	for _, str := range funcLit.Paramemeters {
		params = append(params, str.Inspect())
	}

	for _, str := range funcLit.Body {
		body = append(body, str.Inspect())
	}

	arguments := ""

	if len(params) > 0 {
		arguments = strings.Join(params, " ")
	}

	return "(LAMBDA (" + arguments + ") " + strings.Join(body, " ") + ")"
}

type Cons struct {
	Car Ast
	Cdr []Ast
}

func (cons *Cons) Inspect() string {
	body := []string{}
	for _, str := range cons.Cdr {
		body = append(body, str.Inspect())
	}

	if cons.Cdr == nil {
		return cons.Car.Inspect()
	}
	return "(" + cons.Car.Inspect() + " " + strings.Join(body, " ") + ")"
}
