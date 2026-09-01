package reedsolomon_test

import (
	"acceler/pkg/reedsolomon"
	"bytes"
	"crypto/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
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
	t.Run("{2, 3} 丢失一个数据包", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 3, Datasize: 2}

		// 发送数据,  前两个为数据包, 最后一个为校验包
		var datas = make([][]byte, 2)
		datas[0] = []byte{1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		datas[1] = []byte{4, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		var parity = makes(1, len(datas[0]))

		{ // 编码
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 1) // 丢失 datas[1]

		{ // 恢复
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{0, 2}, reconst)
			require.Equal(t, 1, n)
			require.Equal(t, datas[1], reconst[0])
		}
	})

	t.Run("{3, 5} 丢失一个数据包", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 5, Datasize: 3}

		var datas = make([][]byte, 3)
		datas[0] = []byte{1, 2, 3, 4}
		datas[1] = []byte{4, 3, 0, 0}
		datas[2] = []byte{2, 1, 6, 9}
		var parity = makes(2, len(datas[0]))

		{ // 编码
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
			reedsolomon.Encode(para, datas[2], 2, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 1) // 丢失datas[1]

		{ // 恢复
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{1, 2, 3, 4}, reconst)
			require.Equal(t, 1, n)
			require.Equal(t, datas[1], reconst[0])
		}
	})

	t.Run("{3, 5} 丢失一个数据包和一个校验包", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 5, Datasize: 3}

		var datas = make([][]byte, 3)
		datas[0] = []byte{1, 2, 3, 4}
		datas[1] = []byte{4, 3, 0, 0}
		datas[2] = []byte{2, 1, 6, 9}
		var parity = makes(2, len(datas[0]))

		{ // 编码
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
			reedsolomon.Encode(para, datas[2], 2, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 1, 3)

		{ // 恢复
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{0, 2, 4}, reconst)
			require.Equal(t, 1, n)
			require.Equal(t, datas[1], reconst[0])
		}
	})

	t.Run("{3, 5} 丢失两个paritysize", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 5, Datasize: 3}

		var datas = make([][]byte, 3)
		datas[0] = []byte{1, 2, 3, 4}
		datas[1] = []byte{4, 3, 0, 0}
		datas[2] = []byte{2, 1, 6, 9}
		var parity = makes(2, len(datas[0]))

		{ // 编码
			reedsolomon.Encode(para, datas[0], 0, parity)
			reedsolomon.Encode(para, datas[1], 1, parity)
			reedsolomon.Encode(para, datas[2], 2, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 3, 4)

		{ // 恢复
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{0, 1, 2}, reconst)
			require.Zero(t, n)
		}
	})

	t.Run("datasize ==  1", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 3, Datasize: 1}

		// 发送数据,  前一个为数据包, 最后两个为校验包
		var datas = make([][]byte, 1)
		datas[0] = []byte{1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		var parity = makes(2, len(datas[0]))

		{ // 编码
			reedsolomon.Encode(para, datas[0], 0, parity)
		}

		var blocks = append(datas, parity...)
		blocks = deletes(blocks, 0, 2) // 丢失 datas[0,2]

		{ // 恢复
			var reconst = makes(1, len(datas[0]))
			n := reedsolomon.Reconst(para, blocks, []uint8{1}, reconst)
			require.Zero(t, 0, n)
			require.Equal(t, datas[0], reconst[0])
		}
	})
}

func Test_NotAlign(t *testing.T) {
	// 测试非对其数据块进行编码与回复

	var buildBlocks = func(datas [][]byte, parity [][]byte) [][]byte {
		var blocks [][]byte

		size := 0
		for _, e := range datas {
			blocks = append(blocks, slices.Clone(e))
			size = max(size, len(e))
		}
		for _, e := range parity {
			// 校验块长度等于最长数据块长度
			blocks = append(blocks, slices.Clone(e[:size]))
		}
		return blocks
	}

	t.Run("{2, 3} block长度不一致1", func(t *testing.T) {
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
			require.Equal(t, 1, n)
			require.True(t, bytes.HasPrefix(reconst[0], datas[1]))
		}
	})

	t.Run("{2, 3} block长度不一致2", func(t *testing.T) {
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
			require.Equal(t, 1, n)
			require.True(t, bytes.HasPrefix(reconst[0], datas[0]))
		}
	})

	t.Run("{2, 3} 校验矩阵不为零值, 且datas[0]不是最长数据块", func(t *testing.T) {
		var para = reedsolomon.Para{Groupsize: 3, Datasize: 2}

		var datas = make([][]byte, 2)
		datas[0] = []byte{4, 3}
		datas[1] = []byte{1, 2, 3, 4}

		var parity = makes(0xff, 0xff)
		for _, e := range parity {
			rand.Read(e)
		}
		{
			// 要求 parity[i][len(datas[0]):] 为零值
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
			// 要求 reconst[i][len(blocks[minidx]):] 为零值
			for _, e := range reconst {
				clear(e[len(blocks[0]):])
			}

			n := reedsolomon.Reconst(para, blocks, []uint8{0, 2}, reconst)
			require.Equal(t, 1, n)
			require.True(t, bytes.HasPrefix(reconst[0], datas[1]))
		}
	})

}
