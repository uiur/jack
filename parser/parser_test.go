package parser

import (
	"testing"

	"github.com/uiureo/jack/tokenizer"
)

func parse(source string) *Node {
	return Parse(tokenizer.Tokenize(source))
}

func TestParseLetStatement(t *testing.T) {
	tokens := tokenizer.Tokenize(`let city="Paris";`)
	root := Parse(tokens)

	if !(root.Name == "statements" && root.Children[0].Name == "letStatement") {
		t.Errorf("expect node to have: letStatement, but got: \n%v", root.ToXML())
	}
}

func TestParseIfStatement(t *testing.T) {
	tokens := tokenizer.Tokenize(`
if (x > 153) {
  let city="Paris";
}
`)
	root := Parse(tokens)

	if !(root.Name == "statements" && root.Children[0].Name == "ifStatement") {
		t.Errorf("expect node to have: ifStatement, but got: \n%v", root.ToXML())
	}
}

func TestParseStatements(t *testing.T) {
	root := parse(`
let foo="foo";
let bar="bar";
`)

	if !(root.Name == "statements" && len(root.Children) == 2) {
		t.Errorf("expect statements, got: \n%v", root.ToXML())
	}
}
