//go:build arm64
// +build arm64

package galois

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
)

func Test_Base(t *testing.T) {
	var crossSizes = []int{0, 1, 15, 16, 17, 31, 63, 64, 65, 127, 128, 129, 255, 256, 257, 511, 1500}

	t.Run("mux", func(t *testing.T) {
		for c := 0; c <= 0xff; c++ {
			for _, n := range crossSizes {
				i := make([]byte, n)
				rand.Read(i)
				oGo := make([]byte, n)
				oSmid := make([]byte, n)
				mulVect_go(byte(c), i, oGo)
				mulVect_smid(byte(c), i, oSmid)
				if !bytes.Equal(oGo, oSmid) {
					t.Fatalf("c=%d n=%d", c, n)
				}
			}
		}
	})
	t.Run("mulXor", func(t *testing.T) {
		for c := 0; c <= 0xff; c++ {
			for _, n := range crossSizes {
				i := make([]byte, n)
				o := make([]byte, n)
				rand.Read(i)
				rand.Read(o)
				oGo := make([]byte, n)
				oSmid := make([]byte, n)
				copy(oGo, o)
				copy(oSmid, o)
				mulXorVect_go(byte(c), i, oGo)
				mulXorVect_smid(byte(c), i, oSmid)
				if !bytes.Equal(oGo, oSmid) {
					t.Fatalf("c=%d n=%d", c, n)
				}
			}
		}
	})
}

var benchSizes = []int{
	31, 32,
	64, 65,
	128,
	256,
	1492,
	1500,
}

const benchMulC byte = 0b10101010
const searchByte byte = 0b10101010

func Benchmark_mulVect(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("go/%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulVect_go(benchMulC, i, o)
			}
		})

		b.Run(fmt.Sprintf("smid/%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulVect_smid(benchMulC, i, o)
			}
		})
		fmt.Println()
	}
	fmt.Println()
}

func Benchmark_mulXorVect(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("go/%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulXorVect_go(benchMulC, i, o)
			}
		})

		b.Run(fmt.Sprintf("smid/%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				mulXorVect_smid(benchMulC, i, o)
			}
		})
		fmt.Println()
	}
	fmt.Println()
}

func Benchmark_xorVect(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("go/%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				xorVect_go(i, o)
			}
		})

		b.Run(fmt.Sprintf("smid16/%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				xorVect_smid(i, o)
			}
		})
		fmt.Println()
	}
	fmt.Println()
}

func Benchmark_lastIndex(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("go/%d", n), func(b *testing.B) {
			s := make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				lastIndex_go(s, searchByte)
			}
		})

		b.Run(fmt.Sprintf("smid16/%d", n), func(b *testing.B) {
			s := make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				lastIndex_smid(s, searchByte)
			}
		})
		fmt.Println()
	}
	fmt.Println()
}
