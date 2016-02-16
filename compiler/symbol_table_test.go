package compiler

import "testing"

func TestGet(t *testing.T) {
	table := new(SymbolTable)

	table.Scopes = []map[string]*Symbol{
		{"a": {Number: 0}, "b": {Number: 1}},
		{"c": {Number: 2}, "a": {Number: 3}},
	}

	symbol := table.Get("a")
	if !(symbol != nil && symbol.Number == 0) {
		t.Errorf("table.Get(\"a\") returns %v, want {Number: 0}", symbol)
	}

	symbol = table.Get("c")
	if !(symbol != nil && symbol.Number == 2) {
		t.Errorf("table.Get(\"c\") returns %v, want {Number: 2}", symbol)
	}
}

func TestSet(t *testing.T) {
	table := new(SymbolTable)

	table.Scopes = []map[string]*Symbol{
		{"a": {Kind: "local", Number: 0}, "b": {Kind: "local", Number: 1}},
		{"c": {Kind: "static", Number: 0}, "d": {Kind: "static", Number: 1}},
	}

	table.Set("x", &Symbol{Kind: "static"})
	symbol := table.Get("x")

	if !(symbol != nil && symbol.Number == 0) {
		t.Errorf("%v, want Number: 0", symbol)
	}

	table.Set("y", &Symbol{Kind: "local"})
	symbol = table.Get("y")
	if !(symbol != nil && symbol.Number == 2) {
		t.Errorf("%v, want Number: 2", symbol)
	}
}
