package main

type bitstream struct {
	data []byte
	// dataPtr int
	// bitPtr  int
	ptr int
}

func newBitstream(data []byte) *bitstream {
	s := &bitstream{data: data}
	return s
}

func (s *bitstream) skipToNextByte() {
	for (s.ptr % 8) > 0 {
		debug('x')
		s.ptr++
	}
}

func (s *bitstream) nextBool() bool {

	if s.eof() {
		panic("bitstream EOF")
	}

	bytePos := s.ptr / 8
	bitPos := s.ptr % 8

	s.ptr++

	b := (s.data[bytePos] & (1 << bitPos)) > 0

	if b {
		debug(1)
	} else {
		debug(0)
	}

	if s.ptr%8 == 0 {
		debug(" ")
	}

	return b
}

func (s *bitstream) nextByte() byte {
	return byte(s.nextBits(8))
}

func (s *bitstream) nextBytes(n int) []byte {
	bs := make([]byte, n)
	for i := range n {
		bs[i] = s.nextByte()
	}
	return bs
}

func (s *bitstream) nextBit() uint64 {
	if s.nextBool() {
		return 1
	}
	return 0
}

func (s *bitstream) nextBits(n int) uint64 {
	bits := uint64(0)

	for i := 0; i < n; i++ {
		bits |= (s.nextBit() << i)
	}

	return bits
}

func (s *bitstream) neof() bool {
	return (s.ptr < len(s.data)*8)
}

func (s *bitstream) eof() bool {
	return !s.neof()
}
