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

var DEBUG = false

func init() {
	if err := initFixedHuffmanLitValCodes(); err != nil {
		panic(err)
	}

	if err := initFixedHuffmanDistCodes(); err != nil {
		panic(err)
	}
}

func main() {
	// path := "/home/dominik/projects/coding-challenges/99_gzip/small.txt.gz"

	flag.BoolVar(&DEBUG, "debug", false, "debug")
	flag.Parse()

	path := flag.Arg(0)

	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := decompress(os.Stdout, raw); err != nil {
		panic(err)
	}

}

func decompress(w io.Writer, raw []byte) error {

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

func newDecompressor(writer io.Writer, data []byte) *decompressor {
	return &decompressor{
		istream: newBitstream(data),
		ostream: writer,
		history: []byte{},
	}
}

func (d *decompressor) parseData() error {
	for i := 0; ; i++ {
		debugln("*** BLOCK", i, "***")
		if eof, err := d.parseBlock(); err != nil {
			return err
		} else if eof {
			debugln("*** END OF BLOCK(S) ***")
			return nil
		}
	}
}

func (d *decompressor) parseNoCompression() error {

	debugln(" -> no compression")

	d.istream.skipToNextByte()
	length := binary.LittleEndian.Uint16(d.istream.nextBytes(2))
	nLength := binary.LittleEndian.Uint16(d.istream.nextBytes(2))

	if length != ^nLength {
		return corruptFileError("no-compression check failed")
	}

	debugf("\nLEN: %04X, NLEN: %04X, sum: %04X\n", length, nLength, length-^nLength)

	d.push(d.istream.nextBytes(int(length))...)

	return nil
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

func (d *decompressor) repeat(start int, length int) {
	for offset := range length {
		ptr := start + offset
		d.push(d.history[ptr])
	}
}

func (d *decompressor) parseBlock() (bool, error) {

	// read block header
	bfinal := d.istream.nextBool()

	btype := d.istream.nextBits(2)

	switch btype {
	case 0b00:
		if err := d.parseNoCompression(); err != nil {
			return false, err
		}
	case 0b01:
		if err := d.parseFixedHuffmanCodes(); err != nil {
			return false, err
		}
	case 0b10:
		if err := d.parseDynamicHuffmanCodes(); err != nil {
			return false, err
		}
	case 0b11:
		return false, corruptFileError("unexpected BTYPE 0b11")
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
		return corruptFileError("magic numbers [0x%2X 0x%2X] are not [0x1F 0x8B]", magic[0], magic[1])
	}

	if method := d.istream.nextByte(); method != 0x08 {
		return corruptFileError("unexpected compression method %x", method)
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
		return corruptFileError("unexpected flags 0x%02X", flags)
	}

	mtime := binary.LittleEndian.Uint32(d.istream.nextBytes(4))
	d.info["mtime"] = []byte(fmt.Sprint(time.Unix(int64(mtime), 0)))

	if extra := d.istream.nextByte(); extra != 0x00 {
		return corruptFileError("unexpected extra flags %x", extra)
	}

	d.info["os"] = []byte(fmt.Sprintf("0x%02X", d.istream.nextByte()))

	for _, parseFun := range parseQueue {
		parseFun()
	}

	return nil
}
