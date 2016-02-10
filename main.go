package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// if (x < 153)
// {let city="Paris";}

// -> [
//   &Token{type: "keyword", value: "if"},
//   &Token{type: "symbol", value: "("},
//   &Token{type: "identifier", value: "x"}
// ]

var keywords = []string{
	"class",
	"constructor",
	"function",
	"method",
	"field",
	"static",
	"var",
	"int",
	"char",
	"boolean",
	"void",
	"let",
	"do",
	"if",
	"else",
	"while",
	"return",
}

var symbols = []string{
	"{",
	"}",
	"(",
	")",
	"[",
	"]",
	".",
	",",
	";",
	"+",
	"-",
	"*",
	"/",
	"&",
	"|",
	"<",
	">",
	"=",
	"~",
}

type Token struct {
	TokenType string
	Value     string
}

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

	tokens := tokenize(string(data))

	for _, token := range tokens {
		fmt.Printf("%s %s\n", token.TokenType, token.Value)
	}
}

func tokenize(source string) []*Token {
	tokenRegexpMap := buildTokenRegexpMap()
	tokenRegexp := regexp.MustCompile(
		strings.Join([]string{
			tokenRegexpMap["keyword"],
			tokenRegexpMap["symbol"],
			tokenRegexpMap["integerConstant"],
			tokenRegexpMap["stringConstant"],
			tokenRegexpMap["identifier"],
		}, "|"),
	)

	tokenValues := tokenRegexp.FindAllString(source, -1)

	tokens := make([]*Token, len(tokenValues))
	for i, tokenValue := range tokenValues {
		tokenType := detectTokenType(tokenValue)
		if tokenType == "stringConstant" {
			tokenValue = strings.Trim(tokenValue, `"`)
		}

		tokens[i] = &Token{TokenType: tokenType, Value: tokenValue}
	}

	return tokens
}

func buildTokenRegexpMap() map[string]string {
	return map[string]string{
		"keyword":         buildRegexpFromList(keywords),
		"symbol":          buildRegexpFromList(symbols),
		"integerConstant": `\d+`,
		"stringConstant":  `"[^"\n]+"`,
		"identifier":      `[a-zA-Z_]\w*`,
	}
}

func detectTokenType(token string) string {
	for tokenType, regexpString := range buildTokenRegexpMap() {
		matched := regexp.MustCompile(regexpString).MatchString(token)
		if matched {
			return tokenType
		}
	}

	return ""
}

func buildRegexpFromList(strs []string) string {
	escaped := make([]string, len(strs))
	for i, str := range strs {
		escaped[i] = regexp.QuoteMeta(str)
	}

	return strings.Join(escaped, "|")
}
