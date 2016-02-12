package parser

import (
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
