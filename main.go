package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/uiureo/jack/compiler"
	"github.com/uiureo/jack/parser"
	"github.com/uiureo/jack/tokenizer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no files given")
		os.Exit(1)
	}

	parseMode := os.Args[1] == "parse"

	filename := os.Args[1]
	if parseMode {
		filename = os.Args[2]
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	tokens := tokenizer.Tokenize(string(data))
	tree := parser.Parse(tokens)

	if parseMode {
		fmt.Print(tree.ToXML())
	} else {
		fmt.Print(compiler.Compile(tree))
	}
}
