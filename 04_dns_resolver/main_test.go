package main

import (
	"fmt"
	"testing"
)

func Test_main(t *testing.T) {
	h := encodeHeader()
	q := encodeQuestion("dns.google.com")

	msg := append(h[:], q...)

	have := fmt.Sprintf("%x", msg)
	want := "00160100000100000000000003646e7306676f6f676c6503636f6d0000010001"

	if have != want {
		t.Fatalf("\nhave '%s'\nwant '%s'", have, want)
	}

}
