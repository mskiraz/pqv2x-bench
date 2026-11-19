package dilithium

// nttAVX2_8 performs forward NTT stages len=128 down to len=8 for non-amd64 architectures.
func nttAVX2_8(a *[256]int32, zetas *[256]int32) {
	k := 0
	for len := 128; len >= 8; len >>= 1 {
		start := 0
		for start < 256 {
			k++
			zeta := zetas[k]
			for j := start; j < start+len; j++ {
				t := montgomeryReduce(int64(zeta) * int64(a[j+len]))
				a[j+len] = a[j] - t
				a[j] += t
			}
			start += 2 * len
		}
	}
}
