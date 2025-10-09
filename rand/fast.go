package rand

import "sync/atomic"

var rnd int64

func FastInt() int {
	nVal := atomic.AddInt64(&rnd, 0x61c8864680b583eb)
	return int(nVal >> 32)
}

func FastIntn(n int) int {
	if n <= 0 {
		return 0
	}
	return FastInt() % n
}
