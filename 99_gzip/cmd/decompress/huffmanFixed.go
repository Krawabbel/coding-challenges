package main

var (
	fixedHuffmanlitCodes  *huffmanNode
	fixedHuffmanDistCodes *huffmanNode
)

func (d *decompressor) parseFixedHuffmanCodes() error {
	debugln(" -> fixed huffman compression")
	return d.parseHuffmanCodes(fixedHuffmanlitCodes, fixedHuffmanDistCodes)
}

func initFixedHuffmanlitCodes() (err error) {
	N := 288
	fixedHuffmanlitCodeLengths := make([]int, N)
	for i := 0; i <= 143; i++ {
		fixedHuffmanlitCodeLengths[i] = 8
	}
	for i := 144; i <= 255; i++ {
		fixedHuffmanlitCodeLengths[i] = 9
	}
	for i := 256; i <= 279; i++ {
		fixedHuffmanlitCodeLengths[i] = 7
	}
	for i := 280; i <= 287; i++ {
		fixedHuffmanlitCodeLengths[i] = 8
	}

	fixedHuffmanlitCodes, err = generateTree(fixedHuffmanlitCodeLengths)

	return err
}

func initFixedHuffmanDistCodes() (err error) {
	N := 32
	fixedHuffmanDistCodeLengths := make([]int, N)
	for i := range fixedHuffmanDistCodeLengths {
		fixedHuffmanDistCodeLengths[i] = 5
	}

	fixedHuffmanDistCodes, err = generateTree(fixedHuffmanDistCodeLengths)

	return err
}
