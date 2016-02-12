package parser

import (
	"testing"

	"github.com/uiureo/jack/tokenizer"
)

func TestParse(t *testing.T) {
	tokens := tokenizer.Tokenize(`let city="Paris";`)
	root := Parse(tokens)

	if root.Name != "letStatement" {
		t.Errorf("expect root node: letStatement, got: %v", root.Name)
	}

	if len(root.Children) == 0 {
		t.Errorf("expect root node to have children, but got %v", root)
	}
}
