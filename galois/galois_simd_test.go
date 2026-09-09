//go:build amd64 || arm64
// +build amd64 arm64

package galois

import (
	"fmt"
	"testing"
)

func Benchmark_mulVect_simd(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulVect_simd(benchMulC, i, o)
			}
		})
	}
}
func Benchmark_mulXorVect_simd(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulXorVect_simd(benchMulC, i, o)
			}
		})
	}
}
func Benchmark_xorVect_simd(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				xorVect_simd(i, o)
			}
		})
	}
}
func Benchmark_lastIndex_simd(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			s := make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				lastIndex_simd(s, searchByte)
			}
		})
	}
}
