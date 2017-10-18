package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/amedeiros/go-scheme/evaluator"
	"github.com/amedeiros/go-scheme/lexer"
	"github.com/amedeiros/go-scheme/object"
	"github.com/amedeiros/go-scheme/parser"
)

func main() {
	fmt.Println("Go Schemeing 1.0.0")
	fmt.Println("Type .exit to exit")
	reader := bufio.NewReader(os.Stdin)
	env := object.NewEnvironment() // Global ENV

	for {
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		cleanText := strings.Trim(text, "\n")
		if cleanText == ".exit" {
			break
		}

		lex := lexer.NewLexer(cleanText)
		parse := parser.NewParser(lex)
		program := parse.ParseProgram()
		result := evaluator.Eval(program, env)
		fmt.Println(result.Inspect())
		// fmt.Println(program.Inspect())
	}
}
