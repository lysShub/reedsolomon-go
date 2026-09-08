//go:build arm64
// +build arm64

package galois

import (
	"math/bits"
	"simd/archsimd"
	"unsafe"
)

func init() {
	mulVect = mulVect_smid
	mulXorVect = mulXorVect_smid
	lastIndex = lastIndex_smid
	xorVect = xorVect_smid
}

func mulVect_smid(c byte, i, o []byte) {
	n := len(i)
	switch c {
	case 0:
		clear(o[:n])
	case 1:
		copy(o[:n], i)
	default:
		iPtr := unsafe.Pointer(unsafe.SliceData(i))
		oPtr := unsafe.Pointer(unsafe.SliceData(o))
		if n >= 64 {
			t := &lohiTable()[c]
			loT := archsimd.LoadUint8x16Array(&t.lo)
			hiT := archsimd.LoadUint8x16Array(&t.hi)
			mask := archsimd.BroadcastUint8x16(0x0f)

			for n >= 64 {
				v0 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*0)))
				v1 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*1)))
				v2 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*2)))
				v3 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*3)))

				l0 := v0.And(mask)
				l1 := v1.And(mask)
				l2 := v2.And(mask)
				l3 := v3.And(mask)
				h0 := v0.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)
				h1 := v1.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)
				h2 := v2.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)
				h3 := v3.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)

				loT.LookupOrZero(l0).Xor(hiT.LookupOrZero(h0)).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*0)))
				loT.LookupOrZero(l1).Xor(hiT.LookupOrZero(h1)).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*1)))
				loT.LookupOrZero(l2).Xor(hiT.LookupOrZero(h2)).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*2)))
				loT.LookupOrZero(l3).Xor(hiT.LookupOrZero(h3)).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*3)))

				iPtr = unsafe.Add(iPtr, 64)
				oPtr = unsafe.Add(oPtr, 64)
				n -= 64
			}

			for n >= 16 {
				v0 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 0)))

				l0 := v0.And(mask)
				h0 := v0.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)

				loT.LookupOrZero(l0).Xor(hiT.LookupOrZero(h0)).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 0)))

				iPtr = unsafe.Add(iPtr, 16)
				oPtr = unsafe.Add(oPtr, 16)
				n -= 16
			}
		}

		tailI := unsafe.Slice((*uint8)(iPtr), n)
		tailO := unsafe.Slice((*uint8)(oPtr), n)
		t := _mulTable.Load()[c]
		for k := range tailI {
			tailO[k] = t[tailI[k]]
		}
	}
}
func mulXorVect_smid(c byte, i, o []byte) {
	n := len(i)
	switch c {
	case 0:
	case 1:
		xorVect_smid(i, o)
	default:
		iPtr := unsafe.Pointer(unsafe.SliceData(i))
		oPtr := unsafe.Pointer(unsafe.SliceData(o))
		if n >= 64 {
			t := &lohiTable()[c]
			loT := archsimd.LoadUint8x16Array(&t.lo)
			hiT := archsimd.LoadUint8x16Array(&t.hi)
			mask := archsimd.BroadcastUint8x16(0x0f)

			for n >= 64 {
				v0 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*0)))
				v1 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*1)))
				v2 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*2)))
				v3 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*3)))

				e0 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 16*0)))
				e1 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 16*1)))
				e2 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 16*2)))
				e3 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 16*3)))

				l0 := v0.And(mask)
				l1 := v1.And(mask)
				l2 := v2.And(mask)
				l3 := v3.And(mask)
				h0 := v0.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)
				h1 := v1.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)
				h2 := v2.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)
				h3 := v3.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)

				loT.LookupOrZero(l0).Xor(hiT.LookupOrZero(h0)).Xor(e0).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*0)))
				loT.LookupOrZero(l1).Xor(hiT.LookupOrZero(h1)).Xor(e1).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*1)))
				loT.LookupOrZero(l2).Xor(hiT.LookupOrZero(h2)).Xor(e2).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*2)))
				loT.LookupOrZero(l3).Xor(hiT.LookupOrZero(h3)).Xor(e3).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*3)))

				iPtr = unsafe.Add(iPtr, 64)
				oPtr = unsafe.Add(oPtr, 64)
				n -= 64
			}

			for n >= 16 {
				v0 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 0)))
				e0 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 0)))

				l0 := v0.And(mask)
				h0 := v0.ReshapeToUint64s().ShiftAllRight(4).ReshapeToUint8s().And(mask)

				loT.LookupOrZero(l0).Xor(hiT.LookupOrZero(h0)).Xor(e0).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 0)))

				iPtr = unsafe.Add(iPtr, 16)
				oPtr = unsafe.Add(oPtr, 16)
				n -= 16
			}

		}

		tailI := unsafe.Slice((*uint8)(iPtr), n)
		tailO := unsafe.Slice((*uint8)(oPtr), n)
		t := _mulTable.Load()[c]
		for k := range tailI {
			tailO[k] ^= t[tailI[k]]
		}
	}
}
func xorVect_smid(i, o []byte) {
	n := len(i)
	iPtr := unsafe.Pointer(unsafe.SliceData(i))
	oPtr := unsafe.Pointer(unsafe.SliceData(o))

	for n >= 64 {
		i0 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*0)))
		i1 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*1)))
		i2 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*2)))
		i3 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(iPtr, 16*3)))

		o0 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 16*0)))
		o1 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 16*1)))
		o2 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 16*2)))
		o3 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(oPtr, 16*3)))

		i0.Xor(o0).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*0)))
		i1.Xor(o1).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*1)))
		i2.Xor(o2).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*2)))
		i3.Xor(o3).StoreArray((*[16]uint8)(unsafe.Add(oPtr, 16*3)))

		iPtr = unsafe.Add(iPtr, 64)
		oPtr = unsafe.Add(oPtr, 64)
		n -= 64
	}

	for n >= 16 {
		i := archsimd.LoadUint8x16Array((*[16]uint8)(iPtr))
		o := archsimd.LoadUint8x16Array((*[16]uint8)(oPtr))
		i.Xor(o).StoreArray((*[16]uint8)(oPtr))
		iPtr = unsafe.Add(iPtr, 16)
		oPtr = unsafe.Add(oPtr, 16)
		n -= 16
	}

	tailI := unsafe.Slice((*uint8)(iPtr), n)
	tailO := unsafe.Slice((*uint8)(oPtr), n)
	for i, e := range tailI {
		tailO[i] ^= e
	}
}
func lastIndex_smid(s []byte, v byte) int {
	n := len(s)
	if n == 0 {
		return -1
	}
	src := unsafe.Pointer(unsafe.SliceData(s))
	ptr := unsafe.Add(src, n)
	broad := archsimd.BroadcastUint8x16(v)

	for n >= 64 {
		ptr = unsafe.Add(ptr, -64)
		n -= 64

		v0 := archsimd.LoadUint8x16Array((*[16]uint8)(ptr))
		v1 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(ptr, 16*1)))
		v2 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(ptr, 16*2)))
		v3 := archsimd.LoadUint8x16Array((*[16]uint8)(unsafe.Add(ptr, 16*3)))

		m0 := v0.Equal(broad)
		m1 := v1.Equal(broad)
		m2 := v2.Equal(broad)
		m3 := v3.Equal(broad)

		if m0.Or(m1).Or(m2).Or(m3).ToInt8x16().ConvertToUint8().ReduceSum() != 0 {
			if m := maskToBits(m3); m != 0 {
				return int(uintptr(ptr)-uintptr(src)) + 16*3 + bits.Len16(m) - 1
			}
			if m := maskToBits(m2); m != 0 {
				return int(uintptr(ptr)-uintptr(src)) + 16*2 + bits.Len16(m) - 1
			}
			if m := maskToBits(m1); m != 0 {
				return int(uintptr(ptr)-uintptr(src)) + 16*1 + bits.Len16(m) - 1
			}
			if m := maskToBits(m0); m != 0 {
				return int(uintptr(ptr)-uintptr(src)) + 16*0 + bits.Len16(m) - 1
			}
		}
	}

	for n >= 16 {
		ptr = unsafe.Add(ptr, -16)
		n -= 16
		v := archsimd.LoadUint8x16Array((*[16]uint8)(ptr))
		if m := maskToBits(v.Equal(broad)); m != 0 {
			return int(uintptr(ptr)-uintptr(src)) + bits.Len16(m) - 1
		}
	}

	for i := n - 1; i >= 0; i-- {
		if s[i] == v {
			return i
		}
	}
	return -1
}

func maskToBits(m archsimd.Mask8x16) uint16 {
	u := m.ToInt8x16().ConvertToUint8().ReshapeToUint64s()
	lo := u.GetElem(0)
	hi := u.GetElem(1)
	return uint16(((lo&0x0101010101010101)*0x0102040810204080)>>56 | ((hi&0x0101010101010101)*0x0102040810204080)>>56<<8)
}
