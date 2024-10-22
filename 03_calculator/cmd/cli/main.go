package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"calculator"
)

func main() {

	expr := func() string {

		if len(os.Args) > 1 {
			return strings.Join(os.Args[1:], "")
		}

		blob, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		return string(blob)
	}()

	result, err := calculator.Calculate(expr)
	if err != nil {
		panic(err)
	}

	fmt.Printf(">> %f\n", result)
}
