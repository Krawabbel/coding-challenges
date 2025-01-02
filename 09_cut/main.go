package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	fArg := flag.String("fields", "", "fields")
	dArg := flag.String("delimiter", "\t", "delimiter")
	flag.Parse()

	input, err := readInput(flag.Args())
	if err != nil {
		panic(err)
	}

	fields, err := parseFields(*fArg)
	if err != nil {
		panic(err)
	}

	delim := *dArg

	output := os.Stdout

	err = cut(input, output, delim, fields)
	if err != nil {
		panic(err)
	}

}

func parseFields(s string) ([]int, error) {
	l := strings.ReplaceAll(s, " ", ",")
	fields := make([]int, 0)
	for _, r := range strings.Split(l, ",") {
		fs, err := parseRangeOrValue(r)
		if err != nil {
			return nil, err
		}
		fields = append(fields, fs...)
	}
	return fields, nil
}

func parseRangeOrValue(s string) ([]int, error) {
	if strings.Contains(s, "-") {
		return parseRange(s)
	}
	return parseValue(s)
}

func parseRange(r string) ([]int, error) {

	parts := strings.Split(r, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range '%s': must have length 2", r)
	}
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	b, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	if a > b {
		return nil, fmt.Errorf("invalid range '%s': first argument must be smaller than second", r)
	}

	ids := make([]int, 0, b-a+1)
	for i := a; i <= b; i++ {
		ids = append(ids, i)
	}

	return ids, nil
}

func parseValue(v string) ([]int, error) {
	val, err := strconv.Atoi(v)
	if err != nil {
		return nil, err
	}
	return []int{val}, nil
}

func cut(r io.Reader, w io.Writer, delim string, fs []int) error {
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		line := scan.Text()
		parts := strings.Split(line, delim)
		cols := make([]string, 0)
		for _, f := range fs {
			id := f - 1
			if id >= 0 && id < len(parts) {
				cols = append(cols, parts[id])
			}
		}
		out := strings.Join(cols, delim)
		_, err := fmt.Fprintln(w, out)
		if err != nil {
			return err
		}
	}
	return scan.Err()
}

func readInput(src []string) (io.Reader, error) {
	if len(src) == 0 || len(src) == 1 && src[0] == "-" {
		return os.Stdin, nil
	}
	return readFiles(src)
}

func readFiles(paths []string) (io.Reader, error) {
	rs := make([]io.Reader, len(paths))
	for i, path := range paths {
		r, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		rs[i] = r
	}
	return io.MultiReader(rs...), nil
}
