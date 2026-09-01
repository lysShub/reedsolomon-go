package galois

import (
	"bytes"
	"encoding/binary"
)

func mulVect_go(c byte, i, o []byte) {
	switch c {
	case 0:
		clear(o[:len(i)])
	case 1:
		copy(o[:len(i)], i)
	default:
		t := mulTable()[c]
		for idx := range i {
			o[idx] = t[i[idx]]
		}
	}
}

func mulXorVect_go(c byte, i, o []byte) {
	switch c {
	case 0:
		_ = o[:len(i)]
	case 1:
		xorVect_go(i, o)
	default:
		t := mulTable()[c]
		for idx := range i {
			o[idx] ^= t[i[idx]]
		}
	}
}

func xorVect_go(i, o []byte) {
	n := len(i)
	idx := 0
	for ; idx+8 <= n; idx += 8 {
		wi := binary.LittleEndian.Uint64(i[idx:])
		wo := binary.LittleEndian.Uint64(o[idx:])
		binary.LittleEndian.PutUint64(o[idx:], wi^wo)
	}
	for ; idx < n; idx++ {
		o[idx] ^= i[idx]
	}
}

func lastIndex_go(s []byte, v byte) int {
	return bytes.LastIndexByte(s, v)
}
