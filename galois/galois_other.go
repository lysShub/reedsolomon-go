//go:build !amd64 && !arm64
// +build !amd64,!arm64

package galois

func init() {
	mulVect = mulVect_go
	mulXorVect = mulXorVect_go
	lastIndex = lastIndex_go
	xorVect = xorVect_go
}
