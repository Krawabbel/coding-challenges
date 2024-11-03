package main

import "testing"

func TestNilNode(t *testing.T) {
	n := new(huffmanNode)
	n.insert("101", 12)
	val, err := n.get("101")
	if err != nil {
		t.Fatal(err)
	}
	if val != 12 {
		t.Fatal(val)
	}

}
