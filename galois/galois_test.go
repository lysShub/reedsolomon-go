package galois

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
)

var benchSizes = []int{
	31, 32,
	64, 65,
	128,
	256,
	1492,
	1500,
}
var testSize = []int{
	0, 1, 3, 5, 15, 16, 17, 31, 32,
	33, 63, 64, 65, 127, 128, 129, 1496, 1500,
}

const benchMulC byte = 0b10101010
const searchByte byte = 0b10101010

func test_vect(t *testing.T) {
	for _, n := range testSize {

		t.Run(fmt.Sprintf("MulVect_%d", n), func(t *testing.T) {
			c := byte(123)
			in := make([]byte, n)
			out := make([]byte, n)
			rand.Read(in)
			MulVect(c, in, out)
			for i := range in {
				if out[i] != Mul(c, in[i]) {
					t.Fatalf("n=%d i=%d mismatch", n, i)
				}
			}
		})

		t.Run(fmt.Sprintf("MulXorVect_%d", n), func(t *testing.T) {
			c := byte(123)
			in := make([]byte, n)
			out := make([]byte, n)
			rand.Read(in)
			rand.Read(out)
			exp := append([]byte(nil), out...)
			MulXorVect(c, in, out)
			for i := range in {
				if out[i] != exp[i]^Mul(c, in[i]) {
					t.Fatalf("n=%d i=%d mismatch", n, i)
				}
			}
		})

		t.Run(fmt.Sprintf("XorVect_%d", n), func(t *testing.T) {
			in := make([]byte, n)
			out := make([]byte, n)
			rand.Read(in)
			rand.Read(out)
			exp := append([]byte(nil), out...)
			XorVect(in, out)
			for i := range in {
				if out[i] != exp[i]^in[i] {
					t.Fatalf("n=%d i=%d mismatch", n, i)
				}
			}
		})

		t.Run(fmt.Sprintf("LastIndex_%d", n), func(t *testing.T) {
			s := make([]byte, n)
			if LastIndex(s, 0xff) != bytes.LastIndexByte(s, 0xff) {
				t.Fatalf("n=%d mismatch", n)
			}
			for i := 0; i < n; i++ {
				s[i] = 0xff
				if LastIndex(s, 0xff) != bytes.LastIndexByte(s, 0xff) {
					t.Fatalf("n=%d i=%d mismatch", n, i)
				}
			}
		})
	}
}

func bench_vect(b *testing.B) {
	for _, n := range benchSizes {

		b.Run(fmt.Sprintf("MulVect_%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				MulVect(benchMulC, i, o)
			}
		})

		b.Run(fmt.Sprintf("MulXorVect_%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				MulXorVect(benchMulC, i, o)
			}
		})

		b.Run(fmt.Sprintf("XorVect_%d", n), func(b *testing.B) {
			i, o := make([]byte, n), make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				XorVect(i, o)
			}
		})

		b.Run(fmt.Sprintf("LastIndex_%d", n), func(b *testing.B) {
			s := make([]byte, n)
			b.SetBytes(int64(n))
			b.ResetTimer()
			for range b.N {
				LastIndex(s, searchByte)
			}
		})

	}
}

func Test_Scal(t *testing.T) {
	t.Run("Mul", func(t *testing.T) {
		for a := 0; a <= 0xff; a++ {
			for b := 0; b <= 0xff; b++ {
				if Mul(byte(a), byte(b)) != mulTable[byte(a)][byte(b)] {
					t.Fatalf("a=%d b=%d", a, b)
				}
			}
		}
	})
	t.Run("Add", func(t *testing.T) {
		for a := 0; a <= 0xff; a++ {
			for b := 0; b <= 0xff; b++ {
				if Add(byte(a), byte(b)) != byte(a)^byte(b) {
					t.Fatalf("a=%d b=%d", a, b)
				}
			}
		}
	})
	t.Run("Div", func(t *testing.T) {
		for a := 0; a <= 0xff; a++ {
			for b := 1; b <= 0xff; b++ {
				if Mul(Div(byte(a), byte(b)), byte(b)) != byte(a) {
					t.Fatalf("a=%d b=%d", a, b)
				}
			}
		}
	})
	t.Run("Pow", func(t *testing.T) {
		for a := 0; a <= 0xff; a++ {
			r := byte(1)
			for n := 0; n <= 0xff; n++ {
				if Pow(byte(a), byte(n)) != r {
					t.Fatalf("a=%d n=%d", a, n)
				}
				r = Mul(r, byte(a))
			}
		}
	})
	t.Run("lastIndex", func(t *testing.T) {
		if lastIndex_go(nil, 0xff) != -1 {
			t.Fatal("lastIndex_go(nil) != -1")
		}
		for n := 0; n <= 0xff; n++ {
			for i := 0; i < n; i++ {
				s := make([]byte, n)
				s[i] = 0xff
				if lastIndex_go(s, 0xff) != i {
					t.Fatalf("n=%d i=%d", n, i)
				}
			}
		}
	})
}
