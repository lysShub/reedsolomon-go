package reedsolomon

import (
	"bytes"
	"testing"

	"github.com/lysShub/reedsolomon-go/encodec"
)

func Test_Cache_Base(t *testing.T) {
	c := newCache(1024)
	para := Para{Groupsize: 8, Datasize: 4}

	// first lookup generates and caches the matrix
	m := c.matrix(para)
	if len(c.m) != 1 || c.Bytes() != len(m) {
		t.Fatalf("miss: entries=%d bytes=%d matrix=%d", len(c.m), c.Bytes(), len(m))
	}

	// second lookup hits and returns the same buffer
	before := c.Bytes()
	m2 := c.matrix(para)
	if len(c.m) != 1 || c.Bytes() != before {
		t.Fatalf("hit: entries=%d bytes=%d", len(c.m), c.Bytes())
	}
	if &m[0] != &m2[0] {
		t.Fatal("expected the same cached matrix")
	}

	// content matches a direct encode
	exp := make([]byte, int(para.Groupsize)*int(para.Datasize))
	n := encodec.Encodec(exp, para.Groupsize, para.Datasize)
	if !bytes.Equal(m, exp[:n]) {
		t.Fatal("cached matrix content mismatch")
	}

	// Release clears the cache
	c.Release()
	if len(c.m) != 0 {
		t.Fatalf("Release left %d entries", len(c.m))
	}
	if c.Bytes() != 0 {
		t.Fatalf("Release left bytes=%d", c.Bytes())
	}
}

func Test_Cache_Rotate(t *testing.T) {
	c := newCache(1024)

	// hot entry with a high use count
	hot := Para{Groupsize: 8, Datasize: 4}
	c.matrix(hot)
	for i := 0; i < 16; i++ {
		c.matrix(hot)
	}

	// fill cold entries until rotate evicts
	rotated := false
	for g := uint8(1); g <= 32 && !rotated; g++ {
		for d := uint8(1); d <= g; d++ {
			before := c.Bytes()
			c.matrix(Para{Groupsize: g, Datasize: d})
			if c.Bytes() < before { // rotate evicted cold entries
				rotated = true
				break
			}
		}
	}
	if !rotated {
		t.Fatal("rotate never triggered")
	}

	// the hot entry must survive
	n, b := len(c.m), c.Bytes()
	c.matrix(hot)
	if len(c.m) != n || c.Bytes() != b {
		t.Fatal("hot entry was evicted by rotate")
	}
}

func Test_CacheKey_Not_Reuse_Indexs(t *testing.T) {
	c := newCache(1024)
	para := Para{Groupsize: 8, Datasize: 4}
	idxs := []uint8{0, 1, 4, 5}

	c.matrix(para, idxs...)
	n, b := len(c.m), c.Bytes()

	for i := range idxs {
		idxs[i] = 200 + uint8(i)
	}

	c.matrix(para, []uint8{0, 1, 4, 5}...)
	if len(c.m) != n || c.Bytes() != b {
		t.Fatalf("cachekey aliased caller idxs: entries %d->%d bytes %d->%d", n, len(c.m), b, c.Bytes())
	}
}

func Test_Cache_Rotate_Keeps_NewRecord(t *testing.T) {
	c := newCache(1024)

	var triggering Para
	found := false
loop:
	for g := uint8(1); g <= 32; g++ {
		for d := uint8(1); d <= g; d++ {
			before := c.Bytes()
			c.matrix(Para{Groupsize: g, Datasize: d})
			if c.Bytes() < before { // rotate evicted old entries
				triggering = Para{Groupsize: g, Datasize: d}
				found = true
				break loop
			}
		}
	}
	if !found {
		t.Fatal("rotate never triggered")
	}

	n, b := len(c.m), c.Bytes()
	c.matrix(triggering) // must hit, not regenerate
	if len(c.m) != n || c.Bytes() != b {
		t.Fatalf("newly added entry was evicted by rotate: entries %d->%d bytes %d->%d", n, len(c.m), b, c.Bytes())
	}
}

func Test_Cache_Rotate_When_Count_TooSmall(t *testing.T) {
	c := newCache(1024)

	entries := 0
	for g := uint8(1); g <= 16; g++ {
		for d := uint8(1); d <= g; d++ {
			c.matrix(Para{Groupsize: g, Datasize: d})
			entries++
		}
	}

	if c.count.Load() != 0 {
		t.Fatalf("expected max count 0, got %d", c.count.Load())
	}
	// largest entry for g<=16 is (16-8)*8 = 64
	if c.Bytes() > c.limit+c.limit/8+64 {
		t.Fatalf("cache grew unbounded: entries=%d bytes=%d limit=%d", entries, c.Bytes(), c.limit)
	}
}
