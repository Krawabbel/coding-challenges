package main

func (d *decompressor) parseFixedHuffmanCodes() error {
	return d.parseHuffmanCodes(fixedHuffmanCodes)
}

func (d *decompressor) parseHuffmanCodes(litTree *huffmanNode) error {
	debugln(" -> fixed huffman compression")
	for {

		val := litTree.getElement(d.istream)

		debug(" -> ", val, " -> ")

		if val < 256 {
			d.push(byte(val)) // literal
		} else if val == 256 {
			debugln("end-of-block")
			break // end-of-block
		} else {
			length, err := d.parseFixedHuffmanLength(val)
			if err != nil {
				return err
			}
			debugln(" -> length:", length)

			distance, err := d.parseFixedHuffmanDistance()
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

var fixedHuffmanCodes *huffmanNode

func initFixedHuffmanCodes() (err error) {
	N := 288
	fixedHuffmanCodeLengths := make([]int, N)
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

	fixedHuffmanCodes, err = generateTreeNumbered(fixedHuffmanCodeLengths)

	return err
}

var baseFixedHuffmanDistances = []uint64{
	1, 2, 3, 4, 5, 7, 9, 13, 17, 25,
	33, 49, 65, 97, 129, 193, 257, 385, 513, 769,
	1025, 1537, 2049, 3073, 4097, 6145, 8193, 12289, 16385, 24577,
}

func (d *decompressor) parseFixedHuffmanDistance() (int, error) {
	distcode := d.istream.nextBitsHighFirst(5)
	if distcode > 29 {
		return 0, corruptFileError("unexpected distcode")
	}

	nExtraBits := 0
	if distcode > 1 {
		nExtraBits = int(distcode)/2 - 1
	}

	debugln(" -> distcode:", distcode)
	debugln(" -> # extra bits:", nExtraBits)
	debugln(" -> base dist ", baseFixedHuffmanDistances[distcode])

	extraBits := d.istream.nextBitsLowFirst(nExtraBits)
	debugln(" -> extra bits ", extraBits)

	return int(baseFixedHuffmanDistances[distcode] + extraBits), nil
}

var baseFixedHuffmanLengths = []uint64{
	3, 4, 5, 6, 7, 8, 9, 10,
	11, 13, 15, 17,
	19, 23, 27, 31,
	35, 43, 51, 59,
	67, 83, 99, 115,
	131, 163, 195, 227,
	258,
}

func (d *decompressor) parseFixedHuffmanLength(lencode uint64) (int, error) {

	if lencode < 257 || lencode > 285 {
		return 0, corruptFileError("unexpected lencode")
	}

	nExtraBits := 0
	if lencode > 264 && lencode < 285 {
		nExtraBits = int(lencode-265)/4 + 1
	}

	debugln(" -> lencode", lencode)
	debugln(" -> # extra bits", nExtraBits)

	baseLength := baseFixedHuffmanLengths[lencode-257]
	debugln(" -> base length ", baseLength)

	extraBits := d.istream.nextBitsLowFirst(nExtraBits)
	debugln(" -> extra bits ", extraBits)

	return int(baseLength + extraBits), nil
}
