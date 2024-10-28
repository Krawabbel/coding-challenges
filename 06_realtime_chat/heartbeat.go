package main

import (
	"fmt"
	"io"
	"time"
)

type tickerReader struct {
	ch <-chan time.Time
}

func newChanReader(ch <-chan time.Time) *tickerReader {
	return &tickerReader{ch: ch}
}

func (r *tickerReader) Read(p []byte) (int, error) {
	t, ok := <-r.ch
	if !ok {
		return 0, io.EOF
	}

	data := []byte(t.Format("2006-01-02 15:04:05") + " keep-alive\n")
	if len(data) > len(p) {
		return 0, fmt.Errorf("target buffer too small")
	}

	return copy(p, data), nil
}

func runHeartbeat(brok broker) error {
	tic := time.NewTicker(time.Second)
	return handlePub(brok, newChanReader(tic.C))
}
