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
		"a":  {"boolean", "local", 2},
		"b":  {"boolean", "local", 3},
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

func TestCompile(t *testing.T) {
	result := Compile(parser.Parse(tokenizer.Tokenize(`
    class Main {
      function void main() {
        var SquareGame game;

        let game = SquareGame.new();
        do game.run();
        do game.dispose();

        return;
      }
    }`)))

	vmCode := `
    function Main.main 1
      call SquareGame.new 0
      pop local 0
      push local 0
      call SquareGame.run 1
      pop temp 0
      push local 0
      call SquareGame.dispose 1
      pop temp 0
      push constant 0
      return
  `

	compare(t, result, vmCode)
}

func compare(t *testing.T, code, expected string) {
	codeLines := splitCode(code)
	expectedCodeLines := splitCode(expected)

	for i, expectedLine := range expectedCodeLines {
		line := ""
		if i < len(codeLines) {
			line = codeLines[i]
		}

		if line != expectedLine {
			t.Errorf("line %d: `%v`, want `%v`", i+1, line, expectedLine)
			break
		}
	}
}

func splitCode(code string) []string {
	var lines []string

	for _, line := range strings.Split(code, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	return lines
}
