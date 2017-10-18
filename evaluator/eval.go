package evaluator

import (
	"github.com/amedeiros/go-scheme/object"
	"github.com/amedeiros/go-scheme/parser"
)

// Eval an Ast
func Eval(node parser.Ast, env *object.Environment) object.Object {
	switch node.(type) {
	case *parser.Program:
		return evalProgram(node.(*parser.Program), env)
	case *parser.ProcedureCall:
		procCall := node.(*parser.ProcedureCall)
		if builtin, ok := builtins[procCall.Name]; ok {
			var params []object.Object

			for _, param := range procCall.Arguments {
				params = append(params, Eval(param, env))
			}

			return builtin.Fn(params...)
		}
		return nil
	case *parser.IntegerLiteral:
		return &object.Integer{Value: node.(*parser.IntegerLiteral).Value}
	default:
		panic("No clue man!")
	}
}

func evalProgram(node *parser.Program, env *object.Environment) object.Object {
	var obj object.Object
	for _, statement := range node.Expressions {
		obj = Eval(statement, env)
	}

	return obj
}
