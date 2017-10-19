package object

import (
	"fmt"

	"github.com/amedeiros/go-scheme/parser"
)

// BuiltinFunction type
type BuiltinFunction func(args ...Object) Object

const (
	INT_OBJ     = "INT_OBJ"
	BUILTIN_OBJ = "BUILTIN_OBJ"
	LAMBDA_OBJ  = "LAMBDA_OBJ"
	BOOL_OBJ    = "BOOL_OBJ"
	STRING_OBJ  = "STRING_OBJ"
	ERROR_OBJ   = "ERROR_OBJ"
	CONS_OBJ    = "CONS_OBJ"
	CHAR_OBJ    = "CHAR_OBJ"
	IDENT_OBJ   = "IDENT_OBJ"
)

// Type represents the type of object
type Type string

// Object interface all objects implement
type Object interface {
	Type() Type
	Inspect() string
}

// Integer type
type Integer struct {
	Value int
}

// Type of objcet
func (integer *Integer) Type() Type {
	return INT_OBJ
}

// Inspect object
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

// Builtin function
type Builtin struct {
	Fn BuiltinFunction
}

// Type of builtin
func (builtin *Builtin) Type() Type {
	return BUILTIN_OBJ
}

// Inspect the builtin
func (builtin *Builtin) Inspect() string {
	return "<#procedure>"
}

// Lambda represents a lambda!
type Lambda struct {
	Parameters []*parser.Identifier
	Body       *parser.Cons
	Env        *Environment
}

// Type of lambda
func (builtin *Lambda) Type() Type {
	return LAMBDA_OBJ
}

// Inspect the builtin
func (builtin *Lambda) Inspect() string {
	return "<#procedure>"
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOL_OBJ
}

func (b *Boolean) Inspect() string {
	if b.Value {
		return "#T"
	}

	return "#F"
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

type Error struct {
	Value error
}

func (e *Error) Type() Type {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return e.Value.Error()
}

type Cons struct {
	Car Object
	Cdr Object
}

func (c *Cons) Type() Type {
	return CONS_OBJ
}

func (c *Cons) Inspect() string {
	if c.Cdr == nil {
		return c.Car.Inspect()
	}

	switch c.Cdr.(type) {
	case *Cons:
		return "(" + c.Car.Inspect() + " " + c.Cdr.Inspect() + ")"
	default:
		return "(" + c.Car.Inspect() + " . " + c.Cdr.Inspect() + ")"
	}
}

type Char struct {
	Value string
}

func (c *Char) Type() Type {
	return CHAR_OBJ
}

func (c *Char) Inspect() string {
	return "#\\" + c.Value
}

type Identifier struct {
	Value string
}

func (i *Identifier) Type() Type {
	return IDENT_OBJ
}

func (i *Identifier) Inspect() string {
	return i.Value
}
