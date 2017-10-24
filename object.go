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
	String() string
}

// Integer type
type Integer struct {
	Value int64
}

// Inspect object
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// String object
func (i *Integer) String() string {
	return i.Inspect()
}

// Float type
type Float struct {
	Value float64
}

// Inspect object
func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

// String object
func (f *Float) String() string {
	return f.Inspect()
}

// Builtin function
type Builtin struct {
	Fn BuiltinFunction
}

// Inspect the builtin
func (b *Builtin) Inspect() string {
	return "<#procedure>"
}

// String for builtin
func (b *Builtin) String() string {
	return b.Inspect()
}

// ScopedBuiltin function
type ScopedBuiltin struct {
	Fn  ScopedBuiltinFunction
	Env *Environment
}

// Inspect the builtin
func (b *ScopedBuiltin) Inspect() string {
	return "<#procedure>"
}

// String
func (b *ScopedBuiltin) String() string {
	return b.Inspect()
}

// Lambda represents a lambda!
type Lambda struct {
	Parameters []*Identifier
	Body       Object
	Env        *Environment
}

// Inspect the builtin
func (l *Lambda) Inspect() string {
	str := "(lambda ("
	args := []string{}
	for _, arg := range l.Parameters {
		args = append(args, arg.Inspect())
	}

	str += strings.Join(args, " ") + ") "
	str += l.Body.Inspect() + ")"

	return str
}

func (l *Lambda) String() string {
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

// Inspect the boolean
func (b *Boolean) String() string {
	return b.Inspect()
}

// String represents a string in scheme
type String struct {
	Value string
}

// Inspect the string
func (s *String) Inspect() string {
	return s.Value
}

// String
func (s *String) String() string {
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

// String
func (e *Error) String() string {
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

// String
func (p *Pair) String() string {
	return p.Inspect()
}

// Char representation
type Char struct {
	Value string
}

// Inspect a char
func (c *Char) Inspect() string {
	return "#\\" + c.Value
}

// String
func (c *Char) String() string {
	return c.Inspect()
}

// Identifier is a symbol
type Identifier struct {
	Value string
}

// Inspect the identifier
func (i *Identifier) Inspect() string {
	return i.Value
}

// String the identifier
func (i *Identifier) String() string {
	return i.Inspect()
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

func (v *Vector) String() string {
	return v.Inspect()
}

// Data object
type Data struct {
	Value string
}

// Inspect the data
func (d *Data) Inspect() string {
	return d.Value
}

func (d *Data) String() string {
	return d.Value
}
