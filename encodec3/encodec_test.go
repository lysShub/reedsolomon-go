package encodec3

import (
	"bytes"
	"testing"

	"github.com/lysShub/bytespool-go"
	"github.com/lysShub/reedsolomon-go/encodec"
)

func TestXxxx(t *testing.T) {
	const g = 24
	const maxb = g*(g-1) + 4*(g-1)*(g-1)
}

func Test_Encodec(t *testing.T) {
	cases := []struct {
		g, d uint8
		idxs []uint8
	}{
		{3, 2, nil},
		{5, 3, nil},
		{8, 5, nil},
		{10, 4, nil},
		{20, 8, nil},
		{64, 32, nil},
		{128, 64, nil},
		{255, 254, nil},
		{3, 1, nil},
		{8, 1, nil},

		{3, 2, []uint8{0, 2}},
		{5, 3, []uint8{0, 2, 4}},
		{5, 3, []uint8{1, 2, 3, 4}},
		{8, 5, []uint8{0, 2, 4, 5, 6}},
		{8, 5, []uint8{1, 2, 3, 4, 6, 7}},
		{10, 4, []uint8{0, 1, 3, 4}},
		{20, 8, []uint8{0, 1, 3, 4, 6, 7, 8, 9}},
		{64, 32, []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63}},
		{3, 1, []uint8{1}},
	}

	for _, tc := range cases {
		exp := encodec.Encodec(tc.g, tc.d, tc.idxs...)
		act := make([]byte, int(tc.g)*int(tc.d))
		n := Encodec(act, tc.g, tc.d, tc.idxs...)

		cols := int(tc.d)
		rows := exp.Rows()
		if n != rows*cols || exp.Cols() != cols {
			t.Fatalf("%v: shape mismatch: encodec=%dx%d encode3=%d bytes", tc, exp.Rows(), exp.Cols(), n)
		}
		for i := 0; i < rows; i++ {
			if !bytes.Equal(act[i*cols:(i+1)*cols], exp.Row(i)) {
				t.Fatalf("%v: matrix mismatch at row %d\nencodec:\n%v\nencode3:\n%v", tc, i, exp.Row(i), act[i*cols:(i+1)*cols])
			}
		}
	}
}

func Test_Encodec_String(t *testing.T) {
	dst := make([]byte, 8*5)

	n := Encodec(dst, 8, 5)
	exp1 := `[  7,  7,  6,  6,  1]
[  9,  8,  9,  8,  1]
[ 15, 14, 14, 15,  1]`
	if got := matrixString(dst[:n], n/5); got != exp1 {
		t.Fatalf("act1 mismatch:\n%s", got)
	}

	n = Encodec(dst, 8, 5, []uint8{0, 2, 4, 5, 6}...)
	exp2 := `[ 71, 70, 71,  1, 70]
[174,175,175,  1,174]`
	if got := matrixString(dst[:n], n/5); got != exp2 {
		t.Fatalf("act2 mismatch:\n%s", got)
	}
}

func Test_bytespool(t *testing.T) {
	bytespool.DebugClear()
	const groupsize = 5
	dst := make([]byte, 255*255)
	for datasize := uint8(1); datasize < groupsize; datasize++ {
		Encodec(dst, groupsize, datasize)
		idxs := combination(datasize, datasize)
		for _, e := range idxs {
			if lossDatablocks(datasize, e) > 0 {
				Encodec(dst, groupsize, datasize, e...)
			}
		}
	}
	if n := bytespool.DebugLength(); n != 0 {
		t.Fatalf("DebugLength = %d, want 0", n)
	}
}

func combination[T uint8 | int](n, m T) [][]T {
	var result [][]T
	if m < 0 || m > n || n < 0 {
		return result
	}
	current := make([]T, 0, m)
	var backtrack func(start T)
	backtrack = func(start T) {
		if len(current) == int(m) {
			temp := make([]T, m)
			copy(temp, current)
			result = append(result, temp)
			return
		}
		if int(n)-int(start) < int(m)-len(current) {
			return
		}
		for i := start; i < n; i++ {
			current = append(current, i)
			backtrack(i + 1)
			current = current[:len(current)-1]
		}
	}
	backtrack(0)
	return result
}

func lossDatablocks(datasize uint8, index []uint8) int {
	n := 0
	for _, e := range index {
		if e < datasize {
			n += 1
		}
	}
	return int(datasize) - n
}

var (
	groupsize uint8 = 8
	datasize  uint8 = 5
)

func Benchmark_baseMatrix(b *testing.B) {
	dst := make([]byte, 255*255)
	buf := make([]byte, 255*255+4*255*255)
	b.ReportAllocs()

	for b.Loop() {
		baseMatrixBuf(dst, buf, int(groupsize), int(datasize))
	}
}

func Benchmark_encodeMatrix(b *testing.B) {
	dst := make([]byte, 255*255)
	b.ReportAllocs()

	for b.Loop() {
		encodeMatrix(dst, int(groupsize), int(datasize))
	}
}

func Benchmark_decodeMatrix(b *testing.B) {
	dst := make([]byte, 255*255)
	var indexs []uint8
	for i := range groupsize {
		if i != 0 && i != groupsize-1 {
			indexs = append(indexs, i)
		}
	}
	b.ReportAllocs()

	for b.Loop() {
		decodeMatrix(dst, int(groupsize), int(datasize), indexs)
	}
}
