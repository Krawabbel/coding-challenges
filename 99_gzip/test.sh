#!/bin/bash

gzip -c $1.txt > $1.txt.gz

# cp $1.txt tmp.txt
# echo -e "\n" >> tmp.txt
# go run ./cmd/decompress $1.txt.gz >> tmp.txt
# echo -e "\n" >> tmp.txt
# go run ./cmd/decompress --debug $1.txt.gz >> tmp.txt    

go run ./cmd/decompress $1.txt.gz > out.txt
diff $1.txt out.txt