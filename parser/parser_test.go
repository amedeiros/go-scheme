package parser

import (
	"testing"

	"github.com/amedeiros/go-scheme/lexer"
)

func TestParseIdentifier(t *testing.T) {
	input := "*"
	lex := lexer.NewLexer(input)
	parse := NewParser(lex)
	program := parse.ParseProgram()
	if input != program.Inspect() {
		t.Fatalf("Expected * got %s instead", program.Inspect())
	}
}

func TestParseList(t *testing.T) {
	tests := []struct {
		input string
	}{
		{input: "(+ 1 1 1)"},
		{input: "(+ 1 (+ 1 1))"},
		{input: "(+ (+ 1 1) 1)"},
	}

	for _, test := range tests {
		lex := lexer.NewLexer(test.input)
		parse := NewParser(lex)
		program := parse.ParseProgram()
		if test.input != program.Inspect() {
			t.Fatalf("Expected %s got %s instead", test.input, program.Inspect())
		}
	}
}

// func TestParseString(t *testing.T) {
// 	input := `"Apples!"`

// 	lex := lexer.NewLexer(input)
// 	parse := NewParser(lex)
// 	program := parse.ParseProgram()
// 	if "Apples!" != program.Inspect() {
// 		t.Fatalf("Expected %s got %s instead", input, program.Inspect())
// 	}
// }

// func TestParseLambdaLiteral(t *testing.T) {
// 	tests := []struct {
// 		input, expect string
// 	}{
// 		{input: "(lambda (x y) (+ x y))", expect: "(LAMBDA (X Y) (+ X Y))"},
// 		{input: "(lambda () (+ 1 1))", expect: "(LAMBDA () (+ 1 1))"},
// 	}

// 	for _, test := range tests {
// 		lex := lexer.NewLexer(test.input)
// 		parse := NewParser(lex)
// 		program := parse.ParseProgram()
// 		if test.expect != program.Inspect() {
// 			t.Fatalf("Expected %s got %s instead", test.expect, program.Inspect())
// 		}
// 	}
// }

// func TestParseLambdaCall(t *testing.T) {
// 	tests := []struct {
// 		input, expect string
// 	}{
// 		// {input: "((lambda (x y) (+ x y)))", expect: "((LAMBDA (X Y) (+ X Y)))"},
// 		{input: "((lambda () (+ 1 1)))", expect: "((LAMBDA () (+ 1 1)))"},
// 	}

// 	for _, test := range tests {
// 		lex := lexer.NewLexer(test.input)
// 		parse := NewParser(lex)
// 		program := parse.ParseProgram()
// 		if test.expect != program.Inspect() {
// 			t.Fatalf("Expected %s got %s instead", test.expect, program.Inspect())
// 		}
// 	}
// }
