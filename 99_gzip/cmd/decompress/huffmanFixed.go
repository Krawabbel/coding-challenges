package main

var (
	fixedHuffmanLitValCodes *huffmanNode
	fixedHuffmanDistCodes   *huffmanNode
)

func (d *decompressor) parseFixedHuffmanCodes() error {

	debugln(" -> fixed huffman compression")

	return d.parseHuffmanCodes(fixedHuffmanLitValCodes, fixedHuffmanDistCodes)

}

func initFixedHuffmanLitValCodes() (err error) {
	N := 288
	fixedHuffmanLitValCodeLengths := make([]int, N)
	for i := 0; i <= 143; i++ {
		fixedHuffmanLitValCodeLengths[i] = 8
	}
	for i := 144; i <= 255; i++ {
		fixedHuffmanLitValCodeLengths[i] = 9
	}
	for i := 256; i <= 279; i++ {
		fixedHuffmanLitValCodeLengths[i] = 7
	}
	for i := 280; i <= 287; i++ {
		fixedHuffmanLitValCodeLengths[i] = 8
	}

	fixedHuffmanLitValCodes, err = generateTreeNumbered(fixedHuffmanLitValCodeLengths)

	return err
}

func initFixedHuffmanDistCodes() (err error) {
	N := 32
	fixedHuffmanDistCodeLengths := make([]int, N)
	for i := range fixedHuffmanDistCodeLengths {
		fixedHuffmanDistCodeLengths[i] = 5
	}

	fixedHuffmanDistCodes, err = generateTreeNumbered(fixedHuffmanDistCodeLengths)

	return err
}
