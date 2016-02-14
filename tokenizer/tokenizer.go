package tokenizer

import (
	"regexp"
	"strings"
)

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
	"true",
	"false",
	"null",
	"this",
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

func (token *Token) IsOp() bool {
	if token.TokenType != "symbol" {
		return false
	}

	switch token.Value {
	case "+", "-", "*", "/", "&", "|", "<", ">", "=":
		return true
	default:
		return false
	}
}

func (token *Token) IsUnaryOp() bool {
	return token.TokenType == "symbol" && (token.Value == "-" || token.Value == "~")
}

func (token *Token) IsType() bool {
	switch token.TokenType {
	case "keyword":
		return token.Value == "int" || token.Value == "char" || token.Value == "boolean"
	case "identifier":
		return true
	default:
		return false
	}
}

func (token *Token) IsKeywordConstant() bool {
	if token.TokenType != "keyword" {
		return false
	}

	switch token.Value {
	case "true", "false", "null", "this":
		return true
	default:
		return false
	}
}

func Tokenize(source string) []*Token {
	source = removeComment(source)

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

var tokenTypes = []string{
	"keyword",
	"symbol",
	"integerConstant",
	"stringConstant",
	"identifier",
}

func detectTokenType(token string) string {
	regexpMap := buildTokenRegexpMap()
	for _, tokenType := range tokenTypes {
		regexpString := regexpMap[tokenType]
		matched := regexp.MustCompile(`^(` + regexpString + `)$`).MatchString(token)
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

func removeComment(str string) string {
	str = regexp.MustCompile(`(?m)\s*//.+$`).ReplaceAllString(str, "")
	str = regexp.MustCompile(`(?ms)/\*.*?\*/`).ReplaceAllString(str, "")
	return str
}
