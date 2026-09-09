// Package galois 定义伽罗瓦域 GF(2⁸) 中的运算
package galois

import (
	"bytes"
	"encoding/binary"

	_ "golang.org/x/sys/cpu"
)

func Mul(a, b byte) byte {
	t := lohiTable()[a]
	return t.lo[b&0x0f] ^ t.hi[b>>4]
}
func Add(a, b byte) byte { return a ^ b }

func Div(a, b byte) byte {
	if b == 0 {
		panic("divide by zero")
	}
	switch a {
	case 0:
		return 0
	case 1:
		logR := logTable[b] ^ 255
		if logR == 255 {
			logR = 0
		}
		return expTable[logR]
	default:
		logA := int(logTable[a])
		logB := int(logTable[b])
		logR := logA - logB
		if logR < 0 {
			logR += 255
		}
		return expTable[uint8(logR)]
	}
}

// Pow a**n.
func Pow(a, n byte) byte {
	if n == 0 {
		return 1
	} else if a == 0 {
		return 0
	} else {
		logA := int(logTable[a]) * int(n)
		return expTable[uint8(logA%255)]
	}
}

var (
	mulVect    func(c byte, i, o []byte)
	mulXorVect func(c byte, i, o []byte)
	lastIndex  func(s []byte, v byte) int
	xorVect    func(i, o []byte)
)

func MulVect(c byte, i, o []byte)    { mulVect(c, i, o) }
func MulXorVect(c byte, i, o []byte) { mulXorVect(c, i, o) }
func XorVect(i, o []byte)            { xorVect(i, o) }
func LastIndex(s []byte, v byte) int { return lastIndex(s, v) }

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
