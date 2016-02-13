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

func TestParseIfElseStatement(t *testing.T) {
	root := parse(`
if (x > 153) {
  let city="Paris";
} else {
  let city="Osaka";
}
`)
	statement := root.Children[0]

	if !(root.Name == "statements" && statement.Name == "ifStatement") {
		t.Errorf("expect node to have: ifStatement, but got: \n%v", root.ToXML())
	}

	found := false
	for _, node := range statement.Children {
		if node.Name == "keyword" && node.Value == "else" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expect node to have \"else\" keyword\n%v", root.ToXML())
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

func TestParseWhileStatement(t *testing.T) {
	root := parse(`
    while (i > 100) {
      let foo=0;
    }
  `)

	if len(root.Children) == 0 {
		t.Errorf("expect node to have children, but got:\n%v", root.ToXML())
		return
	}

	statement := root.Children[0]

	if statement.Name != "whileStatement" {
		t.Errorf("expect node to have whileStatement, but got:\n%v", root.ToXML())
	}
}

func TestParseDoStatement(t *testing.T) {
	root := parse(`
    do foo(1, 2, 3);
  `)

	if len(root.Children) == 0 {
		t.Errorf("expect node to have children, but got:\n%v", root.ToXML())
		return
	}

	statement := root.Children[0]

	if statement.Name != "doStatement" {
		t.Errorf("expect node to have whileStatement, but got:\n%v", root.ToXML())
	}
}

func TestParseReturnStatement(t *testing.T) {
	root := parse(`return 1 + 2;`)

	if len(root.Children) == 0 {
		t.Errorf("expect node to have children, but got:\n%v", root.ToXML())
		return
	}

	statement := root.Children[0]

	if statement.Name != "returnStatement" {
		t.Errorf("expect node to have whileStatement, but got:\n%v", root.ToXML())
	}
}
