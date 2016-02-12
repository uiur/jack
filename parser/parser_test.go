package parser

import (
	"testing"

	"github.com/uiureo/jack/tokenizer"
)

func TestParseLetStatement(t *testing.T) {
	tokens := tokenizer.Tokenize(`let city="Paris";`)
	root := Parse(tokens)

	if root.Name != "letStatement" {
		t.Errorf("expect root node: letStatement, got: %v", root.Name)
	}

	if len(root.Children) == 0 {
		t.Errorf("expect root node to have children, but got %v", root)
	}
}

func TestParseIfStatement(t *testing.T) {
	tokens := tokenizer.Tokenize(`
if (x > 153) {
  let city="Paris";
}
`)
	root := Parse(tokens)

	if root.Name != "ifStatement" {
		t.Errorf("expect root node: ifStatement, got: %v", root.ToXML())
	}
}
