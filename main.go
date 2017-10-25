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
	// env := Load()

	for {
		fmt.Print(">> ")
		text, _ := replReader.ReadString('\n')
		cleanText := strings.Trim(text, "\n")
		if cleanText == ".exit" {
			break
		}

		val := Eval(cleanText)
		if val != nil {
			fmt.Println(val.String())
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
