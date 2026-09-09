//go:build amd64
// +build amd64

package galois

import (
	"simd/archsimd"
	"unsafe"

	"golang.org/x/sys/cpu"
)

func init() {
	if cpu.X86.HasAVX512GFNI {
		mulVect = mulVect_gfni
		mulXorVect = mulXorVect_gfni
		lastIndex = lastIndex_simd
		xorVect = xorVect_simd
	} else if cpu.X86.HasSSE2 {
		mulVect = mulVect_simd
		mulXorVect = mulXorVect_simd
		lastIndex = lastIndex_simd
		xorVect = xorVect_simd
	} else {
		mulVect = mulVect_go
		mulXorVect = mulXorVect_go
		lastIndex = lastIndex_go
		xorVect = xorVect_go
	}
}

func maskToBits(m archsimd.Mask8x16) uint16 { return m.ToBits() }
func maskNonZero(m archsimd.Mask8x16) bool  { return m.ToBits() != 0 }
func lookupOrZero(t archsimd.Uint8x16, idx archsimd.Uint8x16) archsimd.Uint8x16 {
	return t.PermuteOrZero(idx.BitsToInt8())
}

func mulVect_gfni(c byte, i, o []byte) {
	n := len(i)
	switch c {
	case 0:
		clear(o[:n])
	case 1:
		copy(o[:n], i)
	default:
		iPtr := unsafe.Pointer(unsafe.SliceData(i))
		oPtr := unsafe.Pointer(unsafe.SliceData(o))

		cVec := archsimd.BroadcastUint8x64(c)
		for n >= 256 {
			v0 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(iPtr, 64*0)))
			v1 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(iPtr, 64*1)))
			v2 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(iPtr, 64*2)))
			v3 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(iPtr, 64*3)))

			v0.GaloisFieldMul(cVec).StoreArray((*[64]uint8)(unsafe.Add(oPtr, 64*0)))
			v1.GaloisFieldMul(cVec).StoreArray((*[64]uint8)(unsafe.Add(oPtr, 64*1)))
			v2.GaloisFieldMul(cVec).StoreArray((*[64]uint8)(unsafe.Add(oPtr, 64*2)))
			v3.GaloisFieldMul(cVec).StoreArray((*[64]uint8)(unsafe.Add(oPtr, 64*3)))

			iPtr = unsafe.Add(iPtr, 256)
			oPtr = unsafe.Add(oPtr, 256)
			n -= 256
		}

		for n >= 64 {
			v0 := archsimd.LoadUint8x64Array((*[64]uint8)(iPtr))
			v0.GaloisFieldMul(cVec).StoreArray((*[64]uint8)(oPtr))

			iPtr = unsafe.Add(iPtr, 64)
			oPtr = unsafe.Add(oPtr, 64)
			n -= 64
		}

		tailI := unsafe.Slice((*uint8)(iPtr), n)
		tailO := unsafe.Slice((*uint8)(oPtr), n)
		t := _mulTable.Load()[c]
		for k := range tailI {
			tailO[k] = t[tailI[k]]
		}
	}
}

func mulXorVect_gfni(c byte, i, o []byte) {
	n := len(i)
	switch c {
	case 0:
	case 1:
		XorVect(i, o)
	default:
		iPtr := unsafe.Pointer(unsafe.SliceData(i))
		oPtr := unsafe.Pointer(unsafe.SliceData(o))

		cVec := archsimd.BroadcastUint8x64(c)
		for n >= 256 {
			v0 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(iPtr, 64*0)))
			v1 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(iPtr, 64*1)))
			v2 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(iPtr, 64*2)))
			v3 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(iPtr, 64*3)))

			e0 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(oPtr, 64*0)))
			e1 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(oPtr, 64*1)))
			e2 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(oPtr, 64*2)))
			e3 := archsimd.LoadUint8x64Array((*[64]uint8)(unsafe.Add(oPtr, 64*3)))

			v0.GaloisFieldMul(cVec).Xor(e0).StoreArray((*[64]uint8)(unsafe.Add(oPtr, 64*0)))
			v1.GaloisFieldMul(cVec).Xor(e1).StoreArray((*[64]uint8)(unsafe.Add(oPtr, 64*1)))
			v2.GaloisFieldMul(cVec).Xor(e2).StoreArray((*[64]uint8)(unsafe.Add(oPtr, 64*2)))
			v3.GaloisFieldMul(cVec).Xor(e3).StoreArray((*[64]uint8)(unsafe.Add(oPtr, 64*3)))

			iPtr = unsafe.Add(iPtr, 256)
			oPtr = unsafe.Add(oPtr, 256)
			n -= 256
		}

		for n >= 64 {
			v0 := archsimd.LoadUint8x64Array((*[64]uint8)(iPtr))
			e0 := archsimd.LoadUint8x64Array((*[64]uint8)(oPtr))
			v0.GaloisFieldMul(cVec).Xor(e0).StoreArray((*[64]uint8)(oPtr))

			iPtr = unsafe.Add(iPtr, 64)
			oPtr = unsafe.Add(oPtr, 64)
			n -= 64
		}

		tailI := unsafe.Slice((*uint8)(iPtr), n)
		tailO := unsafe.Slice((*uint8)(oPtr), n)
		t := _mulTable.Load()[c]
		for k := range tailI {
			tailO[k] ^= t[tailI[k]]
		}
	}
}
