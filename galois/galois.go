// Package galois 定义伽罗瓦域 GF(2⁸) 中的运算
package galois

import (
	"acceler/pkg/debug"

	_ "golang.org/x/sys/cpu"
)

func Div(a, b byte) byte {
	if b == 0 {
		panic("divide by zero")
	}
	switch a {
	case 0:
		return 0
	case 1:
		logResult := logTable[b] ^ 255
		return expTable[logResult]
	default:
		logA := int(logTable[a])
		logB := int(logTable[b])
		logResult := logA - logB
		if logResult < 0 {
			logResult += 255
		}
		return expTable[uint8(logResult)]
	}
}

// Exp a**n.
func Exp(a byte, n int) byte {
	if n == 0 {
		return 1
	}
	if a == 0 {
		return 0
	}

	logA := logTable[a]
	logResult := int(logA) * n
	for logResult >= 255 {
		logResult -= 255
	}
	return expTable[uint8(logResult)]
}

func Mul(a, b byte) byte { return mul(a, b) }
func Add(a, b byte) byte { return a ^ b }

// MulVect out[i] = c * in[i], 要求 len(out)>=len(in)
func MulVect(c byte, in, out []byte) {
	if debug.Debug() {
		debug.GreaterOrEqual(len(out), len(in))
	}
	mulVect(c, in, out)
}

// MulXorVect  out[i] = (c*in[i]) ^ out[i], 要求 len(out)>=len(in)
func MulXorVect(c byte, in, out []byte) {
	if debug.Debug() {
		debug.GreaterOrEqual(len(out), len(in))
	}
	mulXorVect(c, in, out)
}

func LastIndex(s []byte, b byte) int { return lastIndex(s, b) }
