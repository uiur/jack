package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/uiureo/jack/tokenizer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no files given")
		os.Exit(1)
	}

	filename := os.Args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	tokens := tokenizer.Tokenize(string(data))
	tree := parse(tokens)

	fmt.Print(tree.ToXML())
}

// let city="Paris";
// <letStatement>
//   <keyword>let</keyword>
//   <identifier>city</identifier>
//   <symbol>=</symbol>
//   <expression>
//     <term>
//       <stringConstant>Paris</stringConstant>
//     </term>
//   </expression>
//   <symbol>;</symbol>
// </letStatement>

type Node struct {
	Name     string
	Value    string
	Children []*Node
}

func (node *Node) ToXML() string {
	return generateXMLWithIndent(node, 0)
}

func generateXMLWithIndent(node *Node, indent int) string {
	result := ""

	spaces := ""
	for i := 0; i < indent; i++ {
		spaces += " "
	}

	if len(node.Value) > 0 {
		result += fmt.Sprintf(spaces+"<%v>%v</%v>\n", node.Name, node.Value, node.Name)
	} else {
		result += fmt.Sprintf(spaces+"<%v>\n", node.Name)

		for _, n := range node.Children {
			result += generateXMLWithIndent(n, indent+2)
		}

		result += fmt.Sprintf(spaces+"</%v>\n", node.Name)
	}

	return result
}

func (node *Node) appendToken(token *tokenizer.Token) {
	node.Children = append(node.Children, tokenToNode(token))
}

func parse(tokens []*tokenizer.Token) *Node {
	firstToken := tokens[0]

	if firstToken.TokenType == "keyword" && firstToken.Value == "let" {
		node, _ := parseLetStatement(tokens)
		return node
	}

	return nil
}

func parseLetStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	node := &Node{Name: "letStatement", Children: []*Node{}}
	node.appendToken(tokens[0])
	node.appendToken(tokens[1])

	if tokens[2].TokenType == "symbol" && tokens[2].Value == "=" {
		node.appendToken(tokens[2])
		expression, rest := parseExpression(tokens[3:])
		node.Children = append(node.Children, expression)
		node.appendToken(rest[0])

		return node, rest[1:]
	}

	return nil, nil
}

func parseExpression(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	node := &Node{Name: "expression"}
	termNode, restTokens := parseTerm(tokens)
	node.Children = []*Node{termNode}
	return node, restTokens
}

func parseTerm(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if tokens[0].TokenType == "stringConstant" {
		node := &Node{Name: "term", Children: []*Node{tokenToNode(tokens[0])}}
		return node, tokens[1:]
	}
	return nil, nil
}

func tokenToNode(token *tokenizer.Token) *Node {
	return &Node{Name: token.TokenType, Value: token.Value}
}
