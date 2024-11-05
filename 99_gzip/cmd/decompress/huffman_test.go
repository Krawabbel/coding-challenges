package main

import (
	"fmt"
	"io"
	"testing"
)

func (n *huffmanNode) get(code string) (uint64, error) {

	if n == nil {
		return 0, fmt.Errorf("code does not exist")
	}

	if len(code) > 0 {
		switch code[0] {
		case '0':
			return n.left.get(code[1:])
		case '1':
			return n.right.get(code[1:])
		default:
			return 0, fmt.Errorf("invalid code")
		}

	}

	if !n.isLeaf {
		return 0, fmt.Errorf("code does not specify leaf")
	}

	return n.element, nil
}

func TestNilNode(t *testing.T) {
	n := new(huffmanNode)
	n.insertElement(0b101, 3, 12)
	val, err := n.get("101")
	if err != nil {
		t.Fatal(err)
	}
	if val != 12 {
		t.Fatal(val)
	}

}

func Test_generateTree(t *testing.T) {
	tree_len := []int{3, 3, 3, 3, 3, 2, 4, 4}
	root, err := generateTree(tree_len)
	if err != nil {
		t.Fatal(err)
	}

	codes := []string{"010", "011", "100", "101", "110", "00", "1110", "1111"}

	for i, code := range codes {
		want := uint64(i)
		if have, err := root.get(code); err != nil {
			t.Fatal(err)
		} else if have != want {
			t.Fatalf("fail for '%s': have %d, got %d", string(byte(have)+'A'), have, want)
		}
	}
}

func Test_decompressor_parseLength(t *testing.T) {
	type fields struct {
		istream *bitstream
		info    map[string][]byte
		ostream io.Writer
	}
	type args struct {
		val uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"1", fields{istream: newBitstream([]byte{0b00})}, args{269}, 19},
		{"2", fields{istream: newBitstream([]byte{0b01})}, args{269}, 20},
		{"3", fields{istream: newBitstream([]byte{0b10})}, args{269}, 21},
		{"4", fields{istream: newBitstream([]byte{0b11})}, args{269}, 22},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &decompressor{
				istream: tt.fields.istream,
				info:    tt.fields.info,
				ostream: tt.fields.ostream,
			}
			if got, _ := d.parseHuffmanLength(tt.args.val); got != tt.want {
				t.Errorf("decompressor.parseLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
