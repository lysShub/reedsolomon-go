package reedsolomon

import (
	"math"
	"math/rand"
	"testing"
)

func Test_ReedsolomonPL(t *testing.T) {
	if act := ReedsolomonPL(8, 4, -1); act != 0.0 {
		t.Fatalf("act = %v, want 0.0", act)
	}
	if act := ReedsolomonPL(8, 4, 0); act != 0.0 {
		t.Fatalf("act = %v, want 0.0", act)
	}
	if act := ReedsolomonPL(8, 4, 1); act != 1.0 {
		t.Fatalf("act = %v, want 1.0", act)
	}
	if act := ReedsolomonPL(8, 4, 1.1); act != 1.0 {
		t.Fatalf("act = %v, want 1.0", act)
	}

	for groupsize := 2; groupsize < 8; groupsize++ {
		for paritysize := 1; paritysize < groupsize; paritysize++ {
			for _, pl := range []float64{0.01, 0.1, 0.3, 0.6} {

				act := ReedsolomonPL(groupsize, groupsize-paritysize, pl)
				exp := mock(groupsize, paritysize, pl)
				if math.Abs(exp-act) > 0.005 {
					t.Fatalf("exp = %v, act = %v", exp, act)
				}

			}
		}
	}
}

func Benchmark_ReedsolomonPL(b *testing.B) {
	for b.Loop() {
		_ = ReedsolomonPL(8, 4, 0.1)
	}
}

func mock(groupsize int, paritysize int, pl float64) float64 {
	datasize := groupsize - paritysize

	var lossed = func() bool {
		return rand.Int63n(1e6) < int64(pl*1e6)
	}

	var recv int
	const loops = 1e6
	for i := 0; i < loops; i++ {
		var g []bool
		var l int
		for i := 0; i < groupsize; i++ {
			if lossed() {
				g = append(g, false)
				l += 1
			} else {
				g = append(g, true)
			}
		}

		if l <= paritysize {
			// can recover all data blocks
			recv += datasize
		} else {
			// cannot recover, count original data
			for _, e := range g[:datasize] {
				if e {
					recv += 1
				}
			}
		}
	}
	var send = datasize * loops

	return 1 - float64(recv)/float64(send)
}
