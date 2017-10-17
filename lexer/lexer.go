package lexer

import (
	"fmt"
)

type Lexer struct {
	input           string
	row, column, sp int
	currentChar     byte
}

func NewLexer(input string) *Lexer {
	lex := &Lexer{input: input, row: 0, column: -1, sp: 0}
	lex.consume() // prime
	return lex
}

func (lexer *Lexer) NextToken() Token {
	lexer.consumeWS()
	var tok Token

	switch lexer.currentChar {
	case '(':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: LPAREN, Literal: "("}
	case ')':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: RPAREN, Literal: ")"}
	case '+':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: IDENT, Literal: "+"}
	case '-':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: IDENT, Literal: "-"}
	case '*':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: IDENT, Literal: "*"}
	case '/':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: IDENT, Literal: "/"}
	case 0:
		tok = Token{Column: lexer.column, Row: lexer.row, Type: EOF, Literal: "EOF"}
	default:
		if isDigit(lexer.currentChar) {
			row := lexer.row
			column := lexer.column
			literal := lexer.readNumber()
			tok = Token{Column: column, Row: row, Type: DIGIT, Literal: literal}
		} else {
			msg, _ := fmt.Printf("Unkown character %c at %d:%d", lexer.currentChar, lexer.row, lexer.column)
			panic(msg)
		}
	}

	lexer.consume()

	return tok
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (lexer *Lexer) readNumber() string {
	position := lexer.sp - 1

	for isDigit(lexer.currentChar) {
		lexer.consume()
	}

	return lexer.input[position:lexer.sp]
}

func (lexer *Lexer) consumeWS() {
	for isWhiteSpace(lexer.currentChar) {
		lexer.consume()
	}
}

func isWhiteSpace(char byte) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == '\r'
}

func (lexer *Lexer) consume() {
	if lexer.sp < len(lexer.input) {
		lexer.currentChar = byte(lexer.input[lexer.sp])
		lexer.sp += 1

		if lexer.currentChar == '\n' {
			lexer.row += 1
			lexer.column = -1
		} else {
			lexer.column += 1
		}
	} else {
		lexer.currentChar = 0
	}
}
