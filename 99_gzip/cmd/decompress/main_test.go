package main

import (
	"io"
	"testing"
)

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
		{"2", fields{istream: newBitstream([]byte{0b10})}, args{269}, 20},
		{"3", fields{istream: newBitstream([]byte{0b01})}, args{269}, 21},
		{"4", fields{istream: newBitstream([]byte{0b11})}, args{269}, 22},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &decompressor{
				istream: tt.fields.istream,
				info:    tt.fields.info,
				ostream: tt.fields.ostream,
			}
			if got := d.parseLength(tt.args.val); got != tt.want {
				t.Errorf("decompressor.parseLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
