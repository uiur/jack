package compiler

import "fmt"

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

func (table *SymbolTable) FindAll(query *Symbol) []*Symbol {
	result := []*Symbol{}

	for _, scope := range table.Scopes {
		for _, symbol := range scope {
			if symbol.Kind == query.Kind {
				result = append(result, symbol)
			}
		}
	}

	return result
}

func (table *SymbolTable) Find(query *Symbol) *Symbol {
	symbols := table.FindAll(query)

	if len(symbols) > 0 {
		return symbols[0]
	}

	return nil
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

func (table *SymbolTable) String() string {
	result := ""
	for _, scope := range table.Scopes {
		result += fmt.Sprint(scope) + "\n"
	}
	return result
}
