package main

import (
	"fmt"
	"os"
	"sort"

	"golang.org/x/tools/benchmark/parse"
)

type row struct {
	name  string
	delta float64
	ok    bool
}

func main() {
	base := load(os.Args[1])
	head := load(os.Args[2])

	rows := make([]row, 0, len(base))
	for name, ns := range base {
		hns, ok := head[name]
		rows = append(rows, row{name, (hns - ns) / ns * 100, ok})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].name < rows[j].name })

	max := 0
	for _, r := range rows {
		if len(r.name) > max {
			max = len(r.name)
		}
	}

	fail := false
	for _, r := range rows {
		if !r.ok {
			fmt.Printf("%-*s  DELETED\n", max, r.name)
			fail = true
			continue
		}
		fmt.Printf("%-*s  %+.2f%%\n", max, r.name, r.delta)
		if r.delta > 5.0 {
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
