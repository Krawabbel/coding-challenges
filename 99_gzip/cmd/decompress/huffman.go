package main

import (
	"fmt"
)

type huffmanNode struct {
	element     uint64
	isLeaf      bool
	left, right *huffmanNode
}

func (node *huffmanNode) getElement(s *bitstream) uint64 {
	if node.isLeaf {
		return node.element
	}
	if s.nextBool() {
		return node.right.getElement(s)
	}
	return node.left.getElement(s)
}

func (n *huffmanNode) insertElement(code uint64, clen int, element uint64) error {

	if clen == 0 {
		if n.isLeaf {
			return fmt.Errorf("code already in use")
		}
		n.element = element
		n.isLeaf = true
		return nil
	}

	if n.isLeaf {
		return fmt.Errorf("attempting to insert into leaf node")
	}

	if n.right == nil {
		n.right = new(huffmanNode)
	}

	if n.left == nil {
		n.left = new(huffmanNode)
	}

	if (code & (1 << (clen - 1))) > 0 {
		return n.right.insertElement(code, clen-1, element)
	} else {
		return n.left.insertElement(code, clen-1, element)
	}

}

func generateTreeNumbered(lengths []int) (*huffmanNode, error) {
	elements := make([]uint64, len(lengths))
	for i := range elements {
		elements[i] = uint64(i)
	}
	return generateTree(lengths, elements)
}

func generateTree(tree_len []int, elements []uint64) (*huffmanNode, error) {

	// 1) Count the number of codes for each code length.

	MAX_BITS := 0
	for _, l := range tree_len {
		MAX_BITS = max(MAX_BITS, l)
	}

	bl_count := make([]uint64, MAX_BITS+1)
	for _, l := range tree_len {
		bl_count[l]++
	}

	// 2) Find the numerical value of the smallest code for each code length.

	code := uint64(0)
	bl_count[0] = 0
	next_code := make([]uint64, MAX_BITS+1)
	for bits := 1; bits <= MAX_BITS; bits++ {
		code = (code + bl_count[bits-1]) << 1
		next_code[bits] = code
	}

	// 3) Assign numerical values to all codes, using consecutive values for all codes of the same length with the base values determined at step 2. Codes that are never used (which have a bit length of zero) must not be assigned a	value.

	root := new(huffmanNode)

	for n, l := range tree_len {

		if l != 0 {
			tree_code := next_code[l]
			next_code[l]++

			if err := root.insertElement(tree_code, l, elements[n]); err != nil {
				return nil, err
			}

			str_code := fmt.Sprintf("%0"+fmt.Sprint(l)+"b", tree_code)
			debugf("%3d: symbol: %3d, Length: %2d, Code %s\n", n, elements[n], l, str_code)
			if len(str_code) != l {
				panic("len(str_code) != l")
			}

		}
	}

	return root, nil
}

func (d *decompressor) parseHuffmanCodes(litValCodes, distCodes *huffmanNode) error {

	for {

		litValCode := litValCodes.getElement(d.istream)

		debug(" -> ", litValCode, " -> ")

		if litValCode < 256 {
			literal := byte(litValCode)
			d.push(literal) // literal
		} else if litValCode == 256 {
			debugln("end-of-block")
			break // end-of-block
		} else {
			length, err := d.parseHuffmanLength(litValCode)
			if err != nil {
				return err
			}
			debugln(" -> length:", length)

			distcode := distCodes.getElement(d.istream)

			distance, err := d.parseHuffmanDistance(distcode)
			if err != nil {
				return err
			}
			debugln(" -> distance:", distance)

			debugf(" -> <l:%d, d:%d>\n", length, distance)

			d.repeat(len(d.history)-distance, length)

		}
		debugln()
	}

	return nil

}

var baseHuffmanLengths = []uint64{
	3, 4, 5, 6, 7, 8, 9, 10,
	11, 13, 15, 17,
	19, 23, 27, 31,
	35, 43, 51, 59,
	67, 83, 99, 115,
	131, 163, 195, 227,
	258,
}

func (d *decompressor) parseHuffmanLength(lencode uint64) (int, error) {

	if lencode < 257 || lencode > 285 {
		return 0, corruptFileError("unexpected lencode")
	}

	nExtraBits := 0
	if lencode > 264 && lencode < 285 {
		nExtraBits = int(lencode-265)/4 + 1
	}

	debugln(" -> lencode", lencode)
	debugln(" -> # extra bits", nExtraBits)

	baseLength := baseHuffmanLengths[lencode-257]
	debugln(" -> base length ", baseLength)

	extraBits := d.nextBits(nExtraBits)
	debugln(" -> extra bits ", extraBits)

	return int(baseLength + extraBits), nil
}

var baseHuffmanDistances = []uint64{
	1, 2, 3, 4, 5, 7, 9, 13, 17, 25,
	33, 49, 65, 97, 129, 193, 257, 385, 513, 769,
	1025, 1537, 2049, 3073, 4097, 6145, 8193, 12289, 16385, 24577,
}

func (d *decompressor) parseHuffmanDistance(distcode uint64) (int, error) {

	if distcode > 29 {
		return 0, corruptFileError("unexpected distcode")
	}

	nExtraBits := 0
	if distcode > 1 {
		nExtraBits = int(distcode)/2 - 1
	}

	debugln(" -> distcode:", distcode)
	debugln(" -> # extra bits:", nExtraBits)
	debugln(" -> base dist ", baseHuffmanDistances[distcode])

	extraBits := d.nextBits(nExtraBits)
	debugln(" -> extra bits ", extraBits)

	return int(baseHuffmanDistances[distcode] + extraBits), nil
}
