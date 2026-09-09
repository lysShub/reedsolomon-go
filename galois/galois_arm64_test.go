//go:build arm64
// +build arm64

package galois

import "testing"

func Test_Vect(t *testing.T) {
	t.Run("simd", func(t *testing.T) {
		test_vect(t,
			mulVect_simd,
			mulXorVect_simd,
			lastIndex_simd,
			xorVect_simd,
		)
	})
}
