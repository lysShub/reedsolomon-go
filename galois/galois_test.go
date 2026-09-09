package galois

import (
	"bytes"
	"crypto/rand"
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

func test_vect(t *testing.T,
	mulVect func(c byte, i, o []byte),
	mulXorVect func(c byte, i, o []byte),
	lastIndex func(s []byte, v byte) int,
	xorVect func(i, o []byte),
) {
	t.Run("mulVect", func(t *testing.T) {
		for _, n := range testSize {
			var (
				c   = byte(123)
				in  = make([]byte, n)
				out = make([]byte, n)
				exp = make([]byte, n)
			)
			rand.Read(in)
			mulVect_go(c, in, exp)
			mulVect(c, in, out)
			if !bytes.Equal(exp, out) {
				t.Fatalf("n=%d mismatch", n)
			}
		}
	})
	t.Run("mulXorVect", func(t *testing.T) {
		for _, n := range testSize {
			var (
				c   = byte(123)
				in  = make([]byte, n)
				out = make([]byte, n)
				exp = make([]byte, n)
			)
			rand.Read(in)
			rand.Read(out)
			copy(exp, out)
			mulXorVect_go(c, in, exp)
			mulXorVect(c, in, out)
			if !bytes.Equal(exp, out) {
				t.Fatalf("n=%d mismatch", n)
			}
		}
	})
	t.Run("xorVect", func(t *testing.T) {
		for _, n := range testSize {
			var (
				in  = make([]byte, n)
				out = make([]byte, n)
				exp = make([]byte, n)
			)
			rand.Read(in)
			rand.Read(out)
			copy(exp, out)
			xorVect_go(in, exp)
			xorVect(in, out)
			if !bytes.Equal(exp, out) {
				t.Fatalf("n=%d mismatch", n)
			}
		}
	})
	t.Run("lastIndex", func(t *testing.T) {
		for _, n := range testSize {
			s := make([]byte, n)
			if lastIndex(s, 0xff) != -1 {
				t.Fatalf("n=%d mismatch", n)
			}
			for i := 0; i < n; i++ {
				s[i] = 0xff
				if lastIndex(s, 0xff) != lastIndex_go(s, 0xff) {
					t.Fatalf("n=%d i=%d mismatch", n, i)
				}
			}
		}
	})
}

func Test_Scal(t *testing.T) {
	t.Run("Mul", func(t *testing.T) {
		for a := 0; a <= 0xff; a++ {
			for b := 0; b <= 0xff; b++ {
				if Mul(byte(a), byte(b)) != mulTable()[byte(a)][byte(b)] {
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
