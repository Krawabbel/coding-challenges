package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

func unpack(tar []byte) ([]file, error) {

	files := make([]file, 0)

	arch := archive{raw: tar, ptr: 0}

	for {
		eof, err := arch.eof()
		if err != nil {
			return nil, err
		}

		if eof {
			return files, nil
		}

		f, err := arch.parseFile()
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}
}

type archive struct {
	raw   []byte
	ptr   uint64
	files []file
}

type file struct {
	info *fileInfo
	data []byte
}

type fileInfo struct {
	fileName          string
	fileMode          uint64
	ownerUserID       uint64
	ownerGroupID      uint64
	fileSize          uint64
	lastModDate       time.Time
	checkSum          uint64
	linkName          string
	typeFlag          byte
	ustarFlag         bool
	ustarVersion      string
	ownerUserName     string
	ownerGroupName    string
	deviceMajorNumber string
	deviceMinorNumber string
	fileNamePrefix    string
}

func (i *fileInfo) filename() string {
	return i.fileNamePrefix + i.fileName
}

func (a *archive) checkOffset(offset uint64) error {
	if a.ptr+offset < a.size() {
		return nil
	}
	return fmt.Errorf("unexpected end of archive: pos = %d, len = %d", a.ptr+offset, a.size())
}

func (a *archive) size() uint64 {
	return uint64(len(a.raw))
}

func (a *archive) eof() (bool, error) {

	next, err := a.peek(1024)
	if err != nil {
		return true, err
	}

	for _, b := range next {
		if b != 0x00 {
			return false, nil
		}
	}

	return true, nil
}

func (a *archive) pop(n uint64) ([]byte, error) {
	if err := a.checkOffset(n); err != nil {
		return nil, err
	}

	a.ptr += n
	return a.raw[a.ptr-n : a.ptr], nil
}

func (a *archive) peek(n uint64) ([]byte, error) {
	if err := a.checkOffset(n); err != nil {
		return nil, err
	}
	return a.raw[a.ptr : a.ptr+n], nil
}

func (a *archive) parseFile() (file, error) {

	info, err := a.parseHeader()
	if err != nil {
		return file{}, err
	}

	data, err := a.parseData(info.fileSize)
	if err != nil {
		return file{}, err
	}

	return file{info: info, data: data}, nil
}

func (a *archive) parseData(size uint64) ([]byte, error) {

	blocks := ceil512(size)

	padded, err := a.pop(512 * blocks)
	if err != nil {
		return nil, err
	}

	return padded[:size], nil
}

func (a *archive) parseHeader() (*fileInfo, error) {

	header, err := a.pop(512)
	if err != nil {
		return nil, err
	}

	info := new(fileInfo)
	info.fileName = parseText(header, 0, 100)
	info.fileMode = parseOctal(header, 100, 7)
	info.ownerUserID = parseOctal(header, 108, 7)
	info.ownerGroupID = parseOctal(header, 116, 7)
	info.fileSize = parseOctal(header, 124, 11)
	info.lastModDate = parseTime(header, 136, 11)
	info.checkSum = parseOctal(header, 148, 6)
	info.linkName = parseText(header, 157, 100)
	info.typeFlag = header[156]
	info.ustarFlag = parseText(header, 257, 5) == "ustar"
	info.ustarVersion = parseText(header, 263, 2)
	info.ownerUserName = parseText(header, 265, 32)
	info.ownerGroupName = parseText(header, 297, 32)
	info.deviceMajorNumber = parseText(header, 329, 8)
	info.deviceMinorNumber = parseText(header, 337, 8)
	info.fileNamePrefix = parseText(header, 345, 155)

	if err := verifyChecksum(info.checkSum, header); err != nil {
		fmt.Fprintf(os.Stderr, "[WARNING] tar-ball possibly corrupted: %s\n", err.Error())
	}

	return info, nil
}

func verifyChecksum(want uint64, header []byte) error {
	if want == calcCheckSum[uint64](header) {
		return nil
	}

	if int64(want) == calcCheckSum[int64](header) {
		return nil
	}

	return fmt.Errorf("checksum verification failed")
}

func parseOctal(header []byte, offset uint64, size uint64) uint64 {
	octal := string(header[offset : offset+size])
	s, _ := strconv.ParseUint(octal, 8, 64)
	return s
}

func parseTime(header []byte, offset uint64, size uint64) time.Time {
	octal := string(header[offset : offset+size])
	t, _ := strconv.ParseInt(octal, 8, 64)
	return time.Unix(t, 0)
}

func parseText(header []byte, offset uint64, size uint64) string {
	raw := header[offset : offset+size]
	stripped := bytes.ReplaceAll(raw, []byte{0x00}, []byte{})
	return string(stripped)
}
