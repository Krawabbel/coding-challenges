package main

import (
	"fmt"
	"net"
	"strings"
)

type word = uint16

func main() {

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	req := encodeRequest("dns.google.com")

	if _, err := conn.Write(req); err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	resp := buf[:n]

	fmt.Printf("REQUEST:  % x\n", req)
	fmt.Printf("RESPONSE: % x\n", resp)

	fmt.Println("DETAILS:")
	for _, b := range resp {
		fmt.Printf(" 0x%x / 0b%08b/ %d (%s)\n", b, b, b, string(b))
	}
	fmt.Println()

	if err := decodeResponse(resp); err != nil {
		panic(err)
	}
}

func encodeRequest(name string) []byte {
	h := encodeHeader()
	q := encodeQuestion(name)
	return append(h[:], q...)

}

func encodeHeader() [12]byte {

	header := [12]byte{}

	// ID
	id := word(22)
	header[0] = byte(id >> 8)
	header[1] = byte(id)

	// QR(1), OPCODE(4), AA(1), TC(1), RD(1), RA(1), Z(3), RCODE(4)
	flags := word(1) << 8
	header[2] = byte(flags >> 8)
	header[3] = byte(flags)

	// QDCOUNT
	nQuestions := word(0x0001)
	header[4] = byte(nQuestions >> 8)
	header[5] = byte(nQuestions)

	return header
}

func encodeQuestion(name string) []byte {
	q := make([]byte, 0)

	// QNAME (n)
	for _, sub := range strings.Split(name, ".") {
		q = append(q, byte(len(sub)))
		q = append(q, []byte(sub)...)
	}
	q = append(q, 0x00)

	// QTYPE (2)
	q = append(q, 0x00, 0x01) // -> A: host address

	// QCLASS (2)
	q = append(q, 0x00, 0x01) // -> IN: the internet

	return q
}

func decodeResponse(msg []byte) error {
	id := join(msg[0], msg[1])
	fmt.Printf("Question/Response ID: %v\n", id)

	isResponse := msg[2]&(1<<0) > 0
	if !isResponse {
		return fmt.Errorf("not a response")
	}

	qdcount := join(msg[4], msg[5])
	if qdcount != 1 {
		return fmt.Errorf("to many questions")
	}
	ptr := 12

	// qdbegin := ptr
	for range qdcount {
		begin := ptr
		for {
			l := int(msg[ptr])
			ptr++
			if l == 0 {
				break
			}
			ptr += l
		}
		fmt.Printf("QUESTION: % x (%s)\n", msg[begin:ptr], string(msg[begin:ptr]))

		ptr += 2 // qtype
		ptr += 2 // qclass
	}
	// qdend := ptr

	ancount := join(msg[6], msg[7])
	// anbegin := ptr
	for range ancount {
		// begin := ptr
		if (msg[ptr] & 0b11000000) != 0b11000000 {
			return fmt.Errorf("no domain pointer")
		}
		offset := join(msg[ptr]&0b00111111, msg[ptr+1])
		domain := msg[offset : offset+10]

		fmt.Printf("DOMAIN: % x (%s)", domain, string(domain))
		ptr += 2

		if msg[ptr] != 0x00 || msg[ptr+1] != 0x01 {
			return fmt.Errorf("response resource not of type 'A' (host name)")
		}
		ptr += 2 // type

		if msg[ptr] != 0x00 || msg[ptr+1] != 0x01 {
			return fmt.Errorf("response resource not of class 'IN' (internet)")
		}
		ptr += 2 // class

		ptr += 4 // ttl

		rdlength := join(msg[ptr], msg[ptr+1])
		ptr += 2

		rdata := msg[ptr : ptr+int(rdlength)]
		ptr += int(rdlength)
		fmt.Printf(" -> % x\n", rdata)
	}
	// anend := ptr

	// nscount := join(msg[8], msg[9])
	// arcount := join(msg[10], msg[11])

	return nil
}

func join(hi, lo byte) word {
	return (word(hi) << 8) | word(lo)
}
