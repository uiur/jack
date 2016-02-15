package main

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/uiureo/jack/parser"
	"github.com/uiureo/jack/tokenizer"
)

const toolPath = "/Users/zat/Downloads/nand2tetris/tools/TextComparer.sh"

func TestMain(t *testing.T) {
	files, err := filepath.Glob("fixtures/*.jack")
	if err != nil {
		panic(err)
	}

	for _, jackFile := range files {
		testMainOutput(t, jackFile)
	}
}

func testMainOutput(t *testing.T, jackFile string) {
	xmlFile := regexp.MustCompile(`\.jack$`).ReplaceAllString(jackFile, ".xml")

	name := strings.Split(filepath.Base(jackFile), ".")[0]

	code, err := ioutil.ReadFile(jackFile)
	if err != nil {
		t.Error(err)
		return
	}

	parserOutput := parser.Parse(tokenizer.Tokenize(string(code))).ToXML()

	file, _ := ioutil.TempFile("", "")
	file.Write([]byte(parserOutput))

	output, err := exec.Command(toolPath, file.Name(), xmlFile).CombinedOutput()

	if err != nil {
		t.Errorf("%s: %s %v", name, output, err)
	}
}

func TestBuildSymbolTableFromClass(t *testing.T) {
	node := parser.Parse(tokenizer.Tokenize(`
    class Square {
    	field int x, y;
    	static String s;

    	constructor Square new(int Ax, int Ay) {
    		var boolean a;

    		let x = Ax;
    		let y = Ay;

    		return this;
    	}
    }
  `))

	table := buildSymbolTable(node, nil)
	if len(table.Scopes) != 1 {
		t.Errorf("expect table to have 1 scopes, actual: %d", len(table.Scopes))
		return
	}

	expectedScope := map[string]*Symbol{
		"x": {"int", "field", 0},
		"y": {"int", "field", 1},
		"s": {"String", "static", 2},
	}

	scope := table.Scopes[0]

	testScopeMatch(t, scope, expectedScope)
}

func TestBuildSymbolTableFromSubroutine(t *testing.T) {
	node := parser.Parse(tokenizer.Tokenize(`
    class Square {
    	field int x, y;
    	static String s;

    	constructor Square new(int Ax, int Ay) {
    		var boolean a, b;

    		let x = Ax;
    		let y = Ay;

    		return this;
    	}
    }
  `))

	subroutineDec, _ := node.Find(&parser.Node{Name: "subroutineDec"})

	table := buildSymbolTable(subroutineDec, nil)

	testScopeMatch(t, table.Scopes[0], map[string]*Symbol{
		"Ax": {"int", "argument", 0},
		"Ay": {"int", "argument", 1},
		"a":  {"boolean", "var", 2},
		"b":  {"boolean", "var", 3},
	})
}

func testScopeMatch(t *testing.T, scope, expectedScope map[string]*Symbol) {
	if len(scope) != len(expectedScope) {
		t.Errorf("scope should have %v symbols, but actual: %v", len(expectedScope), len(scope))
		return
	}

	for key, expectedSymbol := range expectedScope {
		symbol := scope[key]
		if !(symbol != nil && symbol.SymbolType == expectedSymbol.SymbolType && symbol.Kind == expectedSymbol.Kind && symbol.Number == expectedSymbol.Number) {
			t.Errorf("expect: %v, actual: %v", expectedScope, scope)
		}
	}
}
