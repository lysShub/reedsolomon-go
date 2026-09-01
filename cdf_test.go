package reedsolomon

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ReedsolomonPL(t *testing.T) {
	require.Equal(t, 0.0, ReedsolomonPL(8, 4, -1))
	require.Equal(t, 0.0, ReedsolomonPL(8, 4, 0))
	require.Equal(t, 1.0, ReedsolomonPL(8, 4, 1))
	require.Equal(t, 1.0, ReedsolomonPL(8, 4, 1.1))

	for groupsize := 2; groupsize < 8; groupsize++ {
		for paritysize := 1; paritysize < groupsize; paritysize++ {
			for _, pl := range []float64{0.01, 0.1, 0.3, 0.6} {

				act := ReedsolomonPL(groupsize, groupsize-paritysize, pl)
				exp := mock(groupsize, paritysize, pl)
				require.InDelta(t, exp, act, 0.005)

			}
		}
	}
}

func Benchmark_ReedsolomonPL(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
			// 可以恢复所有datasize
			recv += datasize
		} else {
			// 无法恢复, 统计原始数据
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
