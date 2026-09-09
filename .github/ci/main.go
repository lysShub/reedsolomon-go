package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: ci <bench|cover> base.txt head.txt")
		os.Exit(2)
	}
	mode, base, head := os.Args[1], os.Args[2], os.Args[3]
	switch mode {
	case "bench":
		bench(base, head)
	case "cover":
		cover(base, head)
	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", mode)
		os.Exit(2)
	}
}
