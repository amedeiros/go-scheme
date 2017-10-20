package object

import (
	"fmt"
	"strings"
)

// BuiltinFunction type
type BuiltinFunction func(args ...Object) Object

const (
	INT_OBJ     Type = "INT_OBJ"
	FLOAT_OBJ   Type = "FLOAT_OBJ"
	BUILTIN_OBJ Type = "BUILTIN_OBJ"
	LAMBDA_OBJ  Type = "LAMBDA_OBJ"
	BOOL_OBJ    Type = "BOOL_OBJ"
	STRING_OBJ  Type = "STRING_OBJ"
	ERROR_OBJ   Type = "ERROR_OBJ"
	CONS_OBJ    Type = "CONS_OBJ"
	CHAR_OBJ    Type = "CHAR_OBJ"
	IDENT_OBJ   Type = "IDENT_OBJ"
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
	Value int64
}

// Type of objcet
func (integer *Integer) Type() Type {
	return INT_OBJ
}

// Inspect object
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

// Float type
type Float struct {
	Value float64
}

// Type of objcet
func (f *Float) Type() Type {
	return FLOAT_OBJ
}

// Inspect object
func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
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
	Parameters []*Identifier
	Body       Object
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
		end := c.Cdr.(*Cons)
		cdr := ""

		for {
			cdr += end.Car.Inspect()
			cdr += " "
			if end.Cdr != nil {
				end = end.Cdr.(*Cons)
			} else {
				break
			}
		}
		return "(" + c.Car.Inspect() + " " + strings.TrimSpace(cdr) + ")"
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
