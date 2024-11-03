#!/bin/bash

gzip -c small.txt > small.txt.gz
cp small.txt tmp.txt
go run ./cmd/decompress >> tmp.txt
go run ./cmd/decompress --debug >> tmp.txt    