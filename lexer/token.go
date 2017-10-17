package lexer

type TokenType int

type Token struct {
	Column, Row int
	Type        TokenType
	Literal     string
}

func (tok Token) String() string {
	return tok.Literal
}

const (
	LPAREN TokenType = iota
	RPAREN
	IDENT
	EOF
	DIGIT
)
