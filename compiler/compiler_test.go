package compiler

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uiureo/jack/parser"
	"github.com/uiureo/jack/tokenizer"
)

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
		"s": {"String", "static", 0},
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

	classTable := buildSymbolTable(node, nil)
	table := buildSymbolTable(subroutineDec, classTable)

	testScopeMatch(t, table.Scopes[0], map[string]*Symbol{
		"Ax": {"int", "argument", 0},
		"Ay": {"int", "argument", 1},
		"a":  {"boolean", "local", 0},
		"b":  {"boolean", "local", 1},
	})

	testScopeMatch(t, table.Scopes[1], map[string]*Symbol{
		"x": {"int", "field", 0},
		"y": {"int", "field", 1},
		"s": {"String", "static", 0},
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
			t.Errorf("expect: %v, actual: %v", scopeToString(expectedScope), scopeToString(scope))
		}
	}
}

func scopeToString(scope map[string]*Symbol) string {
	result := []string{}

	for name, symbol := range scope {
		result = append(result, fmt.Sprintf("%v: %v", name, symbol))
	}

	return "{ " + strings.Join(result, ", ") + " }"
}

func TestCompileMain(t *testing.T) {
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

func TestCompileFunctionWithArgument(t *testing.T) {
	result := Compile(parser.Parse(tokenizer.Tokenize(`
    class Number {
      function int plus(int a, int b) {
        var int i;
        let i = a + b;

        return i;
      }
    }`)))

	compare(t, result, `
    function Number.plus 1
    push argument 0
    push argument 1
    add
    pop local 0
    push local 0
    return
  `)
}

func TestCompileSeven(t *testing.T) {
	testCompileFiles(t, "./fixtures/Seven/*.jack")
}

func TestCompileConvertToBin(t *testing.T) {
	testCompileFiles(t, "./fixtures/ConvertToBin/*.jack")
}

func TestCompileSquare(t *testing.T) {
	testCompileFiles(t, "./fixtures/Square/*.jack")
}

func testCompileFiles(t *testing.T, pattern string) {
	jackFiles, _ := filepath.Glob(pattern)

	if len(jackFiles) == 0 {
		t.Error("no files found")
		return
	}

	for _, jackFile := range jackFiles {
		name := strings.Split(filepath.Base(jackFile), ".")[0]
		vmFile := filepath.Dir(jackFile) + "/" + name + ".vm"
		vmData, _ := ioutil.ReadFile(vmFile)

		jackData, _ := ioutil.ReadFile(jackFile)
		compiled := compile(string(jackData))

		compare(t, compiled, string(vmData))
	}
}

func compile(source string) string {
	return Compile(parser.Parse(tokenizer.Tokenize(source)))
}

func compare(t *testing.T, code, expected string) {
	codeLines := splitCode(code)
	expectedCodeLines := splitCode(expected)

	if len(expectedCodeLines) == 0 {
		t.Error("no code specified")
		return
	}

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
