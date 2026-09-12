package galois

import (
	"fmt"
	"testing"
)

/*

goos: windows
goarch: amd64
pkg: acceler/pkg/reedsolomon/galois
cpu: Intel(R) Xeon(R) CPU E5-1650 v4 @ 3.60GHz
Benchmark_mulVect/SIMD15                16450728                71.61 ns/op      209.46 MB/s
Benchmark_mulVect/Go15                  70695942                17.12 ns/op      876.29 MB/s

Benchmark_mulVect/SIMD16                121489448                9.854 ns/op    1623.69 MB/s
Benchmark_mulVect/Go16                  66609308                17.95 ns/op      891.36 MB/s

Benchmark_mulVect/SIMD32                100000000               11.07 ns/op     2890.14 MB/s
Benchmark_mulVect/Go32                  36836595                32.43 ns/op      986.72 MB/s

Benchmark_mulVect/SIMD512               27383130                40.80 ns/op     12548.89 MB/s
Benchmark_mulVect/Go512                  4430280               261.0 ns/op      1961.58 MB/s

Benchmark_mulVect/SIMD1536              12193659                96.82 ns/op     15864.35 MB/s
Benchmark_mulVect/Go1536                 1612784               743.8 ns/op      2065.04 MB/s

Benchmark_mulXorVect/SIMD15             14644951                81.24 ns/op      184.63 MB/s
Benchmark_mulXorVect/Go15               62398602                18.93 ns/op      792.42 MB/s

Benchmark_mulXorVect/SIMD16             100000000               10.15 ns/op     1575.76 MB/s
Benchmark_mulXorVect/Go16               59567539                19.58 ns/op      817.14 MB/s

Benchmark_mulXorVect/SIMD32             98303445                11.46 ns/op     2792.19 MB/s
Benchmark_mulXorVect/Go32               32138925                36.60 ns/op      874.28 MB/s

Benchmark_mulXorVect/SIMD512            28744304                41.16 ns/op     12440.44 MB/s
Benchmark_mulXorVect/Go512               3672314               326.4 ns/op      1568.65 MB/s

Benchmark_mulXorVect/SIMD1536           11540258               104.9 ns/op      14644.63 MB/s
Benchmark_mulXorVect/Go1536              1259648               959.7 ns/op      1600.46 MB/s

Benchmark_xorVect/SIMD15                100000000               10.66 ns/op     1406.74 MB/s
Benchmark_xorVect/Go15                  79680745                14.43 ns/op     1039.66 MB/s

Benchmark_xorVect/SIMD16                215010067                5.564 ns/op    2875.72 MB/s
Benchmark_xorVect/Go16                  73191260                15.75 ns/op     1016.00 MB/s

Benchmark_xorVect/SIMD32                197041099                6.138 ns/op    5213.77 MB/s
Benchmark_xorVect/Go32                  184884006                6.446 ns/op    4964.11 MB/s

Benchmark_xorVect/SIMD512               44116024                27.03 ns/op     18941.74 MB/s
Benchmark_xorVect/Go512                 31478794                37.45 ns/op     13672.77 MB/s

Benchmark_xorVect/SIMD1536              21240216                56.23 ns/op     27314.06 MB/s
Benchmark_xorVect/Go1536                10956122               108.6 ns/op      14138.91 MB/s

Benchmark_lastIndex/SIMD15              317173267                3.772 ns/op    3976.90 MB/s
Benchmark_lastIndex/Go15                129539467                9.259 ns/op    1620.12 MB/s

Benchmark_lastIndex/SIMD16              383559243                3.125 ns/op    5119.50 MB/s
Benchmark_lastIndex/Go16                122508133                9.766 ns/op    1638.34 MB/s

Benchmark_lastIndex/SIMD32              316973368                3.762 ns/op    8506.20 MB/s
Benchmark_lastIndex/Go32                65914870                18.11 ns/op     1766.75 MB/s

Benchmark_lastIndex/SIMD512             90974564                12.38 ns/op     41373.61 MB/s
Benchmark_lastIndex/Go512                7771345               151.9 ns/op      3370.75 MB/s

Benchmark_lastIndex/SIMD1536            34716093                33.75 ns/op     45505.91 MB/s
Benchmark_lastIndex/Go1536               2862682               419.7 ns/op      3659.51 MB/s
PASS
ok      acceler/pkg/reedsolomon/galois        62.754s

*/

func Benchmark_mulVect(b *testing.B) {
	var c byte = 55
	for _, n := range blens {
		println("")

		b.Run(fmt.Sprintf("SIMD%d", n), func(b *testing.B) {
			var in = make([]byte, n)
			var out = make([]byte, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(in)))
				mulVect(c, in, out)
			}
		})
		b.Run(fmt.Sprintf("Go%d", n), func(b *testing.B) {
			var in = make([]byte, n)
			var out = make([]byte, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(in)))
				mulVectGo(c, in, out)
			}
		})
	}
}
func Benchmark_mulXorVect(b *testing.B) {
	var c byte = 55
	for _, n := range blens {
		println("")

		b.Run(fmt.Sprintf("SIMD%d", n), func(b *testing.B) {
			var in = make([]byte, n)
			var out = make([]byte, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(in)))
				mulXorVect(c, in, out)
			}
		})
		b.Run(fmt.Sprintf("Go%d", n), func(b *testing.B) {
			var in = make([]byte, n)
			var out = make([]byte, n)

			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(in)))
				mulXorVectGo(c, in, out)
			}
		})
	}
}
func Benchmark_xorVect(b *testing.B) {
	for _, n := range blens {
		println("")

		b.Run(fmt.Sprintf("SIMD%d", n), func(b *testing.B) {
			var in = make([]byte, n)
			var out = make([]byte, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(in)))
				xorVect(in, out)
			}
		})
		b.Run(fmt.Sprintf("Go%d", n), func(b *testing.B) {
			var in = make([]byte, n)
			var out = make([]byte, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(in)))
				xorVectGo(in, out)
			}
		})
	}
}

func Benchmark_lastIndex(b *testing.B) {
	for _, n := range blens {
		println("")

		b.Run(fmt.Sprintf("SIMD%d", n), func(b *testing.B) {
			var in = make([]byte, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(in)))
				lastIndex(in, 0xff)
			}
		})
		b.Run(fmt.Sprintf("Go%d", n), func(b *testing.B) {
			var in = make([]byte, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(in)))
				lastIndexGo(in, 0xff)
			}
		})
	}
}
