package main

import (
	"errors"
	"fmt"
)

// SET to check for a let call
const SET = "SET"

// Eval an object
func Eval(obj Object, env *Environment) Object {
	switch node := obj.(type) {
	case *Boolean, *Char, *String, *Error, *Integer:
		return obj
	case *Lambda:
		node.Env = env
		return node
	case *Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val
		}

		return newError(fmt.Sprintf("Unkown identifier %s", node.Value))
	case *Cons:
		car := node.Car
		switch carType := car.(type) {
		case *Identifier:
			if builtin, ok := builtins[carType.Value]; ok {
				args, err := evalArgs(node.Cdr.(*Cons), env)
				if err != nil {
					return err
				}

				return builtin.Fn(args...)
			}

			if scopedBuiltin, ok := scopedBuiltins[carType.Value]; ok {
				args, err := evalArgs(node.Cdr.(*Cons), env)
				if err != nil {
					return err
				}

				return scopedBuiltin.Fn(env, args...)
			}

			// Builtin SET
			if carType.Value == SET {
				if cons, ok := node.Cdr.(*Cons); ok {
					if car, ok := cons.Car.(*Identifier); ok {
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
				if lambda, ok := val.(*Lambda); ok {
					args, err := evalArgs(node.Cdr.(*Cons), env)
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

func applyFunction(lambda *Lambda, name string, args []Object) Object {
	extendedEnv := extendFunctionEnv(lambda, name, args)

	return Eval(lambda.Body, extendedEnv)
}

func extendFunctionEnv(lambda *Lambda, name string, args []Object) *Environment {
	env := NewEnclosedEnvironment(lambda.Env)

	for paramIdx, param := range lambda.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

var scopedBuiltins = map[string]*ScopedBuiltin{}

var builtins = map[string]*Builtin{
	"+": &Builtin{
		Fn: func(args ...Object) Object {
			firstArg := args[0].(*Integer)
			intObj := &Integer{Value: firstArg.Value}

			for _, arg := range args[1:len(args)] {
				intArg := arg.(*Integer)
				intObj.Value += intArg.Value
			}

			return intObj
		},
	},
	"-": &Builtin{
		Fn: func(args ...Object) Object {
			firstArg := args[0].(*Integer)
			intObj := &Integer{Value: firstArg.Value}

			for _, arg := range args[1:len(args)] {
				intArg := arg.(*Integer)
				intObj.Value -= intArg.Value
			}

			return intObj
		},
	},
	"*": &Builtin{
		Fn: func(args ...Object) Object {
			firstArg := args[0].(*Integer)
			intObj := &Integer{Value: firstArg.Value}

			for _, arg := range args[1:len(args)] {
				intArg := arg.(*Integer)
				intObj.Value *= intArg.Value
			}

			return intObj
		},
	},
	"/": &Builtin{
		Fn: func(args ...Object) Object {
			firstArg := args[0].(*Integer)
			intObj := &Integer{Value: firstArg.Value}

			for _, arg := range args[1:len(args)] {
				intArg := arg.(*Integer)
				intObj.Value /= intArg.Value
			}

			return intObj
		},
	},
	"QUOTE": &Builtin{
		Fn: func(args ...Object) Object {
			return args[0].(*String)
		},
	},
}

func errorObject(err error) *Error {
	return &Error{Value: err}
}

func newError(msg string) *Error {
	return errorObject(errors.New(msg))
}

func ap(any interface{}) {
	fmt.Println(fmt.Printf("%#v", any))
}

func apMsg(msg string, any interface{}) {
	fmt.Println(fmt.Printf("%s: %#v", msg, any))
}

// Load setups the inital environment and returns it
func Load() *Environment {
	loadScopedBuiltins()
	return NewEnvironment()
}

func evalArgs(cons *Cons, env *Environment) ([]Object, *Error) {
	args := []Object{}

	for {
		car := cons.Car
		val := Eval(car, env)
		if isError(val) {
			return nil, val.(*Error)
		}
		args = append(args, val)
		if cons.Cdr != nil {
			cons = cons.Cdr.(*Cons)
		} else {
			break
		}
	}

	return args, nil
}

func loadScopedBuiltins() {
	eval := &ScopedBuiltin{
		Fn: func(env *Environment, args ...Object) Object {
			r := NewReader(args[0].(*String).Value)
			return Eval(r.Read(), env)
		},
	}

	scopedBuiltins["EVAL"] = eval
}
