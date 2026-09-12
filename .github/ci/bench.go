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

func runBench(masterDir, mergeDir, name, benchtime string, threshold float64, maxRetry int) bool {
	const masterTxt, mergeTxt = "master.txt", "merge.txt"
	args := []string{"test", "-run=^$", "-bench=Benchmark", "-benchmem", "-count=1", "-benchtime=" + benchtime, "./..."}

	for attempt := 0; attempt <= maxRetry; attempt++ {
		fmt.Println()
		fmt.Println()
		fmt.Printf("::group::bench attempt %d/%d\n", attempt+1, maxRetry+1)

		var (
			body string
			over []string
			pass bool
		)
		ok := func() bool {
			if err := runToFile(masterDir, masterTxt, "go", args...); err != nil {
				fmt.Fprintf(os.Stderr, "bench master failed: %v\n", err)
				return false
			}
			if err := runToFile(mergeDir, mergeTxt, "go", args...); err != nil {
				fmt.Fprintf(os.Stderr, "bench merge failed: %v\n", err)
				return false
			}

			body, over, pass = compareBench(masterTxt, mergeTxt, threshold)
			fmt.Print(body)
			if len(over) > 0 {
				fmt.Println("\nover threshold:")
				for _, l := range over {
					fmt.Print(l)
				}
			}
			return true
		}()
		fmt.Println("::endgroup::")

		if !ok {
			return false
		}

		content := body
		if len(over) > 0 {
			content += "\nover threshold:\n" + strings.Join(over, "")
		}
		emit("bench", fmt.Sprintf("%s (try %d/%d)", name, attempt+1, maxRetry+1), content)
		if pass {
			return true
		}
	}
	return false
}

func compareBench(masterPath, mergePath string, threshold float64) (string, []string, bool) {
	master := load(masterPath)
	merge := load(mergePath)

	names := make(map[string]bool, len(master)+len(merge))
	for n := range master {
		names[n] = true
	}
	for n := range merge {
		names[n] = true
	}

	rows := make([]row, 0, len(names))
	for name := range names {
		mns, mOk := master[name]
		rns, rOk := merge[name]
		switch {
		case !mOk:
			rows = append(rows, row{name, +100, rns.mbps})
		case !rOk:
			rows = append(rows, row{name, -100, 0})
		default:
			rows = append(rows, row{name, (mns.ns - rns.ns) / mns.ns * 100, rns.mbps})
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].name < rows[j].name })

	max := 0
	for _, r := range rows {
		if len(r.name) > max {
			max = len(r.name)
		}
	}

	var (
		b    strings.Builder
		over []string
	)
	for _, r := range rows {
		line := fmt.Sprintf("%-*s  %10s  %+.2f%%\n", max, r.name, formatSpeed(r.speed), r.delta)
		fmt.Fprint(&b, line)
		if r.delta < threshold {
			over = append(over, line)
		}
	}
	return b.String(), over, len(over) == 0
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
