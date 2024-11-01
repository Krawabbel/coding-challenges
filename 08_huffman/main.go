package main

import (
	"fmt"
	"os"
	"slices"
)

type node struct {
	weight      uint64
	element     byte
	isLeaf      bool
	depth       int
	left, right *node
}

func sortNodes(m, n *node) int {
	return int(m.weight - n.weight)
}

func encode(text []byte) ([]byte, *node, byte) {

	// count # of occurences
	freq := make([]uint64, 256)
	for _, b := range text {
		freq[b]++
	}

	nodes := []*node{}

	for i, f := range freq {
		if f != 0 {
			n := node{
				weight:  f,
				element: byte(i),
				isLeaf:  true,
				depth:   0,
				left:    nil,
				right:   nil,
			}
			nodes = append(nodes, &n)
		}
	}

	// create Huffman tree
	for len(nodes) > 1 {
		slices.SortFunc(nodes, sortNodes)

		left, right := nodes[0], nodes[1]

		next := &node{
			weight:  left.weight + right.weight,
			element: 0,
			isLeaf:  false,
			depth:   max(left.depth, right.depth) + 1,
			left:    left,
			right:   right,
		}
		nodes = append(nodes[2:], next)
	}

	// create look-up table
	lut := map[byte][]bool{}
	add(lut, []bool{}, nodes[0])

	// replace bytes with Huffman encoding
	data := []byte{}

	buf := []bool{}
	for _, b := range text {
		buf = append(buf, lut[b]...)
		for len(buf) > 8 {
			next := byte(0)
			for i, bit := range buf[0:8] {
				if bit {
					next |= (1 << i)
				}
			}
			data = append(data, next)
			buf = buf[8:]
		}
	}

	last := byte(0)
	for i, bit := range buf {
		if bit {
			last |= (1 << i)
		}
	}
	data = append(data, last)

	rem := byte(8 - len(buf))

	return data, nodes[0], rem
}

func add(lut map[byte][]bool, prefix []bool, n *node) {

	if n.isLeaf {
		code := make([]bool, len(prefix))
		copy(code, prefix)
		lut[n.element] = code
		return
	}

	add(lut, append(prefix, false), n.left)
	add(lut, append(prefix, true), n.right)
}

type bitReader struct {
	code []byte
	buf  []bool
	rem  byte
}

func (r *bitReader) eof() bool {
	return len(r.code) == 0 && len(r.buf) == int(r.rem)
}

func (r *bitReader) pop() bool {

	if r.eof() {
		panic("eof")
	}

	if len(r.buf) == 0 {
		bits := r.code[0]
		r.code = r.code[1:]

		r.buf = make([]bool, 8)
		for i := range r.buf {
			r.buf[i] = (bits & (1 << i)) > 0
		}
	}

	bit := r.buf[0]
	r.buf = r.buf[1:]
	return bit
}

func decode(code []byte, root *node, rem byte) []byte {
	stream := &bitReader{
		code: code,
		rem:  rem,
	}

	text := []byte{}
	for !stream.eof() {
		symbol := nextElement(root, stream)
		text = append(text, symbol)
	}
	return text
}

func nextElement(n *node, stream *bitReader) byte {
	if n.isLeaf {
		return n.element
	}

	bit := stream.pop()

	if bit {
		return nextElement(n.right, stream)
	}

	return nextElement(n.left, stream)
}

func main() {

	path := os.Args[1]

	text, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	code, root, rem := encode(text)

	reconstruct := decode(code, root, rem)

	if err := os.WriteFile(path+".huff", []byte(reconstruct), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("compression rate: %0.2f%%\n", float64(len(code))/float64(len(text))*100)

}
