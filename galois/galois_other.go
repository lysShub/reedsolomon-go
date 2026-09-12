//go:build !amd64
// +build !amd64

package galois

func mul(a, b byte) byte { return mulGo(a, b) }

func mulVect(c byte, in, out []byte) { mulVectGo(c, in, out) }

func mulXorVect(c byte, in, out []byte) { mulXorVectGo(c, in, out) }

func xorVect(in, out []byte) { xorVectGo(in, out) }

func lastIndex(s []byte, v byte) int { return lastIndexGo(s, v) }
