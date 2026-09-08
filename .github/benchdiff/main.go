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
}

func main() {
	base := load(os.Args[1])
	head := load(os.Args[2])

	names := make(map[string]bool, len(base)+len(head))
	for n := range base {
		names[n] = true
	}
	for n := range head {
		names[n] = true
	}

	rows := make([]row, 0, len(names))
	for name := range names {
		bns, bOk := base[name]
		hns, hOk := head[name]
		switch {
		case !bOk:
			rows = append(rows, row{name, +100})
		case !hOk:
			rows = append(rows, row{name, -100})
		default:
			rows = append(rows, row{name, (bns - hns) / bns * 100})
		}
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
		fmt.Printf("%-*s  %+.2f%%\n", max, r.name, r.delta)
		fail = fail || r.delta < -5.0
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
