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
	printNode(parse(tokens))
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

func (node *Node) appendToken(token *tokenizer.Token) {
	node.Children = append(node.Children, tokenToNode(token))
}

func parse(tokens []*tokenizer.Token) *Node {
	firstToken := tokens[0]

	if firstToken.TokenType == "keyword" && firstToken.Value == "let" {
		return parseLetStatement(tokens)
	}

	return nil
}

func tokenToNode(token *tokenizer.Token) *Node {
	return &Node{Name: token.TokenType, Value: token.Value}
}

func parseLetStatement(tokens []*tokenizer.Token) *Node {
	node := &Node{Name: "letStatement", Children: []*Node{}}
	node.appendToken(tokens[0])
	node.appendToken(tokens[1])

	if tokens[2].TokenType == "symbol" && tokens[2].Value == "=" {
		node.appendToken(tokens[2])
		expression := parseExpression(tokens[3:])
		node.Children = append(node.Children, expression)

		// fmt.Println("<symbol>;</symbol>")
		return node
	}

	return nil
}

func parseExpression(tokens []*tokenizer.Token) *Node {
	node := &Node{Name: "expression"}
	node.Children = []*Node{parseTerm(tokens)}
	return node
}

func parseTerm(tokens []*tokenizer.Token) *Node {
	if tokens[0].TokenType == "stringConstant" {
		return &Node{Name: "term", Children: []*Node{tokenToNode(tokens[0])}}
	}
	return nil
}

func printNode(node *Node) {
	printNodeWithIndent(node, 0)
}

func printNodeWithIndent(node *Node, indent int) {
	spaces := ""
	for i := 0; i < indent; i++ {
		spaces += " "
	}

	if len(node.Value) > 0 {
		fmt.Print(spaces)
		fmt.Printf("<%v>%v</%v>\n", node.Name, node.Value, node.Name)
	} else {
		fmt.Print(spaces)
		fmt.Printf("<%v>\n", node.Name)

		for _, n := range node.Children {
			printNodeWithIndent(n, indent+2)
		}

		fmt.Print(spaces)
		fmt.Printf("</%v>\n", node.Name)
	}
}

func printToken(token *tokenizer.Token) {
	fmt.Printf("<%s>%s</%s>\n", token.TokenType, token.Value, token.TokenType)
}
