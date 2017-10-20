package main

import (
	"errors"
	"fmt"

	"github.com/amedeiros/go-scheme/object"
)

// LET to check for a let call
const LET = "LET"

// Eval an object
func Eval(obj object.Object, env *object.Environment) object.Object {
	switch node := obj.(type) {
	case *object.Boolean, *object.Char, *object.String, *object.Error, *object.Integer:
		return obj
	case *object.Lambda:
		node.Env = env
		return node
	case *object.Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val
		}

		return newError(fmt.Sprintf("Unkown identifier %s", node.Value))
	case *object.Cons:
		car := node.Car
		switch carType := car.(type) {
		case *object.Identifier:
			if builtin, ok := builtins[carType.Value]; ok {
				args, err := evalArgs(node.Cdr.(*object.Cons), env)
				if err != nil {
					return err
				}

				return builtin.Fn(args...)
			}

			if scopedBuiltin, ok := scopedBuiltins[carType.Value]; ok {
				args, err := evalArgs(node.Cdr.(*object.Cons), env)
				if err != nil {
					return err
				}

				return scopedBuiltin.Fn(env, args...)
			}

			// Builtin LET
			if carType.Value == LET {
				if cons, ok := node.Cdr.(*object.Cons); ok {
					if car, ok := cons.Car.(*object.Identifier); ok {
						val := Eval(cons.Cdr, env)
						env.Set(car.Value, val)
						return val
					}
				} else {
					return newError("Expecting cons cell")
				}
			}

			// Check the ENV for a Lambda
			if val, ok := env.Get(carType.Value); ok {
				if lambda, ok := val.(*object.Lambda); ok {
					args, err := evalArgs(node.Cdr.(*object.Cons), env)
					if err != nil {
						return err
					}

					return applyFunction(lambda, carType.Value, args)
				}
			}

			return newError(fmt.Sprintf("Unkown proc %s", carType.Value))
		default:
			return Eval(carType, env)
		}
	}

	panic("You just found a bug or an unimplemented feature congrats!")
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

var scopedBuiltins = map[string]*object.ScopedBuiltin{}

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
	"QUOTE": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return args[0].(*object.String)
		},
	},
}

func errorObject(err error) *object.Error {
	return &object.Error{Value: err}
}

func newError(msg string) *object.Error {
	return errorObject(errors.New(msg))
}

func ap(any interface{}) {
	fmt.Println(fmt.Printf("%#v", any))
}

func apMsg(msg string, any interface{}) {
	fmt.Println(fmt.Printf("%s: %#v", msg, any))
}

// Load setups the inital environment and returns it
func Load() *object.Environment {
	loadScopedBuiltins()
	return object.NewEnvironment()
}

func evalArgs(cons *object.Cons, env *object.Environment) ([]object.Object, *object.Error) {
	args := []object.Object{}

	for {
		car := cons.Car
		val := Eval(car, env)
		if isError(val) {
			return nil, val.(*object.Error)
		}
		args = append(args, val)
		if cons.Cdr != nil {
			cons = cons.Cdr.(*object.Cons)
		} else {
			break
		}
	}

	return args, nil
}

func loadScopedBuiltins() {
	eval := &object.ScopedBuiltin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			r := NewReader(args[0].(*object.String).Value)
			return Eval(r.Read(), env)
		},
	}

	scopedBuiltins["EVAL"] = eval
}
