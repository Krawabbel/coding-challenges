package main

import (
	"testing"
)

func Test_bitstream_nextBits(t *testing.T) {
	type fields struct {
		data []byte
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint64
	}{
		{"1", fields{data: []byte{0xAB}}, args{8}, 0xAB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &bitstream{
				data: tt.fields.data,
			}
			if got := s.nextBitsLowFirst(tt.args.n); got != tt.want {
				t.Errorf("bitstream.nextBits() = %08b, want %08b", got, tt.want)
			}
		})
	}
}

func Test_bitstream_nextByte(t *testing.T) {
	type fields struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{"1", fields{data: []byte{0xAB}}, 0xAB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &bitstream{
				data: tt.fields.data,
			}
			if got := s.nextByte(); got != tt.want {
				t.Errorf("bitstream.nextByte() = %v, want %v", got, tt.want)
			}
		})
	}
}
