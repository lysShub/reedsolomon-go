package encodec

import (
	"slices"

	"github.com/lysShub/bytespool-go"
)

var Pooler bytespool.Pooler[[]byte, byte] = bytespool.Pool[[]byte, byte]{}

func Encodec(matrix []byte, grousize, datasize uint8, idxs ...uint8) int {
	g, d := int(grousize), int(datasize)
	if len(idxs) == 0 {
		return encodeMatrix(matrix, g, d)
	}
	return decodeMatrix(matrix, g, d, idxs)
}

const stackAlloc = 1024 + 512

func encodeMatrix(dst []byte, g, d int) int {
	total := g*d + 4*d*d

	if total <= stackAlloc {
		var b [stackAlloc]byte
		return encodeMatrixBuf(dst, b[:], g, d)
	} else {
		b := Pooler.Get(total)
		defer Pooler.Put(b)
		return encodeMatrixBuf(dst, b, g, d)
	}
}
func encodeMatrixBuf(dst, buf []byte, g, d int) int {
	baseMatrixBuf(dst, buf, g, d)

	var idxs [256]int
	for i := range idxs[:d] {
		idxs[i] = i
	}
	rows := delRows(dst, g, idxs[:d]...)
	return rows * d
}

func decodeMatrix(dst []byte, g, d int, indexs []uint8) int {
	total := g*d + 4*d*d

	if total <= stackAlloc {
		var b [stackAlloc]byte
		return decodeMatrixBuf(dst, b[:], g, d, indexs)
	} else {
		b := Pooler.Get(total)
		defer Pooler.Put(b)
		return decodeMatrixBuf(dst, b, g, d, indexs)
	}
}
func decodeMatrixBuf(dst, buf []byte, g, d int, indexs []uint8) int {
	baseMatrixBuf(dst, buf, g, d)

	var del [256]int
	var delN int
	for i := uint8(0); i < uint8(g); i++ {
		if !slices.Contains(indexs[:d], i) {
			del[delN] = int(i)
			delN++
		}
	}
	n := delRows(dst, g, del[:delN]...)

	inv := buf[g*d : g*d+d*d]
	work := buf[g*d+2*d*d : g*d+4*d*d]
	invert(inv[:n*n], work, dst[:n*n], n)

	var del2 [256]int
	var del2N int
	for i := uint8(0); i < uint8(d); i++ {
		if slices.Contains(indexs, i) {
			del2[del2N] = int(i)
			del2N++
		}
	}
	rows := delRows(inv[:n*n], n, del2[:del2N]...)
	copy(dst, inv[:rows*n])
	return rows * n
}

func baseMatrixBuf(dst, buf []byte, g, d int) {
	vm := buf[:g*d]
	sq := buf[g*d : g*d+d*d]
	inv := buf[g*d+d*d : g*d+2*d*d]
	work := buf[g*d+2*d*d : g*d+4*d*d]

	vandermonde(vm, g)
	sub(sq, vm, g, 0, 0, d, d)
	invert(inv, work, sq, d)
	mul(dst, vm, g, inv, d)
}
