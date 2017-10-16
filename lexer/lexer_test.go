package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := 
`()
+-*/
`
	lex := NewLexer(input)

	tests := []struct {
		expectedType  TokenType
		expectedLiteral string
		expectedRow int
		expectedColumn int
	}{
		{ LPAREN, "(", 0, 0 },
		{ RPAREN, ")", 0, 1 },
		{ ADD, "+", 1, 0 },
		{ SUB, "-", 1, 1 },
		{ MUL, "*", 1, 2 },
		{ DIV, "/", 1, 3 },
	}

	for index, test := range tests {
		tok := lex.NextToken()
		if tok.Type != test.expectedType {
			t.Fatalf("tests[%d] - Wrong TokenType expected=%d, got=%d", index, test.expectedType, tok.Type)
		}

		if tok.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - Wrong TokenLiteral expected=%s, got=%s", index, test.expectedLiteral, tok.Literal)
		}

		if tok.Column != test.expectedColumn {
			t.Fatalf("tests[%d] - Wrong Column expected=%d, got=%d", index, test.expectedColumn, tok.Column)
		}

		if tok.Row != test.expectedRow {
			t.Fatalf("tests[%d] - Wrong Row expected=%d, got=%d", index, test.expectedRow, tok.Row)
		}
	} 
}
