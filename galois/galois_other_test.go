//go:build !amd64 && !arm64
// +build !amd64,!arm64

package galois

import (
	"fmt"
	"testing"
)

func Test_Vect(t *testing.T) {
	t.Run("go", func(t *testing.T) {
		test_vect(t,
			mulVect_go,
			mulXorVect_go,
			lastIndex_go,
			xorVect_go,
		)
	})
}

func Benchmark_mulVect_go(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulVect_go(benchMulC, i, o)
			}
		})
	}
}
func Benchmark_mulXorVect_go(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulXorVect_go(benchMulC, i, o)
			}
		})
	}
}
func Benchmark_xorVect_go(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				xorVect_go(i, o)
			}
		})
	}
}
func Benchmark_lastIndex_go(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			s := make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				lastIndex_go(s, searchByte)
			}
		})
	}
}
