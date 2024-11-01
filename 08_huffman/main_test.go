package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_EncodeDecode(t *testing.T) {
	tests := []struct {
		text string
	}{
		{"hello hello world"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			given := []byte(tt.text)
			code, root, rem := encode(given)
			have := decode(code, root, rem)
			want := given
			if !reflect.DeepEqual(have, want) {
				t.Errorf("compress() = %v, want %v", have, want)
			}
		})
	}
}
