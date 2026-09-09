// reedsolomon encode/reconst
package reedsolomon

import (
	"slices"

	"github.com/lysShub/debug-go"
	"github.com/lysShub/reedsolomon-go/galois"
)

type Para struct {
	Groupsize uint8 // total blocks in a group, = data + parity
	Datasize  uint8 // data  blocks in a group
}

func (p Para) Paritysize() uint8 { return p.Groupsize - p.Datasize }
func (p Para) Valid() bool {
	return p.Groupsize > 0 && p.Datasize > 0 && p.Groupsize >= p.Datasize
}

// Encode encoding the idx-th data-block, idx must be 0 on first call.
//
//	para  : reed-solomon encodec parameter
//	data  : idx-th data-block in the group
//	idx   : position of this data-block in the group
//	parity: parity blocks, len(parity) >= para.Paritysize and len(parity[i]) >= len(data)
func Encode(para Para, data []byte, idx uint8, parity [][]byte) {
	if debug.Debug() {
		debug.Less(idx, para.Datasize)
		debug.GreaterOrEqual(len(parity), int(para.Paritysize()))
		debug.GreaterOrEqual(slices.Min(lens(parity[:para.Paritysize()])), len(data))
	}
	parity = parity[:para.Paritysize()]

	matrix := Cache.matrix(para)
	for i, p := range parity {
		// row: i  col: idx
		c := matrix[int(i)*int(para.Datasize)+int(idx)]

		if idx == 0 {
			galois.MulVect(c, data, p)
		} else {
			galois.MulXorVect(c, data, p)
		}
	}
}

// Reconst reconstruct lost data-blocks.
//
//	para   : reed-solomon encodec parameter
//	blocks : received blocks, data and parity block, require are in-ordered
//	indexs : index of each block
//	reconst: reconstructed data-blocks, len(reconst) >= para.Datasize-len(blocks) and len(reconst[i]) >= len(blocks[j])
func Reconst(para Para, blocks [][]byte, indexs []uint8, reconst [][]byte) int {
	if para.Groupsize == para.Datasize {
		return 0 // no rs encoding
	}
	if len(blocks) < int(para.Datasize) {
		return 0 // too many blocks lost
	}
	if debug.Debug() {
		debug.Equal(len(blocks), len(indexs))
		debug.Equal(len(indexs), len(slices.Compact(indexs)))
		debug.Less(slices.Max(indexs), para.Groupsize)
		debug.True(slices.IsSorted(indexs))
	}
	blocks = blocks[:para.Datasize]
	indexs = indexs[:para.Datasize]
	loss := lossDatablocks(para, indexs)
	if loss == 0 {
		return 0 // no lost data block
	}
	if debug.Debug() {
		debug.GreaterOrEqual(len(reconst), loss)
		debug.GreaterOrEqual(slices.Min(lens(reconst)), slices.Max(lens(blocks)))
	}
	reconst = reconst[:loss]

	matrix := Cache.matrix(para, indexs...)
	for bi, b := range blocks {
		for ri, r := range reconst {
			// row: ri  col: bi
			c := matrix[int(ri)*int(para.Datasize)+int(bi)]

			if bi == 0 {
				galois.MulVect(c, b, r)
			} else {
				galois.MulXorVect(c, b, r)
			}
		}
	}
	return loss
}

func lens(s [][]byte) (lens []int) {
	for _, e := range s {
		lens = append(lens, len(e))
	}
	return lens
}

// lossDatablocks counts lost data blocks
func lossDatablocks(p Para, index []uint8) int {
	n := 0
	for _, e := range index {
		if e < p.Datasize {
			n += 1
		}
	}
	return int(p.Datasize) - n
}
