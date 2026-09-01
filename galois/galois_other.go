//go:build !amd64
// +build !amd64

package galois

func init() {
	mulVect = mulVect_go
	mulXorVect = mulXorVect_go
	lastIndex = lastIndex_go
	xorVect = xorVect_go
}
