package encodec

import (
	"testing"

	"github.com/lysShub/bytespool-go"
)

func Test_Encodec(t *testing.T) {

	m1 := Encodec(8, 5)
	act1 := Matrix(m1).String()
	exp1 := `[  7,  7,  6,  6,  1]
[  9,  8,  9,  8,  1]
[ 15, 14, 14, 15,  1]`

	if act1 != exp1 {
		t.Fatalf("act1 mismatch:\n%s", act1)
	}

	m2 := Encodec(8, 5, []uint8{0, 2, 4, 5, 6}...)
	act2 := Matrix(m2).String()
	exp2 := `[ 71, 70, 71,  1, 70]
[174,175,175,  1,174]`

	if act2 != exp2 {
		t.Fatalf("act2 mismatch:\n%s", act2)
	}
}

func Test_bytespool(t *testing.T) {
	bytespool.DebugClear()
	const groupsize = 5
	var ms []Matrix
	for datasize := uint8(1); datasize < groupsize; datasize++ {
		ms = append(ms, Encodec(groupsize, datasize))
		if datasize < groupsize {
			idxs := combination(datasize, datasize)
			for _, e := range idxs {
				if lossDatablocks(datasize, e) > 0 {
					ms = append(ms, Encodec(groupsize, datasize, e...))
				}
			}
		}
	}
	for _, e := range ms {
		e.Release()
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

/*
go test -run=none -bench="Benchmark_.*"
goos: windows
goarch: amd64
pkg: acceler/pkg/reedsolomon
cpu: Intel(R) Xeon(R) CPU E5-1650 v4 @ 3.60GHz
Benchmark_baseMatrix-12            31698             37548 ns/op            2372 B/op         20 allocs/op
Benchmark_encodeMatrix-12          32738             34862 ns/op            1603 B/op          8 allocs/op
Benchmark_decodeMatrix-12          32090             37749 ns/op            1755 B/op          9 allocs/op
PASS
ok      acceler/pkg/reedsolomon       10.717s
*/
var (
	groupsize uint8 = 8
	datasize  uint8 = 5
)

func Benchmark_baseMatrix(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = baseMatrix(groupsize, datasize)
	}
}

func Benchmark_encodeMatrix(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = encodeMatrix(groupsize, datasize)
	}
}

func Benchmark_decodeMatrix(b *testing.B) {
	b.ReportAllocs()

	var indexs []uint8
	for i := range groupsize {
		if i != 0 && i != groupsize-1 {
			indexs = append(indexs, i)
		}
	}

	for i := 0; i < b.N; i++ {
		_ = decodeMatrix(groupsize, datasize, indexs)
	}
}
