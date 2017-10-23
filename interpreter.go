package main

import (
	"errors"
	"fmt"
)

// Load setups the inital environment and returns it
func Load() *Environment {
	loadScopedBuiltins()
	return NewEnvironment()
}

// Eval an object
func Eval(obj Object, env *Environment) Object {
	switch node := obj.(type) {
	case *Boolean, *Char, *String, *Error, *Integer, *Float, *Vector, *Data:
		return obj
	case *Lambda:
		node.Env = env
		return node
	case *Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val
		}

		return newError(fmt.Sprintf("Unkown identifier %s", node.Value))
	case *Pair:
		car := node.Car
		switch carType := car.(type) {
		case *Lambda:
			if node.Cdr != nil {
				args, err := evalArgs(node.Cdr.(*Pair), env)
				if err != nil {
					return err
				}

				if len(args) == len(carType.Parameters) {
					return applyFunction(carType, "#<procedure>", args)
				}

				return newError("arguments do not match")
			}

			return applyFunction(carType, "#<procedure>", []Object{})
		case *Identifier:
			if builtin, ok := builtins[carType.Value]; ok {
				args, err := evalArgs(node.Cdr.(*Pair), env)
				if err != nil {
					return err
				}

				return builtin.Fn(args...)
			}

			if scopedBuiltin, ok := scopedBuiltins[carType.Value]; ok {
				if node.Cdr != nil {
					args, err := evalArgs(node.Cdr.(*Pair), env)
					if err != nil {
						return err
					}

					return scopedBuiltin.Fn(env, args...)
				}

				return scopedBuiltin.Fn(env, []Object{}...)
			}

			// Check the ENV for a Lambda
			if val, ok := env.Get(carType.Value); ok {
				if lambda, ok := val.(*Lambda); ok {
					var params []Object

					if node.Cdr != nil {
						args, err := evalArgs(node.Cdr.(*Pair), env)
						if err != nil {
							return err
						}

						params = args
					}

					return applyFunction(lambda, carType.Value, params)
				}
			}

			if carType.Value == "DEFINE" {
				ident := node.Cdr.(*Pair).Car.(*Identifier)
				value := node.Cdr.(*Pair).Cdr.(*Pair).Car
				env.Set(ident.Value, value)
				return nil
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

func evalArgs(pair *Pair, env *Environment) ([]Object, *Error) {
	args := []Object{}

	for {
		car := pair.Car
		val := Eval(car, env)
		if isError(val) {
			return nil, val.(*Error)
		}
		args = append(args, val)
		if pair.Cdr != nil {
			pair = pair.Cdr.(*Pair)
		} else {
			break
		}
	}

	return args, nil
}
