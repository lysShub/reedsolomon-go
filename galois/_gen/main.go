package main

import (
	"fmt"
	"os"

	"github.com/mmcloughlin/avo/build"
)

func main() {
	if len(os.Args) > 2 {
		panic("require target")
	}

	switch os.Args[1] {
	case "amd64":
		os.Args = append(os.Args[:1], "-out", "../galois_amd64.s", "-pkg=galois")
		amd64()
	default:
		panic(fmt.Sprintf("unknown target %s", os.Args[1]))
	}

	build.Generate()
}
