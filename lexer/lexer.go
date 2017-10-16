package lexer

import ("fmt")

type Lexer struct {
	input string
	row, column, sp int
	currentChar rune
}

func NewLexer(input string) *Lexer {
	lex := &Lexer{input: input, row: 0, column: -1 , sp: 0 }
	lex.consume() // prime
	return lex
}

func (lexer *Lexer) NextToken() Token {
	lexer.consumeWS()
	var tok Token

	switch lexer.currentChar {
	case '(':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: LPAREN, TokenLiteral: "(" }
	case ')':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: RPAREN, TokenLiteral: ")"}
	case '+':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: ADD, TokenLiteral: "+"}
	case '-':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: SUB, TokenLiteral: "-"}
	case '*':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: MUL, TokenLiteral: "*"}
	case '/':
		tok = Token{Column: lexer.column, Row: lexer.row, Type: DIV, TokenLiteral: "/"}
	case -1:
		tok = Token{Column: lexer.column, Row: lexer.row, Type: EOF, TokenLiteral: "EOF"}
	default:
		msg, _ := fmt.Printf("Unkown character %c at %d:%d", lexer.currentChar, lexer.row, lexer.column)
		panic(msg)
	}

	lexer.consume()

	return tok
}

func (lexer *Lexer) consumeWS() {
	for (isWhiteSpace(lexer.currentChar)) { lexer.consume() }
}

func isWhiteSpace(char rune) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == '\r'
}

func (lexer *Lexer) consume() {
	if lexer.sp < len(lexer.input) {
		lexer.currentChar = rune(lexer.input[lexer.sp])
		lexer.sp += 1

		if (lexer.currentChar == '\n') {
			lexer.row += 1
			lexer.column = -1
		} else {
			lexer.column += 1
		}
	} else {
		lexer.currentChar = -1
	}
}