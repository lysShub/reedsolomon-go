package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	coverpkg "golang.org/x/tools/cover"
)

func cover(basePath, headPath string) {
	base := loadCover(basePath)
	head := loadCover(headPath)

	names := make(map[string]bool, len(base)+len(head))
	for n := range base {
		names[n] = true
	}
	for n := range head {
		names[n] = true
	}

	type cRow struct {
		pkg   string
		base  string
		head  string
		delta float64
	}
	rows := make([]cRow, 0, len(names))
	for pkg := range names {
		b, bOk := base[pkg]
		h, hOk := head[pkg]
		switch {
		case !bOk:
			rows = append(rows, cRow{pkg, "-", fmt.Sprintf("%.2f%%", h), +100})
		case !hOk:
			rows = append(rows, cRow{pkg, fmt.Sprintf("%.2f%%", b), "-", -100})
		default:
			rows = append(rows, cRow{pkg, fmt.Sprintf("%.2f%%", b), fmt.Sprintf("%.2f%%", h), h - b})
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].pkg < rows[j].pkg })

	max := 0
	for _, r := range rows {
		if len(r.pkg) > max {
			max = len(r.pkg)
		}
	}

	fail := false
	for _, r := range rows {
		fmt.Printf("%-*s  %6s  %6s  %+.2f%%\n", max, r.pkg, r.base, r.head, r.delta)
		fail = fail || r.delta < -2.0
	}
	if fail {
		os.Exit(1)
	}
}

func loadCover(path string) map[string]float64 {
	profiles, err := coverpkg.ParseProfiles(path)
	if err != nil {
		panic(err)
	}

	type stat struct {
		covered int
		total   int
	}
	stats := make(map[string]*stat)
	for _, p := range profiles {
		pkg := p.FileName
		if slash := strings.LastIndexByte(pkg, '/'); slash >= 0 {
			pkg = pkg[:slash]
		}
		s := stats[pkg]
		if s == nil {
			s = &stat{}
			stats[pkg] = s
		}
		for _, b := range p.Blocks {
			s.total += b.NumStmt
			if b.Count > 0 {
				s.covered += b.NumStmt
			}
		}
	}
	m := make(map[string]float64, len(stats))
	for pkg, s := range stats {
		m[pkg] = float64(s.covered) / float64(s.total) * 100
	}
	return m
}
