package main

func calcCheckSum[T int64 | uint64](header []byte) T {
	have := T(0)
	for i, v := range header {
		if i >= 148 && i < (148+8) {
			have += 32
		} else {
			have += T(v)
		}
	}

	return have
}

func ceil512(size uint64) uint64 {
	blocks := size / 512
	if (size % 512) > 0 {
		return blocks + 1
	}
	return blocks
}
