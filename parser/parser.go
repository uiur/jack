package parser

import (
	"fmt"

	"github.com/uiureo/jack/tokenizer"
)

func Parse(tokens []*tokenizer.Token) *Node {
	node, _ := parseClass(tokens)

	return node
}

func parseClass(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && tokens[0].Value == "class") {
		return nil, tokens
	}

	node := &Node{Name: "class", Children: []*Node{}}
	node.AppendToken(tokens[0])

	expect(tokens[1], "identifier", "")
	node.AppendToken(tokens[1])

	expect(tokens[2], "symbol", "{")
	node.AppendToken(tokens[2])

	tokens = tokens[3:]

	for {
		classVarDec, rest := parseClassVarDec(tokens)
		if classVarDec == nil {
			break
		}

		node.AppendChild(classVarDec)
		tokens = rest
	}

	for {
		subroutineDec, rest := parseSubroutineDec(tokens)

		if subroutineDec == nil {
			break
		}

		node.Children = append(node.Children, subroutineDec)
		tokens = rest
	}

	expect(tokens[0], "symbol", "}")
	node.AppendToken(tokens[0])

	return node, tokens[1:]
}

func parseClassVarDec(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && (tokens[0].Value == "static" || tokens[0].Value == "field")) {
		return nil, tokens
	}

	node := &Node{Name: "classVarDec", Children: []*Node{}}
	node.AppendToken(tokens[0])

	if !tokens[1].IsType() {
		panic(fmt.Sprintf("unexpected token `%s`, expecting type token", tokens[1].TokenType))
	}
	node.AppendToken(tokens[1])

	expect(tokens[2], "identifier", "")
	node.AppendToken(tokens[2])

	tokens = tokens[3:]
	for {
		if !(tokens[0].TokenType == "symbol" && tokens[0].Value == ",") {
			break
		}

		node.AppendToken(tokens[0])

		expect(tokens[1], "identifier", "")
		node.AppendToken(tokens[1])

		tokens = tokens[2:]
	}

	expect(tokens[0], "symbol", ";")
	node.AppendToken(tokens[0])

	return node, tokens[1:]
}

func parseSubroutineDec(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && (tokens[0].Value == "constructor" || tokens[0].Value == "function" || tokens[0].Value == "method")) {
		return nil, tokens
	}

	node := &Node{Name: "subroutineDec", Children: []*Node{}}
	node.AppendToken(tokens[0])
	node.AppendToken(tokens[1]) // "void" || identifier

	expect(tokens[2], "identifier", "")
	node.AppendToken(tokens[2])

	expect(tokens[3], "symbol", "(")
	node.AppendToken(tokens[3])

	tokens = tokens[4:]

	parameterList, tokens := parseParameterList(tokens)
	node.Children = append(node.Children, parameterList)

	expect(tokens[0], "symbol", ")")
	node.AppendToken(tokens[0])

	subroutineBody, tokens := parseSubroutineBody(tokens[1:])
	node.AppendChild(subroutineBody)

	return node, tokens
}

func parseSubroutineBody(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "symbol" && tokens[0].Value == "{") {
		return nil, tokens
	}

	node := &Node{Name: "subroutineBody", Children: []*Node{}}

	node.AppendToken(tokens[0])

	tokens = tokens[1:]

	for {
		varDec, rest := parseVarDec(tokens)
		if varDec == nil {
			break
		}

		node.Children = append(node.Children, varDec)
		tokens = rest
	}

	statements, tokens := ParseStatements(tokens)
	node.Children = append(node.Children, statements)

	expect(tokens[0], "symbol", "}")
	node.AppendToken(tokens[0])

	return node, tokens[1:]
}

func parseParameterList(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	node := &Node{Name: "parameterList", Children: []*Node{}}

	if tokens[0].IsType() {
		node.AppendToken(tokens[0])

		expect(tokens[1], "identifier", "")
		node.AppendToken(tokens[1])

		tokens = tokens[2:]

		for {
			if tokens[0].TokenType == "symbol" && tokens[0].Value == "," {
				node.AppendToken(tokens[0])

				node.AppendToken(tokens[1])

				expect(tokens[2], "identifier", "")
				node.AppendToken(tokens[2])

				tokens = tokens[3:]
			} else {
				break
			}
		}
	}

	return node, tokens
}

func parseVarDec(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && tokens[0].Value == "var") {
		return nil, tokens
	}

	node := &Node{Name: "varDec", Children: []*Node{}}
	node.AppendToken(tokens[0])

	if !tokens[1].IsType() {
		panic("unexpected token `" + tokens[1].Value + "`, expecting type")
	}
	node.AppendToken(tokens[1])

	expect(tokens[2], "identifier", "")
	node.AppendToken(tokens[2])

	tokens = tokens[3:]
	for {
		if !(tokens[0].TokenType == "symbol" && tokens[0].Value == ",") {
			break
		}

		node.AppendToken(tokens[0])
		expect(tokens[1], "identifier", "")
		node.AppendToken(tokens[1])

		tokens = tokens[2:]
	}

	expect(tokens[0], "symbol", ";")
	node.AppendToken(tokens[0])

	return node, tokens[1:]
}

func ParseStatements(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
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

	if node, rest := parseReturnStatement(tokens); node != nil {
		return node, rest
	}

	return nil, tokens
}

func expect(token *tokenizer.Token, tokenType, value string) {
	if len(value) == 0 {
		if token.TokenType != tokenType {
			panic("unexpected token `" + token.TokenType + "`, expecting `" + tokenType + "`")
		}
	} else {
		if !(token.TokenType == tokenType && token.Value == value) {
			panic("unexpected token `" + token.Value + "`, expecting `" + value + "`")
		}
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

	statements, rest := ParseStatements(rest[2:])
	node.Children = append(node.Children, statements)

	expect(rest[0], "symbol", "}")
	node.AppendToken(rest[0]) // }

	rest = rest[1:]

	if len(rest) > 0 && rest[0].TokenType == "keyword" && rest[0].Value == "else" {
		node.AppendToken(rest[0])

		expect(rest[1], "symbol", "{")
		node.AppendToken(rest[1])

		statements, rest := ParseStatements(rest[2:])
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

	expect(tokens[1], "identifier", "")
	node.AppendToken(tokens[1])

	tokens = tokens[2:]
	if tokens[0].TokenType == "symbol" && tokens[0].Value == "[" {
		node.AppendToken(tokens[0])

		expression, rest := parseExpression(tokens[1:])
		node.AppendChild(expression)

		expect(rest[0], "symbol", "]")
		node.AppendToken(rest[0])

		tokens = rest[1:]
	}

	expect(tokens[0], "symbol", "=")
	node.AppendToken(tokens[0])
	expression, rest := parseExpression(tokens[1:])
	node.Children = append(node.Children, expression)

	expect(rest[0], "symbol", ";")
	node.AppendToken(rest[0])

	return node, rest[1:]
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
	node.AppendChild(expression)

	expect(rest[0], "symbol", ")")
	node.AppendToken(rest[0])

	expect(rest[1], "symbol", "{")
	node.AppendToken(rest[1])

	statements, rest := ParseStatements(rest[2:])
	node.AppendChild(statements)

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

	subroutineCallNodes, rest := parseSubroutineCall(tokens[1:])
	for _, n := range subroutineCallNodes {
		node.AppendChild(n)
	}

	expect(rest[0], "symbol", ";")
	node.AppendToken(rest[0])

	return node, rest[1:]
}

func parseReturnStatement(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	if !(tokens[0].TokenType == "keyword" && tokens[0].Value == "return") {
		return nil, tokens
	}

	node := &Node{Name: "returnStatement", Children: []*Node{}}
	node.AppendToken(tokens[0])

	expression, tokens := parseExpression(tokens[1:])
	if expression != nil {
		node.Children = append(node.Children, expression)
	}
	expect(tokens[0], "symbol", ";")
	node.AppendToken(tokens[0])

	return node, tokens[1:]
}

func parseExpression(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	termNode, restTokens := parseTerm(tokens)
	if termNode == nil {
		return nil, tokens
	}

	node := &Node{Name: "expression"}
	node.Children = []*Node{termNode}

	for {
		if !(len(restTokens) > 0 && restTokens[0].IsOp()) {
			break
		}

		node.AppendToken(restTokens[0])
		termNode, rest := parseTerm(restTokens[1:])
		node.Children = append(node.Children, termNode)

		restTokens = rest
	}

	return node, restTokens
}

func parseExpressionList(tokens []*tokenizer.Token) (*Node, []*tokenizer.Token) {
	node := &Node{Name: "expressionList", Children: []*Node{}}

	if tokens[0].TokenType == "symbol" && tokens[0].Value == ")" {
		return node, tokens
	}

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
	case "stringConstant", "integerConstant":
		node := &Node{Name: "term", Children: []*Node{tokenToNode(tokens[0])}}
		return node, tokens[1:]

	case "keyword":
		node := &Node{Name: "term", Children: []*Node{}}
		node.AppendToken(tokens[0])

		return node, tokens[1:]
	}

	subroutineCallNodes, tokens := parseSubroutineCall(tokens)

	if len(subroutineCallNodes) > 0 {
		node := &Node{Name: "term", Children: subroutineCallNodes}
		return node, tokens
	}

	// varName | varName[expression]
	if tokens[0].TokenType == "identifier" {
		node := &Node{Name: "term", Children: []*Node{tokenToNode(tokens[0])}}
		tokens = tokens[1:]
		if len(tokens) > 0 && tokens[0].TokenType == "symbol" && tokens[0].Value == "[" {
			node.AppendToken(tokens[0])

			expression, rest := parseExpression(tokens[1:])
			node.AppendChild(expression)

			expect(rest[0], "symbol", "]")
			node.AppendToken(rest[0])

			tokens = rest[1:]
		}

		return node, tokens
	}

	// ( expression )
	if tokens[0].TokenType == "symbol" && tokens[0].Value == "(" {
		node := &Node{Name: "term", Children: []*Node{}}
		node.AppendToken(tokens[0])

		expression, tokens := parseExpression(tokens[1:])
		node.AppendChild(expression)

		expect(tokens[0], "symbol", ")")
		node.AppendToken(tokens[0])

		tokens = tokens[1:]

		return node, tokens
	}

	// unaryOp term
	if tokens[0].IsUnaryOp() {
		node := &Node{Name: "term", Children: []*Node{}}
		node.AppendToken(tokens[0])

		term, tokens := parseTerm(tokens[1:])
		node.AppendChild(term)

		return node, tokens
	}

	return nil, tokens
}

func parseSubroutineCall(tokens []*tokenizer.Token) ([]*Node, []*tokenizer.Token) {
	if tokens[0].TokenType != "identifier" {
		return []*Node{}, tokens
	}

	if len(tokens) <= 2 {
		return []*Node{}, tokens
	}

	if !((tokens[1].TokenType == "symbol" && tokens[1].Value == "(") || (tokens[1].TokenType == "symbol" && tokens[1].Value == ".")) {
		return []*Node{}, tokens
	}

	node := &Node{Name: "subroutineCall", Children: []*Node{}}
	node.AppendToken(tokens[0])

	tokens = tokens[1:]

	if tokens[0].TokenType == "symbol" && tokens[0].Value == "." {
		node.AppendToken(tokens[0])
		expect(tokens[1], "identifier", "") // subroutineName
		node.AppendToken(tokens[1])
		tokens = tokens[2:]
	}

	expect(tokens[0], "symbol", "(")
	node.AppendToken(tokens[0])

	expression, rest := parseExpressionList(tokens[1:])
	if expression != nil {
		node.Children = append(node.Children, expression)
	}

	expect(rest[0], "symbol", ")")
	node.AppendToken(rest[0])

	return node.Children, rest[1:]
}

func tokenToNode(token *tokenizer.Token) *Node {
	return &Node{Name: token.TokenType, Value: token.Value}
}
