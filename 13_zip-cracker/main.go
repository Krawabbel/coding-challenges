package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/yeka/zip"
)

var SIGNATURE = []byte{0x50, 0x4b, 0x03, 0x04}

func join(bs ...byte) uint64 {
	val := uint64(0)
	for i, b := range bs {
		val |= (uint64(b) << (i * 8))
	}
	return val
}

func readNumericField(name string, data []byte, pos, length uint64) uint64 {
	bs := data[pos : pos+length]
	val := join(bs...)
	fmt.Printf("%s: %v %v\n", name, val, bs)
	return val
}

func readStringField(name string, data []byte, pos, length uint64) string {
	bs := data[pos : pos+length]
	str := string(bs)
	fmt.Printf("%s: %s %v\n", name, str, bs)
	return str
}

func check_password(cipher []byte, pass []byte, checksum []byte) bool {

	// initialize the encryption keys
	keys := []uint32{305419896, 591751049, 878082192}
	for i := range pass {
		update_keys(keys, pass[i])
	}
	fmt.Println(keys)

	// decrypt the encryption header
	plain := make([]byte, len(cipher))

	for i := range cipher {
		t := keys[2] | 2
		decrypt_byte := (t * (t ^ 1)) >> 8
		c := cipher[i] ^ byte(decrypt_byte)
		update_keys(keys, c)
		plain[i] = c
	}

	fmt.Println("(DEFLATED) CONTENT:", string(plain[12:]))

	print_bytes("have", plain)
	print_bytes("want", checksum)

	return reflect.DeepEqual(checksum, plain)
}

func update_keys(keys []uint32, char byte) {
	keys[0] = update_crc(keys[0], char)
	keys[1] += (keys[0] & 0xFF)
	keys[1] = keys[1]*134775813 + 1
	keys[2] = update_crc(keys[2], byte(keys[1]>>24))
}

func update_crc(crc uint32, char byte) uint32 {
	crc ^= uint32(char)
	for range 8 {
		if crc&1 > 0 {
			crc = (crc >> 1) ^ 0xEDB88320
		} else {
			crc >>= 1
		}
	}
	return crc
}

func Run(path string) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	signature := data[0:4]
	if !reflect.DeepEqual(signature, SIGNATURE) {
		return fmt.Errorf("not a zip file: unexpected signature want %v, have %v", SIGNATURE, signature)
	}

	readNumericField("version needed to extract (minimum)", data, 4, 2)
	readNumericField("general purpose bit flag", data, 6, 2)
	readNumericField("compression method", data, 8, 2)
	readNumericField("file last modification time", data, 10, 2)
	readNumericField("file last modification date", data, 12, 2)
	readNumericField("CRC-32 of uncompressed data", data, 14, 4)
	s := readNumericField("compressed size", data, 18, 4)
	readNumericField("uncompressed size", data, 22, 4)
	n := readNumericField("file name length", data, 26, 2)
	m := readNumericField("extra field length", data, 28, 2)
	readStringField("file name", data, 30, n)
	readNumericField("extra field", data, 30+n, m)

	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword("123")
		}
		rd, err := f.Open()
		if err != nil {
			return err
		}
		defer rd.Close()

		content, err := io.ReadAll(rd)
		if err != nil {
			return err
		}

		fmt.Println("content:", string(content))
	}

	encrypt_header := data[30+n+m : 30+n+m+12]
	fmt.Println(30 + n + m)
	print_bytes("encryption header", encrypt_header)
	checksum := data[14:18]
	if !check_password(data[30+n+m:30+n+m+s], []byte("123"), checksum) {
		return fmt.Errorf("wrong password")
	}

	return nil
}

func main() {
	must(Run("/home/dominik/projects/coding-challenges/13_zip-cracker/demo.zip"))
}

func print_bytes(name string, bs []byte) {

	fmt.Printf("%s: ", name)
	for _, b := range bs {
		fmt.Printf("0x%02X ", b)
	}
	fmt.Println()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
