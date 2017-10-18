package object

import "fmt"

// BuiltinFunction type
type BuiltinFunction func(args ...Object) Object

const (
	INT_OBJ     = "INT_OBJ"
	BUILTIN_OBJ = "BUILTIN_OBJ"
)

// Type represents the type of object
type Type string

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
