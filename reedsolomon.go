// Package reedsolomon encode/reconst
package reedsolomon

import (
	"acceler/pkg/debug"
	"acceler/pkg/reedsolomon/galois"
	"slices"
)

type Para struct {
	Groupsize uint8 // 组中所有块数, = 数据块数 + 校验块数
	Datasize  uint8 // 组中数据块数
}

func (p Para) Paritysize() uint8 { return p.Groupsize - p.Datasize }
func (p Para) Valid() bool {
	return p.Groupsize > 0 && p.Datasize > 0 && p.Groupsize >= p.Datasize
}

// Encode 对第idx个数据块进行编码, 首次编码idx必须为0
//
//	para  : reedsolomon编码参数
//	data  : 组中的第idx个数据包
//	idx   : 此数据包在组中的位置
//	parity: 存放校验数据块, 要求 len(parity) >= para.Paritysize, 且 len(parity[i]) >= len(data)
func Encode(para Para, data []byte, idx uint8, parity [][]byte) {
	if debug.Debug() {
		debug.Less(idx, para.Datasize)
		debug.GreaterOrEqual(len(parity), int(para.Paritysize()))
		debug.GreaterOrEqual(slices.Min(lens(parity[:para.Paritysize()])), len(data))
	}
	parity = parity[:para.Paritysize()]

	matrix := Cache.matrix(para)
	for i, p := range parity {
		c := matrix.Row(i)[idx]

		if idx == 0 {
			galois.MulVect(c, data, p)
		} else {
			galois.MulXorVect(c, data, p)
		}
	}
}

// Reconst 恢复丢失的数据块
//
//	para   : 当前reedsolomon编码参数
//	blocks : 接收到的编码块, 包含datablock 和 parityblock, 要求是顺序的
//	indexs : 各编码块对应的位置
//	reconst: 存放被恢复的数据块, 要求 len(reconst) >= para.Datasize-len(blocks) 且, len(reconst[i]) >= len(blocks[j])
func Reconst(para Para, blocks [][]byte, indexs []uint8, reconst [][]byte) int {
	if para.Groupsize == para.Datasize {
		return 0 // 没有进行rs编码
	}
	if len(blocks) < int(para.Datasize) {
		return 0 // 丢失太多编码块, 无法恢复
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
		return 0 // 无丢失数据块
	}
	if debug.Debug() {
		debug.GreaterOrEqual(len(reconst), loss)
		debug.GreaterOrEqual(slices.Min(lens(reconst)), slices.Max(lens(blocks)))
	}
	reconst = reconst[:loss]

	matrix := Cache.matrix(para, indexs...)
	for bi, b := range blocks {
		for ri, r := range reconst {
			c := matrix.Row(ri)[bi]

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

// lossDatablocks 统计丢失数据包个数
func lossDatablocks(p Para, index []uint8) int {
	n := 0
	for _, e := range index {
		if e < p.Datasize {
			n += 1
		}
	}
	return int(p.Datasize) - n
}
