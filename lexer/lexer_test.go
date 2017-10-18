package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input :=
		`()
+-*/
1000 "Apples!"
lambda
`
	lex := NewLexer(input)

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
		expectedRow     int
		expectedColumn  int
	}{
		{LPAREN, "(", 0, 0},
		{RPAREN, ")", 0, 1},
		{IDENT, "+", 1, 0},
		{IDENT, "-", 1, 1},
		{IDENT, "*", 1, 2},
		{IDENT, "/", 1, 3},
		{DIGIT, "1000", 2, 0},
		{STRING, "Apples!", 2, 5},
		{IDENT, "LAMBDA", 3, 0},
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

func TestLexingPairs(t *testing.T) {
	input := `(+ 1 1 1)`
	lex := NewLexer(input)
	tests := []string{"(", "+", "1", "1", "1", ")"}
	for _, test := range tests {
		tok := lex.NextToken()
		if tok.Literal != test {
			t.Fatalf("Expected %s got %s instead", test, tok.Literal)
		}
	}

}
