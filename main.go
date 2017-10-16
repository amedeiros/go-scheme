package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

func main() {
	fmt.Println("Go Schemeing 1.0.0")
	fmt.Println("Type .exit to exit")
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		cleanText := strings.Trim(text, "\n")
		if (cleanText == ".exit") { break }
		fmt.Println(cleanText)	
	}
}