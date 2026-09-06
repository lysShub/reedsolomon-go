package encodec

import (
	"slices"

	"github.com/lysShub/reedsolomon-go/galois"
)

// vandermonde create a vandermonde-matrix
func vandermonde(dst []byte, rows int) {
	cols := len(dst) / rows
	for rowi := 0; rowi < rows; rowi++ {
		row := dst[rowi*cols : (rowi+1)*cols]
		for coli := range row {
			row[coli] = galois.Pow(byte(rowi), byte(coli))
		}
	}
}

// sub extract sub-matrix from m
func sub(dst, m []byte, rows int, rmin, cmin, rmax, cmax int) {
	cols := len(m) / rows
	w := cmax - cmin
	for ri := rmin; ri < rmax; ri++ {
		copy(dst[(ri-rmin)*w:(ri-rmin+1)*w], m[ri*cols+cmin:ri*cols+cmax])
	}
}

func mul(dst, l []byte, lrows int, r []byte, rrows int) {
	lcols := len(l) / lrows
	rcols := len(r) / rrows
	for i := 0; i < lrows; i++ {
		row := dst[i*rcols : (i+1)*rcols]
		lrow := l[i*lcols : (i+1)*lcols]
		for j := 0; j < rcols; j++ {
			row[j] = mulColSum(r, rrows, lrow, j)
		}
	}
}

func mulColSum(m []byte, rows int, s []byte, c int) (sum byte) {
	mCols := len(m) / rows
	for j, e := range s {
		v := galois.Mul(e, m[j*mCols+c])
		sum = galois.Add(sum, v)
	}
	return sum
}

// invert calculate m's invert-matrix
func invert(dst, work, m []byte, n int) {
	for r := 0; r < n; r++ {
		copy(work[r*2*n:(r+1)*2*n], m[r*n:(r+1)*n])
		clear(work[r*2*n+n : (r+1)*2*n])
		work[r*2*n+n+r] = 1
	}
	gaussianElimination(work, n)
	for r := 0; r < n; r++ {
		copy(dst[r*n:(r+1)*n], work[r*2*n+n:(r+1)*2*n])
	}
}

func swapRow(m []byte, rows int, r1, r2 int) {
	if r1 == r2 {
		return
	}
	cols := len(m) / rows
	a := m[r1*cols : (r1+1)*cols]
	b := m[r2*cols : (r2+1)*cols]
	for i := range a {
		a[i], b[i] = b[i], a[i]
	}
}

func gaussianElimination(m []byte, rows int) {
	cols := len(m) / rows
	for r := 0; r < rows; r++ {
		row := m[r*cols : (r+1)*cols]
		if row[r] == 0 {
			for r2 := r + 1; r2 < rows; r2++ {
				if m[r2*cols+r] != 0 {
					swapRow(m, rows, r, r2)
					row = m[r*cols : (r+1)*cols]
					break
				}
			}
		}
		if row[r] == 0 {
			panic("")
		}

		if row[r] != 1 {
			q := galois.Div(1, row[r])
			for c := 0; c < cols; c++ {
				row[c] = galois.Mul(row[c], q)
			}
		}
		for r2 := r + 1; r2 < rows; r2++ {
			row2 := m[r2*cols : (r2+1)*cols]
			q := row2[r]
			if q != 0 {
				for c := 0; c < cols; c++ {
					row2[c] ^= galois.Mul(q, row[c])
				}
			}
		}
	}

	for r := 0; r < rows; r++ {
		row := m[r*cols : (r+1)*cols]
		for r2 := 0; r2 < r; r2++ {
			row2 := m[r2*cols : (r2+1)*cols]
			q := row2[r]
			if q != 0 {
				for c := 0; c < cols; c++ {
					row2[c] ^= galois.Mul(q, row[c])
				}
			}
		}
	}
}

func delRows(m []byte, rows int, idxs ...int) int {
	if len(idxs) == 0 {
		return rows
	}
	cols := len(m) / rows
	w := 0
	for r := 0; r < rows; r++ {
		if slices.Contains(idxs, r) {
			continue
		}
		if w != r {
			copy(m[w*cols:(w+1)*cols], m[r*cols:(r+1)*cols])
		}
		w++
	}
	return w
}
