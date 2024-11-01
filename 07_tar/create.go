package main

import (
	"fmt"
	"os"
)

func createHeader(path string) ([]byte, error) {
	header := make([]byte, 512)

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	copyText(header, 0, 100, stat.Name())

	copyOctal(header, 100, 8, uint64(stat.Mode()))

	// if linuxstat, ok := stat.Sys().(*syscall.Stat_t); ok {
	// 	copyOctal(header, 108, 8, uint64(linuxstat.Uid))
	// 	copyOctal(header, 116, 8, uint64(linuxstat.Gid))
	// }

	copyOctal(header, 124, 12, uint64(stat.Size()))

	copyOctal(header, 136, 12, uint64(stat.ModTime().Unix()))

	if stat.IsDir() {
		return nil, fmt.Errorf("don't know how to pack directories")
	}
	// header[156] = 0x00 // normal file

	// copyText(header, 257, 6, "ustar")
	// copyText(header, 265, 32, "dominik")
	// copyText(header, 297, 32, "dominik")

	checksum := calcCheckSum[uint64](header)
	copyOctal(header, 148, 6, checksum)
	header[155] = ' '

	return header, nil
}

func copyText(header []byte, offset, size uint64, val string) {
	copy(header[offset:offset+size], []byte(val))
}

func copyOctal(header []byte, offset, size uint64, val uint64) {
	format := fmt.Sprintf("%%0%do", size-1)
	s := fmt.Sprintf(format, val)
	copyText(header, offset, size-1, s)
}

func pack(path string) ([]byte, error) {

	tar, err := createHeader(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tar = append(tar, data...)

	padding := make([]byte, 512-(uint64(len(data))%512))
	tar = append(tar, padding...)

	return tar, nil

}
