package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Go Schemeing 1.0.0")
	fmt.Println("Type .exit to exit")
	replReader := bufio.NewReader(os.Stdin)
	env := Load()

	for {
		fmt.Print(">> ")
		text, _ := replReader.ReadString('\n')
		cleanText := strings.Trim(text, "\n")
		if cleanText == ".exit" {
			break
		}

		reader := NewReader(text)
		program := reader.ReadAll()

		for _, obj := range program {
			obj := Eval(obj, env)
			if obj == nil {
				break
			}

			if isError(obj) {
				fmt.Println(obj.Inspect())
				break
			}

			fmt.Println(obj.Inspect())
		}
	}
}

func isError(obj Object) bool {
	switch obj.(type) {
	case *Error:
		return true
	default:
		return false
	}
}
