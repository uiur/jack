package parser

import (
	"strconv"
	"strings"
	"testing"
)

func TestToXML(t *testing.T) {
	node := &Node{Name: "symbol", Value: "<"}
	actual := strings.TrimSpace(node.ToXML())

	if actual != "<symbol> &lt; </symbol>" {
		t.Errorf("the value of node should be escaped, but got \"%v\"", actual)
	}
}

func TestFindAll(t *testing.T) {
	node := &Node{}
	node.Children = []*Node{
		{Name: "expression", Value: "0"},
		{Name: "symbol", Value: ","},
		{Name: "expression", Value: "1"},
		{Name: "symbol", Value: ","},
		{Name: "expression", Value: "2"},
	}

	expressions := node.FindAll(&Node{Name: "expression"})

	passed := true
	for i, expression := range expressions {
		if expression.Name != "expression" || expression.Value != strconv.Itoa(i) {
			passed = false
			break
		}
	}

	if !passed {
		node := &Node{Name: "returnValue", Children: expressions}
		t.Errorf("`FindAll` should return all `expression`: %v", node.ToXML())
	}
}
