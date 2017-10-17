package object

import "fmt"

const (
	INT_OBJ = "INT_OBJ"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int
}

// Type of objcet
func (integer *Integer) Type() string {
	return INT_OBJ
}

// Inspect object
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}
