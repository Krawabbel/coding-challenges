package main

var clenIdxs = []int{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}

func (d *decompressor) parseDynamicHuffmanCodes() error {

	debugln(" -> dynamic huffman compression")

	hlit := d.istream.nextBits(5)
	nlit := int(hlit) + 257
	debugln(" -> hlit:", hlit, "-> nlit:", nlit)

	hdist := d.istream.nextBits(5)
	ndist := int(hdist) + 1
	debugln(" -> hdist:", hdist, "-> ndist:", ndist)

	hclen := d.istream.nextBits(4)
	nclen := int(hclen) + 4
	debugln(" -> hclen:", hclen, "-> nclen:", nclen)

	clens := make([]int, 19)
	for i := range nclen {
		clen := d.istream.nextBits(3)
		idx := clenIdxs[i]
		clens[idx] = int(clen)
	}

	debugln(" => code length tree:", clens)

	clTree, err := generateTree(clens)
	if err != nil {
		return err
	}

	debugln(" => literal/value tree")

	litValTree, err := d.parseAndDecodeDynamicHuffmanTree(clTree, nlit)
	if err != nil {
		return err
	}

	debugln(" => distance tree")

	distTree, err := d.parseAndDecodeDynamicHuffmanTree(clTree, ndist)
	if err != nil {
		return err
	}

	return d.parseHuffmanCodes(litValTree, distTree)
}

func (d *decompressor) parseAndDecodeDynamicHuffmanTree(cl *huffmanNode, n int) (*huffmanNode, error) {

	codeLengths := make([]int, n)
	i := 0
	for i < n {
		codeLengthCode := cl.getElement(d.istream)
		switch {
		case codeLengthCode < 16:
			codeLengths[i] = int(codeLengthCode)
			i++
		case codeLengthCode == 16:
			replen := d.istream.nextBits(2) + 3
			last := codeLengths[len(codeLengths)-1]
			for range replen {
				codeLengths[i] = last
				i++
			}
		case codeLengthCode == 17:
			replen := d.istream.nextBits(3) + 3
			for range replen {
				codeLengths[i] = 0
				i++
			}
		case codeLengthCode == 18:
			replen := d.istream.nextBits(7) + 11
			for range replen {
				codeLengths[i] = 0
				i++
			}
		}
	}

	if i > n {
		panic("i > n")
	}

	debugf("\n -> code lengths (# = %d): %v\n", len(codeLengths), codeLengths)

	return generateTree(codeLengths)
}
