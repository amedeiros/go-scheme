package evaluator

// import (
// 	"fmt"

// 	"github.com/amedeiros/go-scheme/object"
// 	"github.com/amedeiros/go-scheme/parser"
// )

// // Eval an Ast
// func Eval(node parser.Ast, env *object.Environment) object.Object {
// 	switch node := node.(type) {
// 	case *parser.Program:
// 		return evalProgram(node, env)
// 	case *parser.LambdaLiteral:
// 		return &object.Lambda{Parameters: node.Paramemeters, Body: node.Body, Env: env}
// 	case *parser.ProcedureCall:
// 		function := Eval(node.Function, env)
// 		args := evalExpressions(node.Arguments, env)
// 		fmt.Println(fmt.Sprintf("ARGS: %#v", args))
// 		return applyLambda(function, args)
// 	case *parser.IntegerLiteral:
// 		return &object.Integer{Value: node.Value}
// 	case *parser.Identifier:
// 		return evalIdentifier(node, env)
// 	default:
// 		panic("No clue man!")
// 	}
// }

// func evalProgram(node *parser.Program, env *object.Environment) object.Object {
// 	var obj object.Object
// 	for _, statement := range node.Expressions {
// 		obj = Eval(statement, env)
// 	}

// 	return obj
// }

// func evalExpressions(expressions []parser.Ast, env *object.Environment) []object.Object {
// 	var objs []object.Object

// 	for _, expression := range expressions {
// 		objs = append(objs, Eval(expression, env))
// 	}

// 	return objs
// }

// func applyLambda(fn object.Object, args []object.Object) object.Object {
// 	switch fn := fn.(type) {

// 	case *object.Lambda:
// 		extendedEnv := extendLambdaEnv(fn, args)
// 		return Eval(fn.Body, extendedEnv)
// 	case *object.Builtin:
// 		return fn.Fn(args...)
// 	default:
// 		return nil
// 		// return newError("not a function: %s", fn.Type())
// 	}
// }

// func extendLambdaEnv(fn *object.Lambda, args []object.Object) *object.Environment {
// 	env := object.NewEnclosedEnvironment(fn.Env)

// 	for paramIdx, param := range fn.Parameters {
// 		env.Set(param.Value, args[paramIdx])
// 	}

// 	return env
// }

// func evalIdentifier(node *parser.Identifier, env *object.Environment) object.Object {
// 	if val, ok := env.Get(node.Value); ok {
// 		return val
// 	}

// 	if builtin, ok := builtins[node.Value]; ok {
// 		return builtin
// 	}

// 	// return newError("identifier not found: " + node.Value)
// 	return nil
// }
