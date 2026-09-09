package galois

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
	"time"
)

func Test_new_Table(t *testing.T) {
	{
		ta := &mulTable

		var b = make([]byte, 0, 0xffff)
		for _, e := range ta {
			b = append(b, e[:]...)
		}

		exp := [32]byte{20, 161, 231, 231, 124, 168, 163, 11, 91, 181, 62, 99, 16, 116, 140, 224, 73, 142, 185, 224, 74, 183, 138, 68, 219, 239, 182, 235, 250, 200, 168, 75}
		act := sha256.Sum256(b)
		if exp != act {
			t.Fatalf("sha256 mismatch")
		}
	}
	{
		lh := &lohiTable

		var b = make([]byte, 0, 8192)
		for _, e := range lh {
			b = append(b, e.lo[:]...)
			b = append(b, e.hi[:]...)
		}

		exp := [32]byte{184, 32, 66, 66, 250, 216, 134, 44, 232, 30, 221, 162, 147, 133, 189, 197, 92, 76, 35, 112, 253, 64, 163, 61, 96, 187, 148, 151, 230, 165, 47, 246}
		act := sha256.Sum256(b)
		if exp != act {
			t.Fatalf("sha256 mismatch")
		}
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
		mul := mulTable[c]
		for i, e := range in {
			out1[i] = mul[e]
		}

		var out2 = make([]byte, len(in))
		lohimk := lohiTable[c]
		for i, e := range in {
			out2[i] = lohimk.lo[e&0x0f] ^ lohimk.hi[e>>4]
		}

		if !bytes.Equal(out1, out2) {
			t.Fatalf("mul mismatch")
		}
	}
}
