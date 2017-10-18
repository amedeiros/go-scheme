package evaluator

import (
	"github.com/amedeiros/go-scheme/object"
)

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
