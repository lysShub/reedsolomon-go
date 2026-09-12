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
		master := fs.String("master", "../master", "master (base) checkout directory")
		merge := fs.String("merge", "..", "merge (head) checkout directory")
		name := fs.String("name", "", "summary title suffix (e.g. platform name)")
		threshold := fs.Float64("threshold", -2.0, "coverage delta threshold (%)")
		_ = fs.Parse(os.Args[2:])
		if !runTest(*master, *merge, *name, *threshold) {
			os.Exit(1)
		}
	case "bench":
		fs := flag.NewFlagSet("bench", flag.ExitOnError)
		master := fs.String("master", "../master", "master (base) checkout directory")
		merge := fs.String("merge", "..", "merge (head) checkout directory")
		name := fs.String("name", "", "summary title suffix (e.g. platform name)")
		benchtime := fs.String("benchtime", "5s", "benchmark time per case")
		threshold := fs.Float64("threshold", -10.0, "bench delta threshold (%)")
		maxRetry := fs.Int("max-retry", 3, "max retry attempts")
		_ = fs.Parse(os.Args[2:])
		if !runBench(*master, *merge, *name, *benchtime, *threshold, *maxRetry) {
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
