package main

import (
	"fmt"

	"github.com/amedeiros/go-scheme/object"
)

func Eval(obj object.Object, env *object.Environment) object.Object {
	switch node := obj.(type) {
	case *object.Boolean, *object.Char, *object.String:
		return obj
	case *object.Cons:
		car := node.Car
		switch carType := car.(type) {
		case *object.Identifier:
			fmt.Println(fmt.Sprintf("%#v ", carType))
		}
	default:
		panic(fmt.Sprintf("OBJECT: %#v", obj))
	}

	return nil
}
