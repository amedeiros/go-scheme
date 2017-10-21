package main

import "fmt"

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
	"QUOTE": &Builtin{
		Fn: func(args ...Object) Object {
			return args[0].(*String)
		},
	},
}
