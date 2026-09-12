package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "test":
		fs := flag.NewFlagSet("test", flag.ExitOnError)
		base := fs.String("base", "../base", "base checkout directory")
		head := fs.String("head", "..", "head checkout directory")
		name := fs.String("name", "", "summary title suffix (e.g. platform name)")
		threshold := fs.Float64("threshold", -2.0, "coverage delta threshold (%)")
		_ = fs.Parse(os.Args[2:])
		if !runTest(*base, *head, *name, *threshold) {
			os.Exit(1)
		}
	case "bench":
		fs := flag.NewFlagSet("bench", flag.ExitOnError)
		base := fs.String("base", "../base", "base checkout directory")
		head := fs.String("head", "..", "head checkout directory")
		name := fs.String("name", "", "summary title suffix (e.g. platform name)")
		threshold := fs.Float64("threshold", -10.0, "bench delta threshold (%)")
		maxRetry := fs.Int("max-retry", 3, "max retry attempts")
		_ = fs.Parse(os.Args[2:])
		if !runBench(*base, *head, *name, *threshold, *maxRetry) {
			os.Exit(1)
		}
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: ci <test|bench> [flags]")
	os.Exit(2)
}
