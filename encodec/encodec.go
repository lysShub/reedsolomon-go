package encodec

import (
	"slices"

	"github.com/lysShub/debug-go"
)

// Encodec 获取编解码矩阵, idxs为空表示获取编码矩阵
func Encodec(grousize, datasize uint8, idxs ...uint8) Matrix {
	m := encodec(grousize, datasize, idxs...)

	if b := m.raw(); len(b)*2 <= cap(b) {
		defer m.Release()
		return m.Clone()
	} else {
		return m
	}
}
func encodec(grousize, datasize uint8, idxs ...uint8) Matrix {
	if debug.Debug() {
		debug.Greater(grousize, 0)
		debug.Greater(datasize, 0)
		debug.Less(datasize, grousize) // 存在rs编码
	}
	if len(idxs) == 0 {
		return encodeMatrix(grousize, datasize)
	} else {
		if debug.Debug() {
			debug.Greater(lossDatablocks(datasize, idxs), 0) // 存在丢包
			debug.GreaterOrEqual(len(idxs), int(datasize))   // 可以恢复
		}
		return decodeMatrix(grousize, datasize, idxs)
	}
}
func lossDatablocks(datasize uint8, index []uint8) int {
	n := 0
	for _, e := range index {
		if e < datasize {
			n += 1
		}
	}
	if debug.Debug() {
		debug.GreaterOrEqual(int(datasize), n)
	}
	return int(datasize) - n
}

func encodeMatrix(grousize, datasize uint8) (m Matrix) {
	m = baseMatrix(grousize, datasize)

	// 剔除头部的单位矩阵
	idxs := make([]int, m.cols)
	for i := range idxs {
		idxs[i] = i
	}
	m.delRows(idxs...)
	return m
}
func decodeMatrix(grousize, datasize uint8, indexs []uint8) Matrix {
	if debug.Debug() {
		debug.True(slices.IsSorted(indexs))
		debug.Greater(grousize, 0)
		debug.Greater(datasize, 0)
	}
	base := baseMatrix(grousize, datasize)
	defer base.Release()

	// 根据indexs, 获取前datasize个块, 使得m是个方阵
	var del []int
	for i := uint8(0); i < grousize; i++ {
		if !slices.Contains(indexs[:datasize], i) {
			del = append(del, int(i))
		}
	}
	base.delRows(del...)

	// 只是取丢失的数据块对应的行
	m := base.invert()
	var del2 []int
	for i := uint8(0); i < datasize; i++ {
		if slices.Contains(indexs, i) {
			del2 = append(del2, int(i))
		}
	}
	m.delRows(del2...)
	return m
}

func baseMatrix(grousize, datasize uint8) Matrix {
	if debug.Debug() {
		debug.LessOrEqual(datasize, grousize)
	}
	vm := vandermonde(int(grousize), int(datasize))
	defer vm.Release()

	// 获取数据包对应的方阵范德蒙矩阵
	square := vm.sub(0, 0, int(datasize), int(datasize))
	defer square.Release()

	inv := square.invert()
	defer inv.Release()

	m := vm.mul(inv)
	return m
}
