package encodec3_test

import (
	"bytes"
	"testing"

	"github.com/lysShub/reedsolomon-go/encodec"
	"github.com/lysShub/reedsolomon-go/encodec3"
)

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
		{64, 32, []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63}},
		{3, 1, []uint8{1}},
	}

	for _, tc := range cases {
		exp := encodec.Encodec(tc.g, tc.d, tc.idxs...)
		act := make([]byte, int(tc.g)*int(tc.d))
		n := encodec3.Encodec(act, tc.g, tc.d, tc.idxs...)

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
