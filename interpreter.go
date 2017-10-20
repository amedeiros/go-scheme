package main

import (
	"errors"
	"fmt"

	"github.com/amedeiros/go-scheme/object"
)

func Eval(obj object.Object, env *object.Environment) object.Object {
	switch node := obj.(type) {
	case *object.Boolean, *object.Char, *object.String, *object.Error, *object.Integer, *object.Lambda:
		return obj
	case *object.Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val
		}

		return &object.Error{Value: fmt.Errorf("Unkown identifier %s", node.Value)}
	case *object.Cons:
		car := node.Car
		switch carType := car.(type) {
		case *object.Identifier:
			if builtin, ok := builtins[carType.Value]; ok {
				args := []object.Object{}
				node = node.Cdr.(*object.Cons)

				for node != nil {
					car = node.Car
					val := Eval(car, env)
					if isError(val) {
						return val
					}
					args = append(args, val)
					if node.Cdr != nil {
						node = node.Cdr.(*object.Cons)
					} else {
						node = nil
					}
				}

				return builtin.Fn(args...)
			}

			// Builtin LET
			if carType.Value == "LET" {
				if cons, ok := node.Cdr.(*object.Cons); ok {
					if car, ok := cons.Car.(*object.Identifier); ok {
						val := Eval(cons.Cdr, env)
						env.Set(car.Value, val)
						return val
					}
				} else {
					return &object.Error{Value: errors.New("Expecting cons cell")}
				}
			}

			// Check the ENV for a Lambda
			if val, ok := env.Get(carType.Value); ok {
				if lambda, ok := val.(*object.Lambda); ok {
					var args []object.Object

					cdr := node.Cdr
					for cdr != nil {
						cons := cdr.(*object.Cons)
						val := Eval(cons.Car, env)

						if isError(val) {
							return val

						}

						args = append(args, val)

						if cdr != nil && cons.Cdr != nil {
							cdr = cons.Cdr.(*object.Cons)
						} else {
							break
						}
					}

					return applyFunction(lambda, carType.Value, args)
				}
			}

			return &object.Error{Value: fmt.Errorf("Unkown proc %s", carType.Value)}
		default:
			return Eval(carType, env)
		}
	}

	panic(fmt.Sprintf("OBJECT: %#v", obj))
}

func applyFunction(lambda *object.Lambda, name string, args []object.Object) object.Object {
	extendedEnv := extendFunctionEnv(lambda, name, args)

	return Eval(lambda.Body, extendedEnv)
}

func extendFunctionEnv(lambda *object.Lambda, name string, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(lambda.Env)

	for paramIdx, param := range lambda.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

var builtins = map[string]*object.Builtin{
	"+": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			firstArg := args[0].(*object.Integer)
			intObj := &object.Integer{Value: firstArg.Value}

			for _, arg := range args[1:len(args)] {
				intArg := arg.(*object.Integer)
				intObj.Value += intArg.Value
			}

			return intObj
		},
	},
	"-": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			firstArg := args[0].(*object.Integer)
			intObj := &object.Integer{Value: firstArg.Value}

			for _, arg := range args[1:len(args)] {
				intArg := arg.(*object.Integer)
				intObj.Value -= intArg.Value
			}

			return intObj
		},
	},
	"*": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			firstArg := args[0].(*object.Integer)
			intObj := &object.Integer{Value: firstArg.Value}

			for _, arg := range args[1:len(args)] {
				intArg := arg.(*object.Integer)
				intObj.Value *= intArg.Value
			}

			return intObj
		},
	},
	"/": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			firstArg := args[0].(*object.Integer)
			intObj := &object.Integer{Value: firstArg.Value}

			for _, arg := range args[1:len(args)] {
				intArg := arg.(*object.Integer)
				intObj.Value /= intArg.Value
			}

			return intObj
		},
	},
}

func ap(any interface{}) {
	fmt.Println(fmt.Printf("FANCY: %#v", any))
}

func apMsg(msg string, any interface{}) {
	fmt.Println(fmt.Printf("%s: %#v", msg, any))
}
