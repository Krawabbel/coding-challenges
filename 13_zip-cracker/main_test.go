package main

import "testing"

func Test_calcCRC32(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		key  uint32
		char uint32
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := update_crc(tt.key, tt.char)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("calcCRC32() = %v, want %v", got, tt.want)
			}
		})
	}
}
