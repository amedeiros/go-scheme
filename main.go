package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/amedeiros/go-scheme/object"
)

func main() {
	fmt.Println("Go Schemeing 1.0.0")
	fmt.Println("Type .exit to exit")
	replReader := bufio.NewReader(os.Stdin)
	env := object.NewEnvironment()

	for {
		fmt.Print(">> ")
		text, _ := replReader.ReadString('\n')
		cleanText := strings.Trim(text, "\n")
		if cleanText == ".exit" {
			break
		}

		reader := NewReader(cleanText)
		obj := Eval(reader.Read(), env)

		// for {
		// 	if obj == nil {
		// 		break
		// 	}

		if isError(obj) {
			if obj.Inspect() != "EOF" {
				fmt.Println(obj.Inspect())
			}
			// break
		} else {
			fmt.Println(obj.Inspect())
		}

		// fmt.Println(obj.Inspect())
		// obj = reader.Read()
		// }
	}
}

func isError(obj object.Object) bool {
	switch obj.(type) {
	case *object.Error:
		return true
	default:
		return false
	}
}
