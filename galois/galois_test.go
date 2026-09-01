package galois

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
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
			require.Equal(t, Mul(a, b), Mul(a, b))
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

		require.Equal(t, exp, out)
		if n > 0 {
			require.Equal(t, Mul(c, in[0]), out[0])
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

		require.Equal(t, exp, out)
		if n > 0 {
			require.Equal(t, Mul(c, in[0])^o0, out[0])
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

		require.Equal(t, exp, out)
	}
}

func Test_lastIndex(t *testing.T) {
	require.Equal(t, -1, LastIndex(nil, 0xff))

	for n := 0; n <= 0xff; n++ {
		for i := 0; i < +1; i++ {
			s := make([]byte, n)
			if i < len(s) {
				s[i] = 0xff
				require.Equal(t, i, LastIndex(s, 0xff))
			} else {
				require.Equal(t, -1, LastIndex(s, 0xff))
			}
		}
	}
}
