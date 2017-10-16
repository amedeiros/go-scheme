package lexer

type TokenType int

type Token struct {
	Column, Row int
	Type TokenType
	TokenLiteral string
}

const (
	LPAREN TokenType = iota
	RPAREN
	ADD
	SUB
	DIV
	MUL
	EOF
)