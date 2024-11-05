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
	n.insertElement(0b101, 3, 12)
	val, err := n.get("101")
	if err != nil {
		t.Fatal(err)
	}
	if val != 12 {
		t.Fatal(val)
	}

}

func Test_generateTree(t *testing.T) {
	tree_len := []int{3, 3, 3, 3, 3, 2, 4, 4}
	root, err := generateTreeNumbered(tree_len)
	if err != nil {
		t.Fatal(err)
	}

	codes := []string{"010", "011", "100", "101", "110", "00", "1110", "1111"}

	for i, code := range codes {
		want := uint64(i)
		if have, err := root.get(code); err != nil {
			t.Fatal(err)
		} else if have != want {
			t.Fatalf("fail for '%s': have %d, got %d", string(byte(have)+'A'), have, want)
		}
	}
}
