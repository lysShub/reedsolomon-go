//go:build !amd64 && !arm64
// +build !amd64,!arm64

package galois

func initialize() {
	mulVect = mulVect_go
	mulXorVect = mulXorVect_go
	lastIndex = lastIndex_go
	xorVect = xorVect_go
}
