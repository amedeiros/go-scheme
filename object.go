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
	Data       bool
}

// Inspect the builtin
func (l *Lambda) Inspect() string {
	if l.Data {
		str := "(lambda ("
		args := []string{}
		for _, arg := range l.Parameters {
			args = append(args, arg.Inspect())
		}

		str += strings.Join(args, " ") + ") "
		str += l.Body.Inspect() + ")"

		return str
	}

	return "<#procedure>"
}

// Boolean representation
type Boolean struct {
	Value bool
}

// Inspect the boolean
func (b *Boolean) Inspect() string {
	if b.Value {
		return "#T"
	}

	return "#F"
}

// String represents a string in scheme
type String struct {
	Value string
}

// Inspect the string
func (s *String) Inspect() string {
	return "\"" + s.Value + "\""
}

// Error wraps a go error
type Error struct {
	Value error
}

// Inspect the error
func (e *Error) Inspect() string {
	return e.Value.Error()
}

// Pair represents a pair of cons cells
type Pair struct {
	Car Object
	Cdr Object
}

// Inspect the pair
func (c *Pair) Inspect() string {
	if c.Cdr == nil && c.Car == nil {
		return "()"
	} else if c.Cdr == nil {
		return "(" + c.Car.Inspect() + ")"
	}

	switch c.Cdr.(type) {
	case *Pair:
		end := c.Cdr.(*Pair)
		cdr := ""

		for {
			cdr += end.Car.Inspect()
			cdr += " "
			if end.Cdr != nil {
				end = end.Cdr.(*Pair)
			} else {
				break
			}
		}
		return "(" + c.Car.Inspect() + " " + strings.TrimSpace(cdr) + ")"
	default:
		return "(" + c.Car.Inspect() + " . " + c.Cdr.Inspect() + ")"
	}
}

// Char representation
type Char struct {
	Value string
}

// Inspect a char
func (c *Char) Inspect() string {
	return "#\\" + c.Value
}

// Identifier is a symbol
type Identifier struct {
	Value string
}

// Inspect the identifier
func (i *Identifier) Inspect() string {
	return i.Value
}

// Vector wraps a Pair
type Vector struct {
	Value []Object
}

// Inspect return the inspect of the Pair prepending #
func (v *Vector) Inspect() string {
	str := "#("

	for _, obj := range v.Value {
		str += obj.Inspect()
		str += " "
	}

	str = strings.TrimSpace(str)
	str += ")"

	return str
}

// Data object
type Data struct {
	Value string
}

// Inspect the data
func (d *Data) Inspect() string {
	return d.Value
}
