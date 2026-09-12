package reedsolomon_test

import (
	"bytes"
	"crypto/rand"
	"slices"
	"testing"

	"github.com/lysShub/reedsolomon-go"
)

func makes(rows, cols int) (v [][]byte) {
	for i := 0; i < rows; i++ {
		v = append(v, make([]byte, cols))
	}
	return v
}
func deletes[T any](s []T, idxs ...int) []T {
	slices.Sort(idxs)
	for j := len(idxs) - 1; j >= 0; j-- {
		i := idxs[j]
		s = slices.Delete(s, i, i+1)
	}
	return s
}

func Test_Base(t *testing.T) {
	t.Run("{2, 3} one data block lost", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 3, Datasize: 2}

		// data: first 2 data blocks, last 1 parity
		var datas = make([][]byte, 2)
		datas[0] = []byte{1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		datas[1] = []byte{4, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		var parity = makes(1, len(datas[0]))

		{ // encode
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 1) // lose datas[1]

		{ // reconst
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{0, 2}, reconst)
			if n != 1 {
				t.Fatalf("n = %d, want 1", n)
			}
			if !bytes.Equal(datas[1], reconst[0]) {
				t.Fatalf("reconst mismatch")
			}
		}
	})

	t.Run("{3, 5} one data block lost", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 5, Datasize: 3}

		var datas = make([][]byte, 3)
		datas[0] = []byte{1, 2, 3, 4}
		datas[1] = []byte{4, 3, 0, 0}
		datas[2] = []byte{2, 1, 6, 9}
		var parity = makes(2, len(datas[0]))

		{ // encode
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
			reedsolomon.Encode(para, datas[2], 2, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 1) // lose datas[1]

		{ // reconst
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{1, 2, 3, 4}, reconst)
			if n != 1 {
				t.Fatalf("n = %d, want 1", n)
			}
			if !bytes.Equal(datas[1], reconst[0]) {
				t.Fatalf("reconst mismatch")
			}
		}
	})

	t.Run("{3, 5} one data and one parity lost", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 5, Datasize: 3}

		var datas = make([][]byte, 3)
		datas[0] = []byte{1, 2, 3, 4}
		datas[1] = []byte{4, 3, 0, 0}
		datas[2] = []byte{2, 1, 6, 9}
		var parity = makes(2, len(datas[0]))

		{ // encode
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
			reedsolomon.Encode(para, datas[2], 2, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 1, 3)

		{ // reconst
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{0, 2, 4}, reconst)
			if n != 1 {
				t.Fatalf("n = %d, want 1", n)
			}
			if !bytes.Equal(datas[1], reconst[0]) {
				t.Fatalf("reconst mismatch")
			}
		}
	})

	t.Run("{3, 5} two parity blocks lost", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 5, Datasize: 3}

		var datas = make([][]byte, 3)
		datas[0] = []byte{1, 2, 3, 4}
		datas[1] = []byte{4, 3, 0, 0}
		datas[2] = []byte{2, 1, 6, 9}
		var parity = makes(2, len(datas[0]))

		{ // encode
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
			reedsolomon.Encode(para, datas[2], 2, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 3, 4)

		{ // reconst
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{0, 1, 2}, reconst)
			if n != 0 {
				t.Fatalf("n = %d, want 0", n)
			}
		}
	})

	t.Run("datasize ==  1", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 3, Datasize: 1}

		// data: first 1 data block, last 2 parity
		var datas = make([][]byte, 1)
		datas[0] = []byte{1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		var parity = makes(2, len(datas[0]))

		{ // encode
			reedsolomon.Encode(para, datas[0], 0, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 0, 2) // lose datas[0,2]

		{ // reconst
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{1}, reconst)
			if n != 1 {
				t.Fatalf("n = %d, want 1", n)
			}
			if !bytes.Equal(datas[0], reconst[0]) {
				t.Fatalf("reconst mismatch")
			}
		}
	})
}

func Test_NotAlign(t *testing.T) {
	// test encode/reconstruct with unaligned block lengths

	var buildBlocks = func(datas [][]byte, parity [][]byte) [][]byte {
		var blocks [][]byte

		size := 0
		for _, e := range datas {
			blocks = append(blocks, slices.Clone(e))
			size = max(size, len(e))
		}
		for _, e := range parity {
			// parity length = longest data length
			blocks = append(blocks, slices.Clone(e[:size]))
		}
		return blocks
	}

	t.Run("{2, 3} block length mismatch 1", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 3, Datasize: 2}
		var datas = make([][]byte, 2)
		datas[0] = []byte{1, 2, 3, 4}
		datas[1] = []byte{4, 3}

		var parity = makes(0xff, 0xff)
		{
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
		}

		var blocks = buildBlocks(datas, parity[:para.Paritysize()])
		blocks = deletes(blocks, 1)

		{
			var reconst = makes(0xff, 0xff)
			n := reedsolomon.Reconst(para, blocks, []uint8{0, 2}, reconst)
			if n != 1 {
				t.Fatalf("n = %d, want 1", n)
			}
			if !bytes.HasPrefix(reconst[0], datas[1]) {
				t.Fatalf("reconst mismatch")
			}
		}
	})

	t.Run("{2, 3} block length mismatch 2", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 3, Datasize: 2}
		var datas = make([][]byte, 2)
		datas[0] = []byte{1, 2, 3, 4}
		datas[1] = []byte{4, 3}

		var parity = makes(0xff, 0xff)
		{
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
		}

		var blocks = buildBlocks(datas, parity[:para.Paritysize()])
		blocks = deletes(blocks, 0)

		{
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{1, 2}, reconst)
			if n != 1 {
				t.Fatalf("n = %d, want 1", n)
			}
			if !bytes.HasPrefix(reconst[0], datas[0]) {
				t.Fatalf("reconst mismatch")
			}
		}
	})

	t.Run("{2, 3} parity matrix nonzero, datas[0] not longest", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 3, Datasize: 2}

		var datas = make([][]byte, 2)
		datas[0] = []byte{4, 3}
		datas[1] = []byte{1, 2, 3, 4}

		var parity = makes(0xff, 0xff)
		for _, e := range parity {
			rand.Read(e)
		}
		{
			// require parity[i][len(datas[0]):] zero
			for _, e := range parity {
				clear(e[len(datas[0]):])
			}

			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
		}

		var blocks = buildBlocks(datas, parity[:para.Paritysize()])
		blocks = deletes(blocks, 1)

		{
			var reconst = makes(0xff, 0xff)
			for _, e := range reconst {
				rand.Read(e)
			}
			// require reconst[i][len(blocks[0]):] zero
			for _, e := range reconst {
				clear(e[len(blocks[0]):])
			}

			n := reedsolomon.Reconst(para, blocks, []uint8{0, 2}, reconst)
			if n != 1 {
				t.Fatalf("n = %d, want 1", n)
			}
			if !bytes.HasPrefix(reconst[0], datas[1]) {
				t.Fatalf("reconst mismatch")
			}
		}
	})

}
