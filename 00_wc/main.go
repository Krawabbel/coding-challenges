package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const descr = `Print newline, word, and byte counts for each FILE, and a total line if more than one FILE is specified. A word is a non-zero-length sequence of printable characters delimited by white space.
With no FILE, or when FILE is -, read standard input.
The options below may be used to select which counts are printed, always in the following order: newline, word, character, byte, maximum line length.`

func main() {

	fLines := flag.Bool("lines", false, "print the newline counts")
	fWords := flag.Bool("words", false, "print the word counts")
	fChars := flag.Bool("chars", false, "print the character counts")
	fBytes := flag.Bool("bytes", false, "print the byte counts")
	fWidth := flag.Bool("max-line-length", false, "print the maximum display width")

	flag.BoolFunc("version", "output version information and exit", func(s string) error {
		fmt.Println("Word Count Version 0.1")
		os.Exit(0)
		return nil
	})

	flag.BoolFunc("help", "display this help and exit", func(_ string) error {
		fmt.Println(descr)
		flag.Usage()
		os.Exit(0)
		return fmt.Errorf("blablabla")
	})

	flag.Parse()

	paths := flag.Args()
	tLines := uint64(0)
	for _, path := range paths {

		nLines, nWords, nChars, nBytes, maxWidth, err := wcPath(path)
		if err != nil {
			panic(err)
		}
		tLines += nLines

		useDefaultFlags := !*fLines && !*fWords && !*fChars && !*fBytes && !*fWidth
		printStatistic(*fLines || useDefaultFlags, path, "lines:", nLines)
		printStatistic(*fWords || useDefaultFlags, path, "words:", nWords)
		printStatistic(*fChars, path, "chars:", nChars)
		printStatistic(*fBytes || useDefaultFlags, path, "bytes:", nBytes)
		printStatistic(*fWidth, path, "max-width", maxWidth)
	}

	printStatistic(len(paths) > 1, "total", "lines", tLines)
}

func printStatistic(flag bool, path, name string, value uint64) {
	if flag {
		fmt.Fprintf(os.Stdout, "[%s] %s %d\n", path, name, value)
	}
}

func wcPath(path string) (nLines, nWords, nChars, nBytes, maxWidth uint64, err error) {

	file, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	defer file.Close()

	nLines = 0
	nWords = 0
	nChars = 0
	nBytes = 0
	maxWidth = 0

	scanner := bufio.NewScanner(file)
	scanner.Split(splitFunc)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		nLines += 1
		nWords += uint64(len(strings.Fields(line)))
		nBytes += uint64(len(line))
		wChars := uint64(utf8.RuneCountInString(line))
		nChars += uint64(wChars)
		if wChars > maxWidth {
			maxWidth = wChars
		}
	}

	return nLines, nWords, nChars, nBytes, maxWidth, nil
}

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
