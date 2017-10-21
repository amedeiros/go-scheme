package main

import (
	"fmt"
	"strings"
)

// BuiltinFunction type
type BuiltinFunction func(args ...Object) Object

//ScopedBuiltinFunction type
type ScopedBuiltinFunction func(env *Environment, args ...Object) Object

// Object interface all objects implement
type Object interface {
	Inspect() string
}

// Integer type
type Integer struct {
	Value int64
}

// Inspect object
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

// Float type
type Float struct {
	Value float64
}

// Inspect object
func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

// Builtin function
type Builtin struct {
	Fn BuiltinFunction
}

// Inspect the builtin
func (builtin *Builtin) Inspect() string {
	return "<#procedure>"
}

// ScopedBuiltin function
type ScopedBuiltin struct {
	Fn  ScopedBuiltinFunction
	Env *Environment
}

// Inspect the builtin
func (builtin *ScopedBuiltin) Inspect() string {
	return "<#procedure>"
}

// Lambda represents a lambda!
type Lambda struct {
	Parameters []*Identifier
	Body       Object
	Env        *Environment
}

// Inspect the builtin
func (builtin *Lambda) Inspect() string {
	return "<#procedure>"
}

type Boolean struct {
	Value bool
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

func (s *String) Inspect() string {
	return s.Value
}

type Error struct {
	Value error
}

func (e *Error) Inspect() string {
	return e.Value.Error()
}

type Cons struct {
	Car Object
	Cdr Object
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

func (c *Char) Inspect() string {
	return "#\\" + c.Value
}

type Identifier struct {
	Value string
}

func (i *Identifier) Inspect() string {
	return i.Value
}
