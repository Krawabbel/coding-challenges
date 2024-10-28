package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		panic("please provide a nickname")
	}

	name := os.Args[1]

	conn, err := net.Dial("tcp", ":7007")
	if err != nil {
		panic(err)
	}

	go func() {
		scan := bufio.NewScanner(os.Stdin)
		for scan.Scan() {
			data := scan.Bytes()
			_, err := fmt.Fprintf(conn, "%s: %s\n", name, data)
			if err != nil {
				panic(err)
			}
		}
		if err := scan.Err(); err != nil {
			panic(err)
		}
	}()

	for {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(os.Stdout, "%s", buf[:n])
		}
	}
}
