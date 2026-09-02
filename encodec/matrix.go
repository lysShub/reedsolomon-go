package encodec

import (
	"slices"

	"github.com/lysShub/debug-go"
	"github.com/lysShub/reedsolomon-go/galois"

	"github.com/lysShub/bytespool-go"
)

var Pooler bytespool.Pooler[[]byte, byte] = bytespool.Pool[[]byte, byte]{}

// Matrix 行优先存储, 整个矩阵为一块bytespool内存
type Matrix struct {
	b    []byte
	rows int
	cols int
}

func Make(rows, cols int) Matrix {
	if debug.Debug() {
		debug.Greater(rows, 0)
		debug.Greater(cols, 0)
	}
	return Matrix{
		b:    Pooler.Get(rows * cols),
		rows: rows,
		cols: cols,
	}
}
func (m Matrix) raw() []byte { return m.b }
func (m Matrix) Rows() int   { return m.rows }
func (m Matrix) Cols() int   { return m.cols }
func (m Matrix) Len() int    { return len(m.b) }
func (m Matrix) Row(i int) []byte {
	if debug.Debug() {
		debug.GreaterOrEqual(i, 0)
		debug.Less(i, m.rows)
	}
	return m.b[i*m.cols : (i+1)*m.cols]
}
func (m *Matrix) Release() {
	if m.b != nil {
		Pooler.Put(m.b)
		m.b, m.rows, m.cols = nil, 0, 0
	}
}

func (m Matrix) Clone() Matrix {
	m1 := Make(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		copy(m1.Row(i), m.Row(i))
	}
	return m1
}

// vandermonde 创建一个范德蒙矩阵
func vandermonde(rows, cols int) Matrix {
	m := Make(rows, cols)
	for rowi := 0; rowi < rows; rowi++ {
		row := m.Row(rowi)
		for coli := range row {
			row[coli] = galois.Pow(byte(rowi), byte(coli))
		}
	}
	return m
}

func (m *Matrix) delRows(idxs ...int) {
	if len(idxs) == 0 {
		return
	}
	if debug.Debug() {
		debug.Equal(len(m.b), m.rows*m.cols)
		debug.LessOrEqual(len(idxs), m.rows)
	}

	n, w := m.cols, 0
	for r := 0; r < m.rows; r++ {
		if slices.Contains(idxs, r) {
			continue
		}
		if w != r {
			copy(m.Row(w), m.Row(r))
		}
		w++
	}
	m.b = m.b[:w*n]
	m.rows = w
}

func (m *Matrix) swapRow(r1, r2 int) {
	if debug.Debug() {
		debug.GreaterOrEqual(r1, 0)
		debug.Less(r1, m.rows)
		debug.GreaterOrEqual(r2, 0)
		debug.Less(r2, m.rows)
	}
	if r1 == r2 {
		return
	}
	a, b := m.Row(r1), m.Row(r2)
	for i := range a {
		a[i], b[i] = b[i], a[i]
	}
}

func (m Matrix) sub(rmin, cmin, rmax, cmax int) Matrix {
	if debug.Debug() {
		debug.Equal(len(m.b), m.rows*m.cols)
	}
	res := Make(rmax-rmin, cmax-cmin)
	for ri := rmin; ri < rmax; ri++ {
		copy(res.Row(ri-rmin), m.Row(ri)[cmin:cmax])
	}
	return res
}

func (m Matrix) mul(right Matrix) Matrix {
	if debug.Debug() {
		debug.Equal(len(m.b), m.rows*m.cols)
		debug.Equal(len(right.b), right.rows*right.cols)
		debug.Equal(m.cols, right.rows)
	}
	res := Make(m.rows, right.cols)
	for r := 0; r < m.rows; r++ {
		for c := 0; c < right.cols; c++ {
			res.b[r*res.cols+c] = right.mulColSum(m.Row(r), c)
		}
	}
	return res
}

// mulColSum add(s * m.col[i])
func (m Matrix) mulColSum(s []byte, i int) (sum byte) {
	if debug.Debug() {
		debug.Equal(len(s), m.rows)
		debug.Less(i, m.cols)
	}
	for j, e := range s {
		v := galois.Mul(e, m.b[j*m.cols+i])
		sum = galois.Add(sum, v)
	}
	return sum
}

// invert 求取逆矩阵
func (m Matrix) invert() Matrix {
	if debug.Debug() {
		debug.Equal(len(m.b), m.rows*m.cols)
		debug.Equal(m.rows, m.cols, "require square matrix")
	}

	n := m.rows
	// work: [m E]  在m右侧拼接一个单位矩阵
	work := Make(n, n*2)
	for r := 0; r < n; r++ {
		row := work.Row(r)
		copy(row, m.Row(r))
		clear(row[n:])
		row[n+r] = 1
	}
	work.gaussianElimination()

	// 原地裁剪: 将右侧结果搬至左侧
	for r := 0; r < n; r++ {
		copy(work.b[r*n:(r+1)*n], work.Row(r)[n:])
	}
	work.b = work.b[:n*n]
	work.cols = n
	return work
}

// gaussianElimination 高斯消元(原地)
func (m *Matrix) gaussianElimination() {
	if debug.Debug() {
		debug.Equal(len(m.b), m.rows*m.cols)
	}
	cols := m.cols
	for r := 0; r < m.rows; r++ {
		row := m.Row(r)
		if row[r] == 0 {
			for r2 := r + 1; r2 < m.rows; r2++ {
				if m.Row(r2)[r] != 0 {
					m.swapRow(r, r2)
					row = m.Row(r)
					break
				}
			}
		}
		if row[r] == 0 {
			panic("") // 不可逆矩阵(奇异矩阵)
		}

		if row[r] != 1 {
			q := galois.Div(1, row[r])
			for c := 0; c < cols; c++ {
				row[c] = galois.Mul(row[c], q)
			}
		}
		for r2 := r + 1; r2 < m.rows; r2++ {
			row2 := m.Row(r2)
			q := row2[r]
			if q != 0 {
				for c := 0; c < cols; c++ {
					row2[c] ^= galois.Mul(q, row[c])
				}
			}
		}
	}

	for r := 0; r < m.rows; r++ {
		row := m.Row(r)
		for r2 := 0; r2 < r; r2++ {
			row2 := m.Row(r2)
			q := row2[r]
			if q != 0 {
				for c := 0; c < cols; c++ {
					row2[c] ^= galois.Mul(q, row[c])
				}
			}
		}
	}
}
