//go:build amd64
// +build amd64

package galois

import (
	"fmt"
	"testing"

	"golang.org/x/sys/cpu"
)

func Test_Vect(t *testing.T) {
	t.Run("simd", func(t *testing.T) {
		test_vect(t,
			mulVect_simd,
			mulXorVect_simd,
			lastIndex_simd,
			xorVect_simd,
		)
	})
	t.Run("gfni", func(t *testing.T) {
		if !cpu.X86.HasAVX512GFNI {
			t.Skip("AVX512GFNI not supported")
		}
		test_vect(t,
			mulVect_gfni,
			mulXorVect_gfni,
			lastIndex_simd,
			xorVect_simd,
		)
	})
}

func Benchmark_mulVect_gfni(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			if !cpu.X86.HasAVX512GFNI {
				b.Skip("AVX512GFNI not supported")
			}
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulVect_gfni(benchMulC, i, o)
			}
		})
	}
}
