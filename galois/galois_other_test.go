//go:build !amd64 && !arm64
// +build !amd64,!arm64

package galois

import (
	"testing"
)

func Test_Vect_GO(t *testing.T) {
	mulVect = mulVect_go
	mulXorVect = mulXorVect_go
	lastIndex = lastIndex_go
	xorVect = xorVect_go

	test_vect(t)
}

func Benchmark_Vect_GO(b *testing.B) {
	mulVect = mulVect_go
	mulXorVect = mulXorVect_go
	lastIndex = lastIndex_go
	xorVect = xorVect_go

	bench_vect(b)
}
