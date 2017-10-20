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
	env := Load()

	for {
		fmt.Print(">> ")
		text, _ := replReader.ReadString('\n')
		cleanText := strings.Trim(text, "\n")
		if cleanText == ".exit" {
			break
		}

		reader := NewReader(cleanText)

		for {
			obj := Eval(reader.Read(), env)

			if isError(obj) {
				if obj.Inspect() != "EOF" {
					fmt.Println(obj.Inspect())
				} else {
					break
				}
			} else {
				fmt.Println(obj.Inspect())
			}
		}
	}
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}
