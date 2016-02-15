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

func (table *SymbolTable) Get(name string) *Symbol {
	var symbol *Symbol

	for _, scope := range table.Scopes {
		if value := scope[name]; value != nil {
			symbol = value

			break
		}
	}

	return symbol
}

func (table *SymbolTable) Set(name string, symbol *Symbol) {
	currentScope := table.Scopes[0]

	kindCount := 0
	for _, item := range currentScope {
		if symbol.Kind == item.Kind {
			kindCount++
		}
	}
	symbol.Number = kindCount

	currentScope[name] = symbol
}

// compileClass(node *Node, table *SymbolTable) string
// compileSubroutineDec(node *Node, table *SymbolTable) string

func Compile(node *parser.Node) string {
	result := ""
	table := buildSymbolTable(node, nil)

	className := node.Children[1].Value

	for _, node := range node.Children {
		if node.Name == "subroutineDec" {
			result += compileSubroutineDec(node, table, className)
		}
	}

	return result
}

func pushSymbol(symbol *Symbol) string {
	return fmt.Sprintf("push %s %d\n", symbol.Kind, symbol.Number)
}

func popSymbol(symbol *Symbol) string {
	return fmt.Sprintf("pop %s %d\n", symbol.Kind, symbol.Number)
}

func compileSubroutineDec(node *parser.Node, table *SymbolTable, className string) string {
	result := ""

	table = buildSymbolTable(node, table)
	name := node.Children[2].Value

	localVarCount := 0
	for _, symbol := range table.Scopes[0] {
		if symbol.Kind == "local" {
			localVarCount++
		}
	}

	result += fmt.Sprintf("function %s.%s %d\n", className, name, localVarCount)

	subroutineBody, _ := node.Find(&parser.Node{Name: "subroutineBody"})
	statements, _ := subroutineBody.Find(&parser.Node{Name: "statements"})

	for _, statement := range statements.Children {
		switch statement.Name {
		case "letStatement":
			identifier, _ := statement.Find(&parser.Node{Name: "identifier"})
			symbol := table.Get(identifier.Value)
			if symbol == nil {
				panic(fmt.Sprintf("variable `%v` is not defined", identifier.Value))
			}

			expression, _ := statement.Find(&parser.Node{Name: "expression"})
			result += pushExpression(expression, table)

			result += popSymbol(symbol)

		case "doStatement":
			subroutineCall := &parser.Node{Name: "subroutineCall", Children: statement.Children[1 : len(statement.Children)-1]}
			result += compileSubroutineCall(subroutineCall, table)
			result += "pop temp 0\n"

		case "returnStatement":
			expression, _ := statement.Find(&parser.Node{Name: "expression"})

			if expression != nil {
				result += pushExpression(expression, table)
			} else {
				result += "push constant 0\n"
			}

			result += "return\n"
		}
	}

	return result
}

func pushExpression(expression *parser.Node, table *SymbolTable) string {
	leftTerm, _ := expression.Find(&parser.Node{Name: "term"})

	result := compileTerm(leftTerm, table)

	if len(expression.Children) > 1 {
		operator := expression.Children[1]
		expression.Children = expression.Children[2:]
		result += pushExpression(expression, table)
		result += compileOperator(operator.Value)
	}

	return result
}

func compileOperator(operator string) string {
	switch operator {
	case "+":
		return "add\n"
	case "-":
		return "sub\n"
	default:
		return ""
	}

}

func compileTerm(term *parser.Node, table *SymbolTable) string {
	firstChild := term.Children[0]

	lastChild := term.Children[len(term.Children)-1]

	isSubroutineCall := !(firstChild.Name == "symbol" && firstChild.Value == "(") && (lastChild.Name == "symbol" && lastChild.Value == ")")

	if isSubroutineCall {
		return compileSubroutineCall(term, table)
	}

	switch firstChild.Name {
	case "integerConstant":
		return fmt.Sprintf("push constant %s\n", firstChild.Value)
	case "identifier":
		symbol := table.Get(firstChild.Value)

		return pushSymbol(symbol)
	}

	return ""
}

func compileSubroutineCall(node *parser.Node, table *SymbolTable) string {
	result := ""
	argSize := 0
	_, i := node.Find(&parser.Node{Name: "symbol", Value: "("})

	var functionName string
	if i == 1 {
		functionName = node.Children[0].Value
	} else if i == 3 {
		classOrVarName := node.Children[0].Value

		var className string
		if symbol := table.Get(classOrVarName); symbol != nil {
			className = symbol.SymbolType
			argSize++

			result += fmt.Sprintf("push %v %v\n", symbol.Kind, symbol.Number)
		} else {
			className = classOrVarName
		}

		subroutineName := node.Children[2].Value

		functionName = fmt.Sprintf("%s.%s", className, subroutineName)
	}

	expressionList, _ := node.Find(&parser.Node{Name: "expressionList"})
	argSize += len(expressionList.Children)

	result += fmt.Sprintf("call %s %d\n", functionName, argSize)

	return result
}

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
					table.Set(name, &Symbol{SymbolType: symbolType, Kind: kind})
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

				table.Set(name, &Symbol{SymbolType: keyword.Value, Kind: "argument"})
			}
		}

		subroutineBody, _ := node.Find(&parser.Node{Name: "subroutineBody"})
		for _, node := range subroutineBody.Children {
			if node.Name == "varDec" {
				symbolType := node.Children[1].Value

				names := []string{}
				for _, node := range node.Children[2:] {
					if node.Name == "identifier" {
						names = append(names, node.Value)
					}
				}

				for _, name := range names {
					table.Set(name, &Symbol{
						SymbolType: symbolType,
						Kind:       "local",
					})
				}
			}
		}
	}

	return table
}
