//go:build amd64 || arm64
// +build amd64 arm64

package galois

import (
	"testing"
)

func Test_Vect(t *testing.T) {
	t.Run("SIMD", func(t *testing.T) {
		defer initialize()
		mulVect = mulVect_simd
		mulXorVect = mulXorVect_simd
		lastIndex = lastIndex_simd
		xorVect = xorVect_simd

		test_vect(t)
	})

	t.Run("GO", func(t *testing.T) {
		defer initialize()
		mulVect = mulVect_go
		mulXorVect = mulXorVect_go
		lastIndex = lastIndex_go
		xorVect = xorVect_go

		test_vect(t)
	})
}

func Benchmark_Vect_GO(b *testing.B) {
	defer initialize()
	mulVect = mulVect_go
	mulXorVect = mulXorVect_go
	lastIndex = lastIndex_go
	xorVect = xorVect_go

	bench_vect(b)
}

func Benchmark_Vect_SIMD(b *testing.B) {
	defer initialize()
	mulVect = mulVect_simd
	mulXorVect = mulXorVect_simd
	lastIndex = lastIndex_simd
	xorVect = xorVect_simd

	bench_vect(b)
}
