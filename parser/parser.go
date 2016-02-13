package parser

import "github.com/uiureo/jack/tokenizer"

func Parse(tokens []*tokenizer.Token) *Node {
	node, _ := parseStatements(tokens)

	return node
}

func parseStatements(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	node := &Node{Name: "statements", Children: []*Node{}}

	for {
		statement, rest := parseStatement(tokens)
		if statement != nil {
			node.Children = append(node.Children, statement)

			tokens = rest
			if len(tokens) == 0 {
				break
			}
		} else {
			break
		}
	}

	return node, tokens
}

func parseStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if node, rest := parseLetStatement(tokens); node != nil {
		return node, rest
	}

	if node, rest := parseIfStatement(tokens); node != nil {
		return node, rest
	}

	if node, rest := parseWhileStatement(tokens); node != nil {
		return node, rest
	}

	if node, rest := parseDoStatement(tokens); node != nil {
		return node, rest
	}

	return nil, tokens
}

func expect(token *tokenizer.Token, tokenType, value string) {
	if !(token.TokenType == tokenType && token.Value == value) {
		panic("unexpected token `" + token.Value + "`, expecting `" + value + "`")
	}
}

func parseIfStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && tokens[0].Value == "if") {
		return nil, tokens
	}

	node := &Node{Name: "ifStatement", Children: []*Node{}}
	node.AppendToken(tokens[0]) // if

	expect(tokens[1], "symbol", "(")
	node.AppendToken(tokens[1]) // (

	expression, rest := parseExpression(tokens[2:])
	node.Children = append(node.Children, expression)

	expect(rest[0], "symbol", ")")
	node.AppendToken(rest[0]) // )
	expect(rest[1], "symbol", "{")
	node.AppendToken(rest[1]) // {

	statements, rest := parseStatements(rest[2:])
	node.Children = append(node.Children, statements)

	expect(rest[0], "symbol", "}")
	node.AppendToken(rest[0]) // }

	rest = rest[1:]

	if len(rest) > 0 && rest[0].TokenType == "keyword" && rest[0].Value == "else" {
		node.AppendToken(rest[0])

		expect(rest[1], "symbol", "{")
		node.AppendToken(rest[1])

		statements, rest := parseStatements(rest[2:])
		node.Children = append(node.Children, statements)

		expect(rest[0], "symbol", "}")
		node.AppendToken(rest[0])

		rest = rest[1:]
	}

	return node, rest
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

func parseWhileStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && tokens[0].Value == "while") {
		return nil, tokens
	}

	node := &Node{Name: "whileStatement", Children: []*Node{}}
	node.AppendToken(tokens[0]) // while

	expect(tokens[1], "symbol", "(")
	node.AppendToken(tokens[1])

	expression, rest := parseExpression(tokens[2:])
	node.Children = append(node.Children, expression)

	expect(rest[0], "symbol", ")")
	node.AppendToken(rest[0])

	expect(rest[1], "symbol", "{")
	node.AppendToken(rest[1])

	statements, rest := parseStatements(rest[2:])
	node.Children = append(node.Children, statements)

	expect(rest[0], "symbol", "}")
	node.AppendToken(rest[0])

	return node, rest[1:]
}

func parseDoStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && tokens[0].Value == "do") {
		return nil, tokens
	}

	node := &Node{Name: "doStatement", Children: []*Node{}}
	node.AppendToken(tokens[0]) // do

	subroutineCall, rest := parseSubroutineCall(tokens[1:])
	node.Children = append(node.Children, subroutineCall)

	expect(rest[0], "symbol", ";")
	node.AppendToken(rest[0])

	return node, rest[1:]
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

func parseExpressionList(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if tokens[0].TokenType == "symbol" && tokens[0].Value == ")" {
		return nil, tokens
	}

	node := &Node{Name: "expressionList", Children: []*Node{}}

	expression, rest := parseExpression(tokens)
	node.Children = append(node.Children, expression)

	for {
		if rest[0].TokenType == "symbol" && rest[0].Value == "," {
			node.AppendToken(rest[0])
			expression, tokens := parseExpression(rest[1:])
			node.Children = append(node.Children, expression)

			rest = tokens
		} else {
			break
		}
	}

	return node, rest
}

func parseTerm(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	switch tokens[0].TokenType {
	case "stringConstant", "integerConstant", "identifier":
		node := &Node{Name: "term", Children: []*Node{tokenToNode(tokens[0])}}
		return node, tokens[1:]
	}

	return nil, tokens
}

func parseSubroutineCall(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if tokens[0].TokenType == "identifier" {
		node := &Node{Name: "subroutineCall", Children: []*Node{}}
		node.AppendToken(tokens[0])

		expect(tokens[1], "symbol", "(")
		node.AppendToken(tokens[1]) // (

		expression, rest := parseExpressionList(tokens[2:])
		node.Children = append(node.Children, expression)

		expect(rest[0], "symbol", ")")
		node.AppendToken(rest[0])

		return node, rest[1:]
	}

	return nil, tokens
}

func tokenToNode(token *tokenizer.Token) *Node {
	return &Node{Name: token.TokenType, Value: token.Value}
}
