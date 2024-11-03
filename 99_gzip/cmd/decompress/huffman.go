package main

import (
	"fmt"
)

var (
	fixedHuffmanCodes *huffmanNode
	// fixedHuffmanDistanceCodes *huffmanNode
)

// func initFixedHuffmanDistanceCodes() {
// 	fixedHuffmanDistanceLengths := make([]int, 32)
// 	for i := range fixedHuffmanDistanceLengths {
// 		fixedHuffmanDistanceLengths[i] = 5
// 	}
// 	fixedHuffmanDistanceCodes = genTree(fixedHuffmanDistanceLengths)
// }

func initFixedHuffmanCodes() {
	fixedHuffmanCodeLengths := make([]int, 288)
	for i := 0; i <= 143; i++ {
		fixedHuffmanCodeLengths[i] = 8
	}
	for i := 144; i <= 255; i++ {
		fixedHuffmanCodeLengths[i] = 9
	}
	for i := 256; i <= 279; i++ {
		fixedHuffmanCodeLengths[i] = 7
	}
	for i := 280; i <= 287; i++ {
		fixedHuffmanCodeLengths[i] = 8
	}

	fixedHuffmanCodes = genTree(fixedHuffmanCodeLengths)
}

type huffmanNode struct {
	element     uint64
	isLeaf      bool
	left, right *huffmanNode
}

func (n *huffmanNode) insert(code string, element uint64) error {

	if len(code) > 0 {
		if n.isLeaf {
			return fmt.Errorf("attempting to insert into leaf node")
		}
		switch code[0] {
		case '0':
			if n.left == nil {
				n.left = new(huffmanNode)
			}
			return n.left.insert(code[1:], element)
		case '1':
			if n.right == nil {
				n.right = new(huffmanNode)
			}
			return n.right.insert(code[1:], element)
		default:
			return fmt.Errorf("invalid code")
		}

	}

	if n.element != 0 {
		return fmt.Errorf("code already in use")
	}

	n.element = element
	n.isLeaf = true

	return nil
}

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

func genTree(lengths []int) *huffmanNode {

	// 1) Count the number of codes for each code length.

	maxLength := 0
	for _, l := range lengths {
		maxLength = max(maxLength, l)
	}
	maxLength++

	freqLength := make([]uint64, maxLength)
	for _, l := range lengths {
		freqLength[l]++
	}

	// 2) Find the numerical value of the smallest code for each code length.

	nextCode := make([]uint64, maxLength)
	prevCode := uint64(0)
	freqLength[0] = 0
	for l := 1; l < maxLength; l++ {
		prevCode = (prevCode + freqLength[l-1]) << 1
		nextCode[l] = prevCode
	}

	// 3) Assign numerical values to all codes, using consecutive values for all codes of the same length with the base values determined at step 2. Codes that are never used (which have a bit length of zero) must not be assigned a	value.

	root := new(huffmanNode)

	for n, l := range lengths {
		if l != 0 {

			code := fmt.Sprintf("%0"+fmt.Sprint(l)+"b", nextCode[l])
			nextCode[l]++
			element := uint64(n)

			if err := root.insert(code, element); err != nil {
				panic(err)
			}

			// debugf("symbol: %s Length: %d, Code %"+fmt.Sprint(maxLength)+"s\n", fmt.Sprint(n), l, code)
		}
	}

	return root
}
