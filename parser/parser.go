package parser

import "github.com/uiureo/jack/tokenizer"

func Parse(tokens []*tokenizer.Token) *Node {
	firstToken := tokens[0]

	if firstToken.TokenType == "keyword" && firstToken.Value == "let" {
		node, _ := parseLetStatement(tokens)
		return node
	}

	return nil
}

func parseLetStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	node := &Node{Name: "letStatement", Children: []*Node{}}
	node.AppendToken(tokens[0])
	node.AppendToken(tokens[1])

	if tokens[2].TokenType == "symbol" && tokens[2].Value == "=" {
		node.AppendToken(tokens[2])
		expression, rest := parseExpression(tokens[3:])
		node.Children = append(node.Children, expression)
		node.AppendToken(rest[0])

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
