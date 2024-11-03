package main

import "fmt"

func debug(a ...any) {
	if DEBUG {
		fmt.Print(a...)
	}
}

func debugf(format string, a ...any) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}

func debugln(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}

func hex(v uint64) string {
	return fmt.Sprintf("0x%X", v)
}

func bin(v uint64) string {
	return fmt.Sprintf("0b%b", v)
}

func bits(b byte) []bool {
	bs := make([]bool, 8)
	for i := range bs {
		// bs[i] = (b & (1 << (8 - i))) > 0
		bs[i] = (b & (1 << i)) > 0
	}
	return bs
}
