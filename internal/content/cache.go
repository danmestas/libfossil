package content

import (
	"container/list"
	"sync"

	libfossil "github.com/danmestas/libfossil/internal/fsltype"
	"github.com/danmestas/libfossil/db"
)

// Cache is a concurrency-safe LRU cache for expanded blob content.
// It reduces redundant delta-chain walks by caching the fully-expanded
// result of [Expand], keyed by rid.
//
// A nil *Cache is valid and acts as a passthrough to [Expand].
type Cache struct {
	mu      sync.Mutex
	items   map[libfossil.FslID]*list.Element
	order   *list.List
	curSize int64
	maxSize int64
	hits    int64
	misses  int64
}

type cacheEntry struct {
	rid  libfossil.FslID
	data []byte
}

// CacheStats reports cache hit/miss statistics and current memory usage.
type CacheStats struct {
	Hits    int64
	Misses  int64
	Size    int64
	MaxSize int64
	Entries int
}

// NewCache creates a cache bounded by maxBytes of expanded content.
func NewCache(maxBytes int64) *Cache {
	if maxBytes <= 0 {
		panic("content.NewCache: maxBytes must be > 0")
	}
	return &Cache{
		items:   make(map[libfossil.FslID]*list.Element),
		order:   list.New(),
		maxSize: maxBytes,
	}
}

// Expand returns the expanded content for rid, serving from cache when possible.
// On a cache miss, it delegates to the package-level [Expand] and caches the result.
// A nil receiver falls through to [Expand] directly.
func (c *Cache) Expand(q db.Querier, rid libfossil.FslID) ([]byte, error) {
	if c == nil {
		return Expand(q, rid)
	}

	c.mu.Lock()
	if elem, ok := c.items[rid]; ok {
		c.order.MoveToFront(elem)
		data := elem.Value.(*cacheEntry).data
		c.hits++
		c.mu.Unlock()
		out := make([]byte, len(data))
		copy(out, data)
		return out, nil
	}
	c.misses++
	c.mu.Unlock()

	data, err := Expand(q, rid)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Another goroutine may have cached it while we were expanding.
	if elem, ok := c.items[rid]; ok {
		c.order.MoveToFront(elem)
		return data, nil
	}

	cached := make([]byte, len(data))
	copy(cached, data)

	elem := c.order.PushFront(&cacheEntry{rid: rid, data: cached})
	c.items[rid] = elem
	c.curSize += int64(len(cached))

	for c.curSize > c.maxSize && c.order.Len() > 0 {
		c.evictOldest()
	}

	return data, nil
}

func (c *Cache) evictOldest() {
	back := c.order.Back()
	if back == nil {
		return
	}
	e := back.Value.(*cacheEntry)
	c.order.Remove(back)
	delete(c.items, e.rid)
	c.curSize -= int64(len(e.data))
}

// Invalidate removes a single rid from the cache.
func (c *Cache) Invalidate(rid libfossil.FslID) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[rid]; ok {
		e := elem.Value.(*cacheEntry)
		c.curSize -= int64(len(e.data))
		c.order.Remove(elem)
		delete(c.items, rid)
	}
}

// Clear removes all entries from the cache.
func (c *Cache) Clear() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[libfossil.FslID]*list.Element)
	c.order.Init()
	c.curSize = 0
}

// Stats returns a snapshot of cache statistics.
func (c *Cache) Stats() CacheStats {
	if c == nil {
		return CacheStats{}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return CacheStats{
		Hits:    c.hits,
		Misses:  c.misses,
		Size:    c.curSize,
		MaxSize: c.maxSize,
		Entries: c.order.Len(),
	}
}
