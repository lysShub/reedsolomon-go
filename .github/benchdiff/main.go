package main

import (
	"fmt"
	"os"

	"golang.org/x/tools/benchmark/parse"
)

func main() {
	base := load(os.Args[1])
	head := load(os.Args[2])

	fail := false
	for name, ns := range base {
		hns, ok := head[name]
		if !ok {
			fmt.Printf("DELETED: %s\n", name)
			fail = true
			continue
		}
		delta := (hns - ns) / ns * 100
		fmt.Printf("%s  %+.2f%%\n", name, delta)
		if delta > 5.0 {
			fail = true
		}
	}
	if fail {
		os.Exit(1)
	}
}

func load(path string) map[string]float64 {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s, err := parse.ParseSet(f)
	if err != nil {
		panic(err)
	}
	m := make(map[string]float64)
	for name, bs := range s {
		if len(bs) > 0 {
			m[name] = bs[0].NsPerOp
		}
	}
	return m
}