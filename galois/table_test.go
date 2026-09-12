package galois

import (
	"crypto/sha256"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_new_Table(t *testing.T) {
	{
		ta := mulTable()

		var b = make([]byte, 0, 0xffff)
		for _, e := range ta {
			b = append(b, e[:]...)
		}

		exp := [32]byte{0, 61, 26, 96, 151, 131, 210, 116, 11, 155, 63, 0, 176, 205, 158, 67, 228, 44, 79, 62, 237, 197, 255, 84, 236, 23, 9, 153, 109, 82, 225, 224}
		act := sha256.Sum256(b)
		require.Equal(t, exp, act)
	}
	{
		lh := lohiTable()

		var b = make([]byte, 0, 8192)
		for _, e := range lh {
			b = append(b, e.lo[:]...)
			b = append(b, e.hi[:]...)
		}

		exp := [32]byte{118, 220, 196, 252, 39, 178, 191, 152, 246, 192, 245, 14, 113, 87, 1, 1, 246, 28, 194, 172, 152, 60, 54, 111, 109, 128, 214, 75, 101, 90, 47, 34}
		act := sha256.Sum256(b)
		require.Equal(t, exp, act)
	}
}

func Test_mul_lowhigh_cross_verify(t *testing.T) {
	// galois mul
	//
	// for i, e := range in {
	// 	  out[i] = mul[e]
	// 	  out[i] = lo[e&0x0f] ^ hi[e>>4]
	// }

	seed := time.Now().UnixNano()
	t.Logf("seed: %d", seed)
	r := rand.New(rand.NewSource(seed))

	for range 0xff {
		c := byte(r.Int31n(0xff))
		var in = make([]byte, r.Int31n(1536)+1)
		r.Read(in)

		var out1 = make([]byte, len(in))
		mul := mulTable()[c]
		for i, e := range in {
			out1[i] = mul[e]
		}

		var out2 = make([]byte, len(in))
		lohimk := lohiTable()[c]
		for i, e := range in {
			out2[i] = lohimk.lo[e&0x0f] ^ lohimk.hi[e>>4]
		}

		require.Equal(t, out1, out2)
	}
}
