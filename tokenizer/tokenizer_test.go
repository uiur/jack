package tokenizer

import "testing"

func TestTokenize(t *testing.T) {
	expected := [][]string{
		{"if", "keyword"},
		{"(", "symbol"},
		{"x", "identifier"},
		{"<", "symbol"},
		{"153", "integerConstant"},
		{")", "symbol"},
		{"{", "symbol"},
		{"let", "keyword"},
		{"city", "identifier"},
		{"=", "symbol"},
		{"Paris", "stringConstant"},
		{";", "symbol"},
		{"}", "symbol"},
	}

	tokens := Tokenize(`
if (x < 153)
{let city="Paris";}
`)

	if len(tokens) != len(expected) {
		t.Errorf("expect length: %d, got %d", len(tokens), len(expected))
		t.FailNow()
	}

	for i, token := range tokens {
		expectedToken := expected[i]
		if !(token.Value == expectedToken[0] && token.TokenType == expectedToken[1]) {
			t.Errorf("expect {type: %s, value: %s}, got {type: %s, value: %s}", expectedToken[1], expectedToken[0], token.TokenType, token.Value)
		}
	}
}

func TestTokenizeCodeWithComment(t *testing.T) {
	tokens := Tokenize(`
let foo=0; // foo bar
`)

	expected := [][]string{
		{"let", "keyword"},
		{"foo", "identifier"},
		{"=", "symbol"},
		{"0", "integerConstant"},
		{";", "symbol"},
	}

	if len(tokens) != len(expected) {
		t.Errorf("expect length: %d, got %d", len(tokens), len(expected))
		t.FailNow()
	}

	for i, token := range tokens {
		expectedToken := expected[i]
		if !(token.Value == expectedToken[0] && token.TokenType == expectedToken[1]) {
			t.Errorf("expect {type: %s, value: %s}, got {type: %s, value: %s}", expectedToken[1], expectedToken[0], token.TokenType, token.Value)
		}
	}
}
