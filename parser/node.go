package parser

import (
	"fmt"
	"html"

	"github.com/uiureo/jack/tokenizer"
)

type Node struct {
	Name     string
	Value    string
	Children []*Node
}

func (node *Node) ToXML() string {
	return generateXMLWithIndent(node, 0)
}

func (node *Node) AppendToken(token *tokenizer.Token) {
	node.Children = append(node.Children, tokenToNode(token))
}

func generateXMLWithIndent(node *Node, indent int) string {
	if node == nil {
		return "nil"
	}
	result := ""

	spaces := ""
	for i := 0; i < indent; i++ {
		spaces += " "
	}

	if len(node.Value) > 0 {
		result += fmt.Sprintf(spaces+"<%v> %v </%v>\n", node.Name, html.EscapeString(node.Value), node.Name)
	} else {
		result += fmt.Sprintf(spaces+"<%v>\n", node.Name)

		for _, n := range node.Children {
			result += generateXMLWithIndent(n, indent+2)
		}

		result += fmt.Sprintf(spaces+"</%v>\n", node.Name)
	}

	return result
}
