package main

import (
	"errors"
	"fmt"
	"io/ioutil"
)

// Load setups the inital environment and returns it
func Load() *Environment {
	env := NewEnvironment()
	loadScopedBuiltins()

	b, err := ioutil.ReadFile("./lib/builtins.scm") // just pass the file name
	if err != nil {
		panic(err)
	}

	r := NewReader(string(b))
	program := r.ReadAll()
	for _, prog := range program {
		Eval(prog, env)
	}

	return env
}

// Eval an object
func Eval(obj Object, env *Environment) Object {
	switch node := obj.(type) {
	case *Boolean, *Char, *String, *Error, *Integer, *Float, *Vector:
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
			carType.Env = env

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
			if carType.Value == "DEFINE" {
				ident := node.Cdr.(*Pair).Car.(*Identifier)
				value := node.Cdr.(*Pair).Cdr.(*Pair).Car
				env.Set(ident.Value, value)
				return nil
			}

			if carType.Value == "QUOTE" {
				return node.Cdr
			}

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

					lambda.Env = env
					return applyFunction(lambda, carType.Value, params)
				}
			}

			return newError(fmt.Sprintf("Unkown proc %s", carType.Value))
		default:
			// Empty Pair
			if carType == nil {
				return node
			}

			obj = Eval(node.Cdr, env)
			ap(obj)
			return obj

			// for {
			// 	obj = Eval(node.Cdr, env)

			// 	if node.Cdr != nil {
			// 		node.Cdr = node.Cdr.(*Pair).Car
			// 		ap(node.Cdr.Inspect())
			// 	} else {
			// 		break
			// 	}
			// }

			// return obj
		}
	}

	panic("You just found a bug or an unimplemented feature congrats!")
}

func applyFunction(lambda *Lambda, name string, args []Object) Object {
	extendedEnv := extendFunctionEnv(lambda, name, args)

	// ap(lambda.Body.Inspect())
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
