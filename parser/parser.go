package parser

import "github.com/uiureo/jack/tokenizer"

func Parse(tokens []*tokenizer.Token) *Node {
	node, _ := parseStatement(tokens)

	return node
}

func parseStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if node, rest := parseLetStatement(tokens); node != nil {
		return node, rest
	}

	if node, rest := parseIfStatement(tokens); node != nil {
		return node, rest
	}

	return nil, tokens
}

func parseIfStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && tokens[0].Value == "if") {
		return nil, tokens
	}

	node := &Node{Name: "ifStatement", Children: []*Node{}}
	node.AppendToken(tokens[0]) // if
	node.AppendToken(tokens[1]) // (

	expression, rest := parseExpression(tokens[2:])
	node.Children = append(node.Children, expression)

	node.AppendToken(rest[0]) // )
	node.AppendToken(rest[1]) // {

	statement, rest := parseStatement(rest[2:])
	node.Children = append(node.Children, statement)

	node.AppendToken(rest[0]) // }

	return node, rest[1:]
}

func parseLetStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && tokens[0].Value == "let") {
		return nil, tokens
	}
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

	for {
		if restTokens[0].IsOp() {
			node.AppendToken(restTokens[0])
			termNode, rest := parseTerm(restTokens[1:])
			node.Children = append(node.Children, termNode)

			restTokens = rest
		} else {
			break
		}
	}

	return node, restTokens
}

func parseTerm(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	switch tokens[0].TokenType {
	case "stringConstant", "integerConstant", "identifier":
		node := &Node{Name: "term", Children: []*Node{tokenToNode(tokens[0])}}
		return node, tokens[1:]
	}

	return nil, tokens
}

func tokenToNode(token *tokenizer.Token) *Node {
	return &Node{Name: token.TokenType, Value: token.Value}
}
