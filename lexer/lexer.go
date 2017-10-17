package lexer

import (
	"fmt"
)

type Lexer struct {
	input                               string
	row, column, position, readPosition int
	currentChar                         byte
}

func NewLexer(input string) *Lexer {
	lex := &Lexer{input: input, row: 0, column: -1}
	lex.consume() // prime
	return lex
}

func (lexer *Lexer) NextToken() Token {
	var tok Token
	lexer.consumeWS()

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
	case '"':
		row := lexer.row
		column := lexer.column
		value := lexer.consumeString()
		return Token{Column: column, Row: row, Type: STRING, Literal: value}
	case 0:
		tok = Token{Column: lexer.column, Row: lexer.row, Type: EOF, Literal: "EOF"}
	default:
		if isDigit(lexer.currentChar) {
			row := lexer.row
			column := lexer.column
			literal := lexer.readNumber()
			return Token{Column: column, Row: row, Type: DIGIT, Literal: literal}
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

// TODO: This does not consider escape sequences. This is an easy fix just being lazy.
func (lexer *Lexer) consumeString() string {
	position := lexer.readPosition

	for {
		lexer.consume()
		if lexer.currentChar == '"' || lexer.currentChar == 0 {
			break
		}
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position

	for isDigit(lexer.currentChar) {
		lexer.consume()
	}

	return lexer.input[position:lexer.position]
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
	if lexer.readPosition >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		lexer.currentChar = lexer.input[lexer.readPosition]

		if lexer.currentChar == '\n' {
			lexer.row++
			lexer.column = -1
		} else {
			lexer.column++
		}
	}

	lexer.position = lexer.readPosition
	lexer.readPosition++
}
