package galois

import (
	"bytes"
	"crypto/rand"
	"testing"
)

var testSize = []int{
	0, 1, 3, 5, 15, 16, 17, 31, 32,
	33, 63, 64, 65, 127, 128, 129, 1496, 1500,
}

func Test_Mul(t *testing.T) {
	var d = make([]byte, 0xff)
	rand.Read(d)
	for _, a := range d {
		for _, b := range d {
			if Mul(a, b) != Mul(a, b) {
				t.Fatalf("a=%d b=%d", a, b)
			}
		}
	}
}

func Test_MulVect(t *testing.T) {
	for _, n := range testSize {
		var (
			c   = byte(123)
			in  = make([]byte, n)
			out = make([]byte, n)
			exp = make([]byte, n)
		)
		rand.Read(in)
		mulVect_go(c, in, exp)

		MulVect(c, in, out)

		if !bytes.Equal(exp, out) {
			t.Fatalf("n=%d mismatch", n)
		}
		if n > 0 {
			if Mul(c, in[0]) != out[0] {
				t.Fatalf("n=%d mismatch", n)
			}
		}
	}
}

func Test_MulXorVect(t *testing.T) {
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
		o0 := byte(0)
		if n > 0 {
			o0 = out[0]
		}

		MulXorVect(c, in, out)

		if !bytes.Equal(exp, out) {
			t.Fatalf("n=%d mismatch", n)
		}
		if n > 0 {
			if Mul(c, in[0])^o0 != out[0] {
				t.Fatalf("n=%d mismatch", n)
			}
		}
	}
}

func Test_xorVect(t *testing.T) {
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

		XorVect(in, out)

		if !bytes.Equal(exp, out) {
			t.Fatalf("n=%d mismatch", n)
		}
	}
}

func Test_lastIndex(t *testing.T) {
	if LastIndex(nil, 0xff) != -1 {
		t.Fatalf("LastIndex(nil) != -1")
	}

	for n := 0; n <= 0xff; n++ {
		for i := 0; i < +1; i++ {
			s := make([]byte, n)
			if i < len(s) {
				s[i] = 0xff
				if LastIndex(s, 0xff) != i {
					t.Fatalf("n=%d i=%d", n, i)
				}
			} else {
				if LastIndex(s, 0xff) != -1 {
					t.Fatalf("n=%d i=%d", n, i)
				}
			}
		}
	}
}
