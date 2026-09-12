//go:build amd64
// +build amd64

package galois

import (
	"testing"

	"golang.org/x/sys/cpu"
)

func Test_Vect_GFNI(t *testing.T) {
	if !cpu.X86.HasAVX512GFNI {
		t.Skip("AVX512GFNI not supported")
	}
	defer initialize()
	mulVect = mulVect_gfni
	mulXorVect = mulXorVect_gfni
	lastIndex = lastIndex_simd
	xorVect = xorVect_simd

	test_vect(t)
}

func Benchmark_Vect_GFNI(b *testing.B) {
	if !cpu.X86.HasAVX512GFNI {
		b.Skip("AVX512GFNI not supported")
	}
	defer initialize()
	mulVect = mulVect_gfni
	mulXorVect = mulXorVect_gfni
	lastIndex = lastIndex_simd
	xorVect = xorVect_simd

	bench_vect(b)
}
