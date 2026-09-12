package galois

import (
	"golang.org/x/sys/cpu"
)

var hasAVX2 = cpu.X86.HasAVX2
var _ = hasAVX2

// copy from https://github.com/golang/go/issues/36891
func _LastIndexByte(b []byte, c byte) int
