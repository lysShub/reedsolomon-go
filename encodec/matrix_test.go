package encodec

import (
	"bytes"
	"crypto/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lysShub/bytespool-go"
)

func fromRows(rows [][]byte) []byte {
	cols := len(rows[0])
	m := make([]byte, len(rows)*cols)
	for i, r := range rows {
		copy(m[i*cols:(i+1)*cols], r)
	}
	return m
}

func matrixString(m []byte, rows int) string {
	cols := len(m) / rows
	lines := make([]string, 0, rows)
	var num = make([]byte, 0, 3)
	for r := 0; r < rows; r++ {
		row := m[r*cols : (r+1)*cols]
		var b = make([]byte, 0, len(row)*4+2)
		b = append(b, '[')
		for i, v := range row {
			num = strconv.AppendInt(num[:0], int64(v), 10)
			for i := len(num); i < 3; i++ {
				b = append(b, ' ')
			}
			b = append(b, num...)
			if i < len(row)-1 {
				b = append(b, ',')
			}
		}
		b = append(b, ']')
		lines = append(lines, string(b))
	}
	return strings.Join(lines, "\n")
}

func Test_Matrix(t *testing.T) {
	t.Run("delRows", func(t *testing.T) {
		m := fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
		})
		rows := delRows(m, 2, 0)
		if got := matrixString(m[:rows*3], rows); got != "[152,  2,  3]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("String", func(t *testing.T) {
		m := fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
		})
		if got := matrixString(m, 2); got != "[  0,115,255]\n[152,  2,  3]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("mul", func(t *testing.T) {
		m1 := fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
		})
		m2 := fromRows([][]byte{
			{34, 67},
			{77, 12},
			{111, 1},
		})
		dst := make([]byte, 2*2)
		mul(dst, m1, 2, m2, 3)
		if got := matrixString(dst, 2); got != "[209,157]\n[134,191]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("invert", func(t *testing.T) {
		m1 := fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
			{111, 1, 77},
		})
		dst := make([]byte, 3*3)
		work := make([]byte, 3*3*2)
		invert(dst, work, m1, 3)
		if got := matrixString(dst, 3); got != "[172, 26, 17]\n[104, 44,209]\n[126,130,215]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("sub", func(t *testing.T) {
		m := fromRows([][]byte{
			{0, 115, 255},
			{152, 112, 3},
			{111, 1, 177},
		})
		m1 := make([]byte, 1*1)
		sub(m1, m, 3, 1, 1, 2, 2)
		if got := matrixString(m1, 1); got != "[112]" {
			t.Fatalf("mismatch:\n%s", got)
		}

		m2 := make([]byte, 2*2)
		sub(m2, m, 3, 1, 1, 3, 3)
		if got := matrixString(m2, 2); got != "[112,  3]\n[  1,177]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("swapRow", func(t *testing.T) {
		bytespool.DebugClear()
		defer func() {
			if n := bytespool.DebugLength(); n != 0 {
				t.Fatalf("DebugLength = %d, want 0", n)
			}
		}()

		m := make([]byte, 2*2)
		rand.Read(m)
		bak := make([]byte, 2*2)
		copy(bak, m)

		swapRow(m, 2, 0, 1)

		if !bytes.Equal(m[2:4], bak[0:2]) {
			t.Fatalf("Row(0) mismatch")
		}
		if !bytes.Equal(m[0:2], bak[2:4]) {
			t.Fatalf("Row(1) mismatch")
		}
	})

	t.Run("delRows 2", func(t *testing.T) {
		bytespool.DebugClear()
		defer func() {
			if n := bytespool.DebugLength(); n != 0 {
				t.Fatalf("DebugLength = %d, want 0", n)
			}
		}()

		m := make([]byte, 2*2)
		rand.Read(m)
		bak := make([]byte, 2*2)
		copy(bak, m)

		rows := delRows(m, 2, 0)

		if rows != 1 {
			t.Fatalf("Rows = %d, want 1", rows)
		}
		if !bytes.Equal(m[:2], bak[2:4]) {
			t.Fatalf("Row(0) mismatch")
		}
	})

	t.Run("delRows 3", func(t *testing.T) {
		bytespool.DebugClear()
		defer func() {
			if n := bytespool.DebugLength(); n != 0 {
				t.Fatalf("DebugLength = %d, want 0", n)
			}
		}()

		m := make([]byte, 2*2)
		rand.Read(m)

		rows := delRows(m, 2, 0)
		rows = delRows(m[:rows*2], rows, 0)

		if rows != 0 {
			t.Fatalf("Rows = %d, want 0", rows)
		}
	})
}
