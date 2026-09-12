package galois

import "golang.org/x/sys/cpu"

func mul(a byte, b byte) byte {
	t := &lohiTable()[a]
	return t.lo[b&0x0f] ^ t.hi[b>>4]
}

//go:generate go run -C ./_gen . amd64
func _mulVect_16(lohi *lohi, in, out []byte)
func _mulXorVect_16(lohi *lohi, in, out []byte)
func _xorVect_16(in, out []byte)

func mulVect(c byte, in []byte, out []byte) {
	if cpu.X86.HasAVX {
		if len(in) < 16 {
			for i, e := range in {
				out[i] = mul(c, e)
			}
		} else {
			t := &lohiTable()[c]
			_mulVect_16(t, in, out)
		}
	} else {
		mulVectGo(c, in, out)
	}
}

func mulXorVect(c byte, in []byte, out []byte) {
	if cpu.X86.HasAVX {
		if len(in) < 16 {
			for i, e := range in {
				out[i] ^= mul(c, e)
			}
		} else {
			t := &lohiTable()[c]
			_mulXorVect_16(t, in, out)
		}
	} else {
		mulXorVectGo(c, in, out)
	}
}

func xorVect(in, out []byte) {
	if cpu.X86.HasAVX {
		if len(in) < 16 {
			for i, e := range in {
				out[i] ^= e
			}
		} else {
			_xorVect_16(in, out)
		}
	} else {
		xorVectGo(in, out)
	}
}

func lastIndex(s []byte, b byte) int {
	return _LastIndexByte(s, b)
}
