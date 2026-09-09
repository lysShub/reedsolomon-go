//go:build arm64
// +build arm64

package galois

import (
	"simd/archsimd"
)

func init() {
	mulVect = mulVect_simd
	mulXorVect = mulXorVect_simd
	lastIndex = lastIndex_simd
	xorVect = xorVect_simd
}

func maskToBits(m archsimd.Mask8x16) uint16 {
	u := m.ToInt8x16().ConvertToUint8().ReshapeToUint64s()
	lo := u.GetElem(0)
	hi := u.GetElem(1)
	return uint16(((lo&0x0101010101010101)*0x0102040810204080)>>56 | ((hi&0x0101010101010101)*0x0102040810204080)>>56<<8)
}
func maskNonZero(m archsimd.Mask8x16) bool { return m.ToInt8x16().ConvertToUint8().ReduceSum() != 0 }
func lookupOrZero(t archsimd.Uint8x16, idx archsimd.Uint8x16) archsimd.Uint8x16 {
	return t.LookupOrZero(idx)
}
