package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

const (
	maxHistoryLength = 0x8000
)

var DEBUG = true

func init() {
	initFixedHuffmanCodes()
}

func main() {
	// path := os.Args[1]
	path := "/home/dominik/projects/coding-challenges/99_gzip/small.txt.gz"

	flag.BoolVar(&DEBUG, "debug", false, "debug")
	flag.Parse()

	if err := decompress(os.Stdout, path); err != nil {
		panic(err)
	}

}

func decompress(w io.Writer, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	d := newDecompressor(w, raw)

	if err := d.parseHeader(); err != nil {
		return err
	}

	for key, val := range d.info {
		debugf("[INFO] %s: '%s' %v\n", key, string(val), val)
	}

	if err := d.parseData(); err != nil {
		return err
	}

	return nil

}

type decompressor struct {
	istream *bitstream
	info    map[string][]byte
	ostream io.Writer
	history []byte
}

func newDecompressor(w io.Writer, data []byte) *decompressor {
	return &decompressor{
		istream: newBitstream(data),
		ostream: w,
		history: []byte{},
	}
}

func (d *decompressor) parseData() error {

	for {
		if eof, err := d.parseBlock(); err != nil {
			return err
		} else if eof {
			return nil
		}
	}

}

func (d *decompressor) parseNoCompression() {
	panic("not yet implemented: no compression")
	d.istream.skipBits()
}

func (d *decompressor) push(bs ...byte) {

	debug("'")
	if n, err := d.ostream.Write(bs); err != nil {
		panic(err)
	} else if n != len(bs) {
		panic("output buffer too small")
	}
	debug("'")

	d.history = append(d.history, bs...)
	for len(d.history) >= maxHistoryLength {
		d.history = d.history[1:]
	}

}

func (d *decompressor) parseFixedHuffmanCodes() {

	for {

		val := d.parseValue(fixedHuffmanCodes)

		debug(" -> ", val, " -> ")

		if val < 256 {
			d.push(byte(val)) // literal
		} else if val == 256 {
			debugln("end-of-block")
			break // end-of-block
		} else {
			length := d.parseLength(val)

			distance := d.parseDistance()

			debugf(" -> <length: %d, distance: %d> -> ", length, distance)

			d.repeat(len(d.history)-distance, length)

		}
		debugln()
	}
}

func (d *decompressor) repeat(start int, length int) {
	for offset := range length {
		ptr := start + offset
		d.push(d.history[ptr])
	}
}

var baseDistances = []uint64{
	1, 2, 3, 4, 5, 7, 9, 13, 17, 25,
	33, 49, 65, 97, 129, 193, 257, 385, 513, 769,
	1025, 1537, 2049, 3073, 4097, 6145, 8193, 12289, 16385, 24577,
}

func (d *decompressor) parseDistance() int {
	distcode := d.istream.nextBitsRev(5)
	if distcode > 29 {
		panic("unexpected distcode")
	}

	nExtraBits := 0
	if distcode > 1 {
		nExtraBits = int(distcode)/2 - 1
	}

	debugln("\ndistcode", distcode)

	debugln("# extra bits", nExtraBits)

	extraBits := d.istream.nextBits(nExtraBits)

	debugln("\nbase dist ", baseDistances[distcode])
	debugln("extra bits ", extraBits)

	return int(baseDistances[distcode] + extraBits)
}

var baseLengths = []uint64{
	3, 4, 5, 6, 7, 8, 9, 10,
	11, 13, 15, 17,
	19, 23, 27, 31,
	35, 43, 51, 59,
	67, 83, 99, 115,
	131, 163, 195, 227,
	0,
}

func (d *decompressor) parseLength(lencode uint64) int {

	if lencode < 257 || lencode > 285 {
		panic("unexpected lencode")
	}

	nExtraBits := 0
	if lencode > 264 {
		nExtraBits = int(lencode-265)/4 + 1
	}

	debugln("\nlencode", lencode)
	debugln("# extra bits", nExtraBits)

	extraBits := d.istream.nextBits(nExtraBits)

	baseLength := baseLengths[lencode-257]

	debugln("\nbase length ", baseLength)
	debugln("extra bits ", extraBits)

	return int(baseLength + extraBits)
}

// func (d *decompressor) parseLengthOld(val uint64) int {

// 	defer debugln("|")

// 	switch val {
// 	case 257, 258, 259, 260, 261, 262, 263, 264:
// 		return int(val) - 257 + 3
// 	case 265, 266, 267, 268:
// 		extra := d.istream.nextBitsRev(1)
// 		return 2*(int(val)-265) + 11 + int(extra)
// 	case 269, 270, 271, 272:
// 		extra := d.istream.nextBitsRev(2)
// 		return 4*(int(val)-269) + 19 + int(extra)
// 	case 273, 274, 275, 276:
// 		extra := d.istream.nextBitsRev(3)
// 		return 8*(int(val)-273) + 35 + int(extra)
// 	case 277, 278, 279, 280:
// 		extra := d.istream.nextBitsRev(4)
// 		return 16*(int(val)-277) + 67 + int(extra)
// 	case 281, 282, 283, 284:
// 		extra := d.istream.nextBitsRev(5)
// 		return 32*(int(val)-281) + 115 + int(extra)
// 	case 285:
// 		return 258
// 	}
// 	panic("unexpected length code")
// }

func (d *decompressor) parseValue(node *huffmanNode) uint64 {
	if node.isLeaf {
		return node.element
	}
	if d.istream.nextBool() {
		return d.parseValue(node.right)
	}
	return d.parseValue(node.left)
}

func (d *decompressor) parseDynamicHuffmanCodes() {
	panic("not yet implemented: dynamic Huffman codes")
}

func (d *decompressor) parseBlock() (bool, error) {

	// read block header
	bfinal := d.istream.nextBool()

	btype := d.istream.nextBits(2)

	switch btype {
	case 0b00:
		d.parseNoCompression()
	case 0b01:
		d.parseFixedHuffmanCodes()
	case 0b10:
		d.parseDynamicHuffmanCodes()
	case 0b11:
		return false, fmt.Errorf("not a GZIP file: unexpected BTYPE 0b11")
	}

	return bfinal, nil
}

func (d *decompressor) parseFEXTRA() {
	size := uint64(binary.LittleEndian.Uint16(d.istream.nextBytes(2)))
	d.info["fextra"] = d.istream.nextBytes(int(size))
}

func (d *decompressor) parseCSTRING(key string) func() {
	return func() {
		cstr := []byte{}
		for {
			c := d.istream.nextByte()
			if c == 0x00 {
				break
			}
			cstr = append(cstr, c)
		}
		d.info[key] = cstr
	}
}

func (d *decompressor) parseFHCRC() {
	d.info["fhcrc"] = d.istream.nextBytes(2)
}

func (d *decompressor) parseHeader() error {

	parseQueue := []func(){}

	d.info = make(map[string][]byte)

	if magic := d.istream.nextBytes(2); !slices.Equal(magic, []byte{0x1F, 0x8B}) {
		return fmt.Errorf("not a gzip file: magic numbers [0x%2X 0x%2X] are not [0x1F 0x8B]", magic[0], magic[1])
	}

	if method := d.istream.nextByte(); method != 0x08 {
		return fmt.Errorf("unexpected compression method %x", method)
	}

	flags := d.istream.nextByte()

	if (flags & 0x01) != 0 {
		fmt.Println("[WARNING]", "FTEXT: If set the uncompressed data needs to be treated as text instead of binary data.")
	}

	if (flags & 0x04) != 0 {
		parseQueue = append(parseQueue, d.parseFEXTRA)
	}

	if (flags & 0x08) != 0 {
		parseQueue = append(parseQueue, d.parseCSTRING("fname"))
	}

	if (flags & 0x10) != 0 {
		parseQueue = append(parseQueue, d.parseCSTRING("fcomment"))
	}

	if (flags & 0x02) != 0 {
		parseQueue = append(parseQueue, d.parseFHCRC)
	}

	if (flags & (0x20 | 0x40 | 0x80)) != 0 {
		return fmt.Errorf("unexpected flags 0x%02X", flags)
	}

	mtime := binary.LittleEndian.Uint32(d.istream.nextBytes(4))
	d.info["mtime"] = []byte(fmt.Sprint(time.Unix(int64(mtime), 0)))

	if extra := d.istream.nextByte(); extra != 0x00 {
		return fmt.Errorf("unexpected extra flags %x", extra)
	}

	d.info["os"] = []byte(fmt.Sprintf("0x%02X", d.istream.nextByte()))

	for _, parseFun := range parseQueue {
		parseFun()
	}

	return nil
}
