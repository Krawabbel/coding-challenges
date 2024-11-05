package main

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestSmall(t *testing.T) {
	testhelper_decompress(t, "../../small.txt")
}

func TestRandom(t *testing.T) {
	testhelper_decompress(t, "../../random.txt")
}

func TestLong(t *testing.T) {
	// testhelper_decompress(t, "../../long.txt")
}

func testhelper_decompress(t *testing.T, path string) {
	cmd := exec.Command("gzip", "-c", path)
	given, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	w := new(bytes.Buffer)

	DEBUG = false
	if err := decompress(w, given); err != nil {
		t.Fatal(err)
	}

	have := w.Bytes()

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(have, want) {
		for i := range min(len(have), len(want)) {
			t.Log("have", have[i], "want", want[i])
			if have[i] != want[i] {
				t.Log(" <---")
			}
			t.Log("\n")
		}
		t.Fail()
	}
}
