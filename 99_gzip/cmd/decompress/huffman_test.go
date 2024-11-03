package main

import (
	"fmt"
	"testing"
)

func (n *huffmanNode) get(code string) (uint64, error) {

	if n == nil {
		return 0, fmt.Errorf("code does not exist")
	}

	if len(code) > 0 {
		switch code[0] {
		case '0':
			return n.left.get(code[1:])
		case '1':
			return n.right.get(code[1:])
		default:
			return 0, fmt.Errorf("invalid code")
		}

	}

	if !n.isLeaf {
		return 0, fmt.Errorf("code does not specify leaf")
	}

	return n.element, nil
}

func TestNilNode(t *testing.T) {
	n := new(huffmanNode)
	n.insertElement("101", 12)
	val, err := n.get("101")
	if err != nil {
		t.Fatal(err)
	}
	if val != 12 {
		t.Fatal(val)
	}

}
