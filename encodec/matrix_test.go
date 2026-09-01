package encodec

import (
	"crypto/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lysShub/bytespool-go"

	"github.com/stretchr/testify/require"
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

		require.Equal(t, "[152,  2,  3]", m.String())
	})

	t.Run("String", func(t *testing.T) {
		var m = fromRows([][]byte{
			{0, 115, 255},
			{152, 2, 3},
		})
		defer m.Release()
		require.Equal(t, "[  0,115,255]\n[152,  2,  3]", m.String())
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
		require.Equal(t, "[209,157]\n[134,191]", r.String())
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
		require.Equal(t, "[172, 26, 17]\n[104, 44,209]\n[126,130,215]", inv.String())
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
		require.Equal(t, "[112]", m1.String())

		m2 := m.sub(1, 1, 3, 3)
		defer m2.Release()
		require.Equal(t, "[112,  3]\n[  1,177]", m2.String())
	})

	t.Run("swapRow", func(t *testing.T) {
		bytespool.DebugClear()
		defer func() {
			require.Zero(t, bytespool.DebugLength())
		}()

		m := Make(2, 2)
		defer m.Release()
		rand.Read(m.Row(0))
		rand.Read(m.Row(1))
		bak := m.Clone()
		defer bak.Release()

		m.swapRow(0, 1)

		require.Equal(t, bak.Row(1), m.Row(0))
		require.Equal(t, bak.Row(0), m.Row(1))
	})

	t.Run("delRows 2", func(t *testing.T) {
		bytespool.DebugClear()
		defer func() {
			require.Zero(t, bytespool.DebugLength())
		}()

		m := Make(2, 2)
		defer m.Release()
		rand.Read(m.Row(0))
		rand.Read(m.Row(1))
		bak := m.Clone()
		defer bak.Release()

		m.delRows(0)

		require.Equal(t, 1, m.Rows())
		require.Equal(t, bak.Row(1), m.Row(0))
	})

	t.Run("delRows 3", func(t *testing.T) {
		bytespool.DebugClear()
		defer func() {
			require.Zero(t, bytespool.DebugLength())
		}()

		m := Make(2, 2)
		defer m.Release()
		rand.Read(m.Row(0))
		rand.Read(m.Row(1))

		m.delRows(0)
		m.delRows(0)

		require.Equal(t, 0, m.Rows())
	})

}
