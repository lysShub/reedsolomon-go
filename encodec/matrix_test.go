package encodec

import (
	"bytes"
	"crypto/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lysShub/bytespool-go"
)

func fromRows(rows [][]byte) Matrix {
	m := Make(len(rows), len(rows[0]))
	for i := range rows {
		copy(m.Row(i), rows[i])
	}
	return m
}

func (m Matrix) String() string {
	rows := make([]string, 0, m.rows)

	var num = make([]byte, 0, 3)
	for r := 0; r < m.rows; r++ {
		row := m.Row(r)
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

		rows = append(rows, string(b))
	}
	return strings.Join(rows, "\n")
}

func Test_Matrix(t *testing.T) {

	t.Run("delRows", func(t *testing.T) {
		var m = fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
		})
		defer m.Release()
		m.delRows(0)

		if got := m.String(); got != "[152,  2,  3]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("String", func(t *testing.T) {
		var m = fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
		})
		defer m.Release()
		if got := m.String(); got != "[  0,115,255]\n[152,  2,  3]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("mul", func(t *testing.T) {
		var m1 = fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
		})
		var m2 = fromRows([][]byte{
			{34, 67},
			{77, 12},
			{111, 1},
		})
		defer m1.Release()
		defer m2.Release()

		r := m1.mul(m2)
		defer r.Release()
		if got := r.String(); got != "[209,157]\n[134,191]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("invert", func(t *testing.T) {
		var m1 = fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
			{111, 1, 77},
		})
		defer m1.Release()

		inv := m1.invert()
		defer inv.Release()
		if got := inv.String(); got != "[172, 26, 17]\n[104, 44,209]\n[126,130,215]" {
			t.Fatalf("mismatch:\n%s", got)
		}
	})

	t.Run("sub", func(t *testing.T) {
		var m = fromRows([][]byte{
			{0, 115, 255},
			{152, 112, 3},
			{111, 1, 177},
		})
		defer m.Release()

		m1 := m.sub(1, 1, 2, 2)
		defer m1.Release()
		if got := m1.String(); got != "[112]" {
			t.Fatalf("mismatch:\n%s", got)
		}

		m2 := m.sub(1, 1, 3, 3)
		defer m2.Release()
		if got := m2.String(); got != "[112,  3]\n[  1,177]" {
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

		m := Make(2, 2)
		defer m.Release()
		rand.Read(m.Row(0))
		rand.Read(m.Row(1))
		bak := m.Clone()
		defer bak.Release()

		m.swapRow(0, 1)

		if !bytes.Equal(bak.Row(1), m.Row(0)) {
			t.Fatalf("Row(0) mismatch")
		}
		if !bytes.Equal(bak.Row(0), m.Row(1)) {
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

		m := Make(2, 2)
		defer m.Release()
		rand.Read(m.Row(0))
		rand.Read(m.Row(1))
		bak := m.Clone()
		defer bak.Release()

		m.delRows(0)

		if n := m.Rows(); n != 1 {
			t.Fatalf("Rows = %d, want 1", n)
		}
		if !bytes.Equal(bak.Row(1), m.Row(0)) {
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

		m := Make(2, 2)
		defer m.Release()
		rand.Read(m.Row(0))
		rand.Read(m.Row(1))

		m.delRows(0)
		m.delRows(0)

		if n := m.Rows(); n != 0 {
			t.Fatalf("Rows = %d, want 0", n)
		}
	})

}
