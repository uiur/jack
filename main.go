package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/uiureo/jack/parser"
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
	tree := parser.Parse(tokens)

	fmt.Print(tree.ToXML())
}

// class Square {
// 	field int x, y;
// 	field int size;
//
// 	constructor Square new(int Ax, int Ay) {
// 		var boolean a;
//
// 		let x = Ax;
// 		let y = Ay;
//
// 		return this;
// 	 }
// }
// SymbolTable
// table.look("a") //=> {"local", "boolean", 0}
// [
// 	{Ax: {"argument", "int", 0}, Ay: {"argument", "int", 1}, a: {"local", "boolean", 2} },
// 	{x: {"field", "int", 0}, y: {"field", "int", 1}, size: {"field", "int", 2}},
// ]

type Symbol struct {
	// Kind: var, argument, static, field, class, subroutine
	SymbolType, Kind string
	Number           int
}

type SymbolTable struct {
	Scopes []map[string]*Symbol
}

// compile(node *Node) string
// compileClass(node *Node, table *SymbolTable) string
// compileSubroutineDec(node *Node, table *SymbolTable) string

func buildSymbolTable(node *parser.Node, base *SymbolTable) *SymbolTable {
	if base == nil {
		base = &SymbolTable{}
	}

	table := &SymbolTable{}

	currentScope := map[string]*Symbol{}

	baseScopes := []map[string]*Symbol{}
	copy(baseScopes, base.Scopes)

	table.Scopes = append([]map[string]*Symbol{currentScope}, baseScopes...)

	switch node.Name {
	case "class":
		for _, node := range node.Children {
			if node.Name == "classVarDec" {
				kind := node.Children[0].Value
				symbolType := node.Children[1].Value

				names := []string{}
				for _, node := range node.Children[2:] {
					if node.Name == "identifier" {
						names = append(names, node.Value)
					}
				}

				for _, name := range names {
					currentScope[name] = &Symbol{SymbolType: symbolType, Kind: kind, Number: len(currentScope)}
				}
			}
		}
	case "subroutineDec":
		var parameterList *parser.Node

		for _, node := range node.Children {
			if node.Name == "parameterList" {
				parameterList = node
				break
			}
		}

		for i, node := range parameterList.Children {
			if node.Name == "keyword" {
				keyword := node
				identifier := parameterList.Children[i+1]
				name := identifier.Value

				currentScope[name] = &Symbol{
					SymbolType: keyword.Value,
					Kind:       "argument",
					Number:     len(currentScope),
				}
			}
		}

		subroutineBody, _ := node.Find(&parser.Node{Name: "subroutineBody"})
		for _, node := range subroutineBody.Children {
			if node.Name == "varDec" {
				symbolType := node.Children[1].Value

				names := []string{}
				for _, node := range node.Children {
					if node.Name == "identifier" {
						names = append(names, node.Value)
					}
				}

				for _, name := range names {
					currentScope[name] = &Symbol{
						SymbolType: symbolType,
						Kind:       "var",
						Number:     len(currentScope),
					}
				}
			}
		}
	}

	return table
}
