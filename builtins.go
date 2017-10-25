package main

import "fmt"
import "reflect"

var scopedBuiltins = map[string]*ScopedBuiltin{}

var builtins = map[string]*Builtin{
	"+": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						obj.Value += intArg.Value
					} else {
						return newError("Expecting an Integer")
					}
				}

				return obj
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						obj.Value += floatArg.Value
					} else {
						return newError("Expecting a Float")
					}
				}

				return obj
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	"-": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						obj.Value -= intArg.Value
					} else {
						return newError("Expecting an Integer")
					}
				}

				return obj
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						obj.Value -= floatArg.Value
					} else {
						return newError("Expecting a Float")
					}
				}

				return obj
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	"*": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						obj.Value *= intArg.Value
					} else {
						return newError("Expecting an Integer")
					}
				}

				return obj
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						obj.Value *= floatArg.Value
					} else {
						return newError("Expecting a Float")
					}
				}

				return obj
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	"/": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						obj.Value /= intArg.Value
					} else {
						return newError("Expecting an Integer")
					}
				}

				return obj
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						obj.Value /= floatArg.Value
					} else {
						return newError("Expecting a Float")
					}
				}

				return obj
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	"<": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						if !(obj.Value < intArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting an Integer")
					}
				}

				return TRUE
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						if !(obj.Value < floatArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting a Float")
					}
				}

				return TRUE
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	"<=": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						if !(obj.Value <= intArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting an Integer")
					}
				}

				return TRUE
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						if !(obj.Value <= floatArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting a Float")
					}
				}

				return TRUE
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	">": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						if !(obj.Value > intArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting an Integer")
					}
				}

				return TRUE
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						if !(obj.Value > floatArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting a Float")
					}
				}

				return TRUE
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	">=": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						if !(obj.Value >= intArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting an Integer")
					}
				}

				return TRUE
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						if !(obj.Value >= floatArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting a Float")
					}
				}

				return TRUE
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	"=": &Builtin{
		Fn: func(args ...Object) Object {
			switch obj := args[0].(type) {
			case *Integer:
				for _, rightSide := range args[1:len(args)] {
					if intArg, ok := rightSide.(*Integer); ok {
						if !(obj.Value == intArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting an Integer")
					}
				}

				return TRUE
			case *Float:
				for _, rightSide := range args[1:len(args)] {
					if floatArg, ok := rightSide.(*Float); ok {
						if !(obj.Value == floatArg.Value) {
							return FALSE
						}
					} else {
						return newError("Expecting a Float")
					}
				}

				return TRUE
			default:
				return errorObject(fmt.Errorf("Unexpected %s expecting one of Integer or Float", obj.Inspect()))
			}
		},
	},
	"DISPLAY": &Builtin{
		Fn: func(args ...Object) Object {
			fmt.Print(args[0].Inspect())
			return nil
		},
	},
	"CALL-METHOD": &Builtin{
		Fn: func(args ...Object) Object {
			if methodName, ok := args[0].(*String); ok {
				obj := args[1]
				values := reflect.ValueOf(obj).MethodByName(methodName.Value).Call([]reflect.Value{})
				return &String{Value: values[0].String()}
			}

			return newError("Expecting a string as the first argument")
		},
	},
	"CALL-FIELD": &Builtin{
		Fn: func(args ...Object) Object {
			if fieldName, ok := args[0].(*String); ok {
				obj := args[1]

				value := reflect.ValueOf(obj).Elem()
				field := value.FieldByName(fieldName.Value).Interface().(Object)
				return field
			}

			return newError("Expecting a string as the first argument")
		},
	},
	"BOOLEAN?": &Builtin{
		Fn: func(args ...Object) Object {
			switch args[0].(type) {
			case *Boolean:
				return TRUE
			default:
				return FALSE
			}
		},
	},
	"PAIR?": &Builtin{
		Fn: func(args ...Object) Object {
			switch kind := args[0].(type) {
			case *Pair:
				if kind.Car != nil {
					return TRUE
				}

				return FALSE
			default:
				return FALSE
			}
		},
	},
	"SYMBOL?": &Builtin{
		Fn: func(args ...Object) Object {
			switch args[0].(type) {
			case *Identifier:
				return TRUE
			default:
				return FALSE
			}
		},
	},
}

func loadScopedBuiltins() {
	eval := &ScopedBuiltin{
		Fn: func(env *Environment, args ...Object) Object {
			switch kind := args[0].(type) {
			case *String:
				return kind
			default:
				ap(kind)
				return Eval(kind, env)
			}
		},
	}

	env := &ScopedBuiltin{
		Fn: func(env *Environment, args ...Object) Object {
			return env
		},
	}

	scopedBuiltins["EVAL"] = eval
	scopedBuiltins["ENV"] = env
}
