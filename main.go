package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/amedeiros/go-scheme/lexer"
	"github.com/amedeiros/go-scheme/parser"
)

func main() {
	fmt.Println("Go Schemeing 1.0.0")
	fmt.Println("Type .exit to exit")
	reader := bufio.NewReader(os.Stdin)

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
		fmt.Println(program.Inspect())
	}
}
