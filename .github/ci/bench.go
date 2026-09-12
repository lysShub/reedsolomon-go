package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/tools/benchmark/parse"
)

type row struct {
	name  string
	delta float64
	speed float64
}

// runBench runs the bench job: benchmark base and head on the same runner,
// compare, retry the whole round on failure, and print a step-summary section.
func runBench(baseDir, headDir, name string, threshold float64, maxRetry int) bool {
	const baseTxt, headTxt = "base.txt", "head.txt"
	args := []string{"test", "-run=^$", "-bench=Benchmark", "-benchmem", "-count=1", "-benchtime=15s", "./..."}

	var (
		body string
		pass bool
	)
	for attempt := 0; attempt <= maxRetry; attempt++ {
		fmt.Fprintf(os.Stderr, "\n======================= bench attempt %d/%d =======================\n", attempt+1, maxRetry+1)
		if err := runToFile(baseDir, baseTxt, "go", args...); err != nil {
			fmt.Fprintf(os.Stderr, "bench base failed: %v\n", err)
			return false
		}
		if err := runToFile(headDir, headTxt, "go", args...); err != nil {
			fmt.Fprintf(os.Stderr, "bench head failed: %v\n", err)
			return false
		}
		body, pass = compareBench(baseTxt, headTxt, threshold)
		if pass {
			break
		}
	}
	emit("bench", name, body)
	return pass
}

func compareBench(basePath, headPath string, threshold float64) (string, bool) {
	base := load(basePath)
	head := load(headPath)

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
			rows = append(rows, row{name, +100, hns.mbps})
		case !hOk:
			rows = append(rows, row{name, -100, 0})
		default:
			rows = append(rows, row{name, (bns.ns - hns.ns) / bns.ns * 100, hns.mbps})
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].name < rows[j].name })

	max := 0
	for _, r := range rows {
		if len(r.name) > max {
			max = len(r.name)
		}
	}

	var b strings.Builder
	fail := false
	for _, r := range rows {
		fmt.Fprintf(&b, "%-*s  %10s  %+.2f%%\n", max, r.name, formatSpeed(r.speed), r.delta)
		fail = fail || r.delta < threshold
	}
	return b.String(), !fail
}

func formatSpeed(mbps float64) string {
	if mbps <= 0 {
		return "-"
	}
	b := mbps * 1e6 // B/s
	switch {
	case b >= 1e9:
		return fmt.Sprintf("%.1f GB/s", b/1e9)
	case b >= 1e6:
		return fmt.Sprintf("%.1f MB/s", b/1e6)
	case b >= 1e3:
		return fmt.Sprintf("%.1f KB/s", b/1e3)
	default:
		return fmt.Sprintf("%.1f B/s", b)
	}
}

type benchInfo struct {
	ns   float64
	mbps float64
}

func load(path string) map[string]benchInfo {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s, err := parse.ParseSet(f)
	if err != nil {
		panic(err)
	}
	m := make(map[string]benchInfo)
	for name, bs := range s {
		if len(bs) > 0 {
			m[name] = benchInfo{bs[0].NsPerOp, bs[0].MBPerS}
		}
	}
	return m
}
