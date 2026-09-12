package reedsolomon

import (
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/lysShub/reedsolomon-go/encodec"
)

var Cache = newCache((1024 * 1024) * 8)

type cache struct {
	limit int           // cache size limit
	count atomic.Uint64 // max count
	mu    sync.RWMutex
	bytes int
	m     map[cachekey]*cacheval
}

// update

func newCache(cacheSize int) *cache {
	if cacheSize < 1024 {
		panic(cacheSize)
	}
	return &cache{
		limit: cacheSize,
		m:     map[cachekey]*cacheval{},
	}
}

func (c *cache) Bytes() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bytes
}

func (c *cache) Release() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, e := range c.m {
		e.release()
	}
	clear(c.m)
}

func (c *cache) matrix(para Para, idxs ...uint8) []byte {
	key := cachekey{para: para, indexs: newIndexs(idxs)}
	c.mu.RLock()
	v := c.m[key]
	c.mu.RUnlock()
	if v != nil {
		cnt := v.count.Add(1)
		max := max(cnt, c.count.Load())
		c.count.Store(max)
		return v.m
	}

	// [newIndexs] is referenced underlying memory, read-only
	key.indexs = indexs(strings.Clone(string(key.indexs)))
	val := newCacheval(para, idxs)
	c.mu.Lock()
	{
		if v, ok := c.m[key]; ok {
			val.release()
			val = v
		} else {
			c.m[key] = val
			c.bytes += len(val.m)
			c.clear()
		}
	}
	c.mu.Unlock()
	return val.m
}
func (c *cache) clear() {
	if c.bytes > c.limit+c.limit/8 {
		limit := c.count.Load() / 8
		for k, v := range c.m {
			if v.count.Load() < limit {
				c.bytes -= len(v.m)
				v.release()
				delete(c.m, k)
			} else {
				v.count.Store(0)
			}
		}
		c.count.Store(0)
	}
}

type cacheval struct {
	m     []byte
	count atomic.Uint64
}

func newCacheval(para Para, idxs []uint8) *cacheval {
	buf := encodec.Pooler.Get(int(para.Groupsize) * int(para.Datasize))
	n := encodec.Encodec(buf, para.Groupsize, para.Datasize, idxs...)
	return &cacheval{m: buf[:n]}
}
func (c *cacheval) release() { encodec.Pooler.Put(c.m) }

type indexs string

func newIndexs(idxs []uint8) indexs {
	if len(idxs) == 0 {
		return ""
	}
	const word = unsafe.Sizeof(idxs[0])
	ptr := (*byte)(unsafe.Pointer(unsafe.SliceData(idxs)))
	return indexs(unsafe.String(ptr, len(idxs)*int(word)))
}

type cachekey struct {
	para   Para
	indexs indexs // only decodeMatrix
}
