# reedsolomon-go

_设置 `go env -w GOEXPERIMENT=simd`_

这是专为网络设计 reedsolomon fec, 与其他实现的主要区别是：

1. 无状态的，支持流式编解码
2. 不要求 block 长度一致
3. 系统码 -- 不会额外引入延时
4. 高性能 -- `SIMD` 指令支持

以上特性为网络数据包的编解码提供便利。

##### benchmark

https://github.com/lysShub/reedsolomon-go/actions

##### example

```go
package main

import (
	"fmt"
	"slices"

	"github.com/lysShub/reedsolomon-go"
	"github.com/lysShub/reedsolomon-go/galois"
)

// go env -w GOEXPERIMENT=simd
func main() {
	var para = reedsolomon.Para{Groupsize: 5, Datasize: 3}

	var datas = make([][]byte, 3)
	datas[0] = []byte{1, 2, 3, 4}
	datas[1] = []byte{5, 6}
	datas[2] = []byte{7, 8, 9}
	var parity = genParityblocks(para, datas)

	{ // encode
		encode(para, datas[0], 0, parity)
		encode(para, datas[1], 1, parity)
		encode(para, datas[2], 2, parity)
	}
	var blocks = append(datas, parity...)

	// mock packet loss, lost second data-block and first parity-block
	blocks = slices.Delete(blocks, 3, 4)
	blocks = slices.Delete(blocks, 1, 2)

	{ // reconst
		var reconst = [][]byte{
			make([]byte, len(blocks[len(blocks)-1])),
		}
		n := reedsolomon.Reconst(para, blocks, []uint8{0, 2, 4}, reconst)
		if n != 1 {
			panic(n)
		}

		i := galois.LastIndex(reconst[0], tailByte)
		if i < 0 {
			panic(i)
		}
		fmt.Println(reconst[0][:i])
		// [5 6]
	}
}

func genParityblocks(para reedsolomon.Para, datablocks [][]byte) [][]byte {
	var maxDatablock int
	for _, e := range datablocks {
		maxDatablock = max(maxDatablock, len(e))
	}
	maxDatablock += 1 // tail-byte

	var p = [][]byte{}
	for range para.Paritysize() {
		p = append(p, make([]byte, maxDatablock))
	}
	return p
}

const tailByte = 0xff

func encode(para reedsolomon.Para, data []byte, idx uint8, parity [][]byte) {
	data = append(data, tailByte)
	reedsolomon.Encode(para, data, idx, parity)
}
```

##### reference

https://github.com/klauspost/reedsolomon

https://github.com/templexxx/reedsolomon

https://vearne.cc/archives/39331
