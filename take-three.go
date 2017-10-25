package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// GLOBAL environment
var GLOBAL = NewEnvironment()

func Eval(code string) Object {
	tokens := tokens(code)
	tokens = rewrite(tokens)
	obj := parse(&tokens)
	return eval(obj, GLOBAL)
}

func eval(obj Object, env *Environment) Object {
	switch kind := obj.(type) {
	case *String, *Boolean, *Char, *Integer, *Float:
		return obj
	case *Identifier:
		if val, ok := env.Get(kind.Value); ok {
			return val
		}

		return newError(fmt.Sprintf("unkown identifier %s", kind.Value))
	case *Begin:
		var obj Object
		for _, sexpression := range kind.Body {
			obj = eval(sexpression, env)
		}

		return obj
	case *Pair:
		carNode := car(kind)
		switch carType := carNode.(type) {
		case *Identifier:
			switch carType.Value {
			case "begin":
				pair := cdr(kind).(*Pair)
				body := []Object{}

				for {
					sexpression := car(pair)
					body = append(body, sexpression)

					if pair.Cdr == nil {
						break
					}

					pair = cdr(pair).(*Pair)
				}

				return &Begin{Body: body}
			case "define":
				name := car(cdr(kind)).(*Identifier)
				body := car(cdr(cdr(kind)))
				env.Set(name.Value, eval(body, env))
				return nil
			case "lambda":
				args := car(cdr(kind)).(*Pair)
				params := []*Identifier{}
				body := car(cdr(cdr(kind)))
				begin := eval(body, env).(*Begin)
				ap(begin)

				for {
					param := car(args).(*Identifier)
					params = append(params, param)
					if args.Cdr == nil {
						break
					}

					args = cdr(args).(*Pair)
				}

				apMsg("PARAMS: ", params)
				return &Lambda{Parameters: params, Body: begin, Env: env}
			default:
				// Check the env
				if val, ok := env.Get(carType.Value); ok {
					pair := cdr(kind).(*Pair)
					args := []Object{}

					for {
						arg := car(pair)
						args = append(args, eval(arg, env))

						if pair.Cdr == nil {
							break
						}

						pair = cdr(pair).(*Pair)
					}

					apMsg("VAL: ", val)
					apMsg("PAIR: ", pair.Inspect())

					return applyFunction(val.(*Lambda), "lambda", args)
				}

				return newError(fmt.Sprintf("unknown identifier %s", carType.Value))
			}
		}
	default:
		panic(kind.Inspect())
	}
	return obj
}

// Syntatical
func parse(tokens *[]string) Object {
	token := (*tokens)[0]
	*tokens = (*tokens)[1:]

	switch token {
	case "(": // List
		pair := &Pair{Car: parse(tokens)}
		lastPair := pair

		if len(*tokens) < 1 {
			return pair
		}

		for (*tokens)[0] != ")" {
			if (*tokens)[0] == "." {
				*tokens = (*tokens)[1:]
				lastPair.Cdr = parse(tokens)
				break
			}

			obj := parse(tokens)
			lastPair.Cdr = &Pair{Car: obj}
			lastPair = lastPair.Cdr.(*Pair)
		}

		*tokens = (*tokens)[1:]
		return pair
	default: // Atom
		i, err := strconv.ParseInt(token, 0, 64)

		if err != nil {
			f, err := strconv.ParseFloat(token, 64)

			if err != nil {
				ident := strings.ToLower(token)
				if ident == "#f" {
					return FALSE
				} else if ident == "#t" {
					return TRUE
				} else if ident == "#(" {
					vector := &Vector{}
					if len(*tokens) <= 0 {
						return newError("missing closing )")
					}

					for (*tokens)[0] != ")" {
						vector.Value = append(vector.Value, parse(tokens))
					}

					return vector
				} else if len(ident)-1 >= 1 && string(ident[1]) == "\\" {
					return &Char{Value: string(ident[len(ident)-1])}
				}
				return &Identifier{Value: ident}
			}

			return &Float{Value: f}
		}

		return &Integer{Value: i}
	}
}

// Lexical
func tokens(code string) []string {
	strippedComments := stripComments(code)
	tokens := strings.Split(
		strings.Replace(strings.Replace(strippedComments, "(", "( ",
			-1), ")", " )",
			-1), " ")

	return tokens
}

func stripComments(code string) string {
	lines := strings.Split(code, "\n")
	output := make([]string, len(lines))

	for _, line := range lines {
		newLine := ""
		for _, char := range line {
			if char == ';' {
				break
			}

			newLine += string(char)
		}

		output = append(output, newLine)
	}

	return strings.Join(output, "")
}

func rewrite(values []string) []string {
	rewritten := []string{}
	sp := 0

	for sp < len(values) {
		token := values[sp]
		sp++

		// Rewrite define into lambda
		if token == "define" && values[sp] == "(" {
			sp++
			variable := values[sp]
			variables := []string{}
			body := []string{}
			sp++

			for values[sp] != ")" {
				variables = append(variables, values[sp])
				sp++
			}

			sp++

			for values[sp] != ")" {
				body = append(body, values[sp])
				sp++
			}

			sp++

			lambda := fmt.Sprintf("define %s (lambda (%s) (begin %s)))", variable, strings.Join(variables, " "), strings.Join(body, " "))
			lambdaTokens := tokens(lambda)

			for _, v := range lambdaTokens {
				rewritten = append(rewritten, v)
			}
		} else {
			rewritten = append(rewritten, token)
		}
	}

	return rewritten
}

func applyFunction(lambda *Lambda, name string, args []Object) Object {
	extendedEnv := extendFunctionEnv(lambda, name, args)

	return eval(lambda.Body, extendedEnv)
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
