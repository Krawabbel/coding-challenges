package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

func main() {
	path := os.Args[1]

	if err := decompress(path); err != nil {
		panic(err)
	}
}

type decompressor struct {
	data []byte
	ptr  uint64
	info map[string][]byte
}

func decompress(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	d := decompressor{data: raw}

	if err := d.parseHeader(); err != nil {
		return err
	}

	for key, val := range d.info {
		fmt.Println("[DEBUG]", key, string(val), val)
	}

	return nil

}

func (d *decompressor) parseHeader() error {

	d.ptr = uint64(10)

	d.info = make(map[string][]byte)

	if d.data[0] != 0x1F || d.data[1] != 0x8B {
		return fmt.Errorf("not a gzip file")
	}

	if d.data[2] != 0x08 {
		return fmt.Errorf("unexpected compression method %x", d.data[2])
	}

	if (d.data[3] & 0x01) != 0 {
		fmt.Println("[WARNING]", "FTEXT: If set the uncompressed data needs to be treated as text instead of binary data.")
	}

	if (d.data[3] & 0x04) != 0 {
		size := uint64(binary.LittleEndian.Uint16(d.data[d.ptr : d.ptr+2]))
		d.ptr += 2
		d.info["fextra"] = d.data[d.ptr : d.ptr+size]
		d.ptr += size
	}

	if (d.data[3] & 0x08) != 0 {
		end := d.ptr
		for d.data[end] != 0x00 {
			end++
		}
		d.info["fname"] = d.data[d.ptr:end]
		d.ptr = end + 1
	}

	if (d.data[3] & 0x10) != 0 {
		end := d.ptr
		for d.data[end] != 0x00 {
			end++
		}
		d.info["fcomment"] = d.data[d.ptr:end]
		d.ptr = end + 1
	}

	if (d.data[3] & 0x02) != 0 {
		d.info["fhcrc"] = d.data[d.ptr : d.ptr+2]
		d.ptr += 2
	}

	if (d.data[3] & (0x20 | 0x40 | 0x80)) != 0 {
		return fmt.Errorf("unexpected flags 0x%02X", d.data[3])
	}

	unixtime := binary.LittleEndian.Uint32(d.data[4:8])
	d.info["mtime"] = []byte(fmt.Sprint(time.Unix(int64(unixtime), 0)))

	if d.data[8] != 0x00 {
		return fmt.Errorf("unexpected extra flags %x", d.data[8])
	}

	d.info["os"] = []byte(fmt.Sprintf("0x%02X", d.data[9]))

	return nil
}
