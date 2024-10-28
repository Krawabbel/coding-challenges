package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func logErr(err error) {
	if err != nil {
		log.Println("[ERROR]", err.Error())
	}
}

func logDbg(format string, v ...any) {
	log.Printf("%s "+format, append([]any{"[DEBUG]"}, v...)...)
}

func genId() string {
	blob := make([]byte, SECRET_LENGTH)
	n, err := rand.Reader.Read(blob)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(blob[:n])
}
