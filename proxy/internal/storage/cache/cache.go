package cache

import "sync"

// SegmentCache holds a thread-safe in-memory record of cached segment CIDs.
type SegmentCache struct {
	data map[string]bool
	mu   sync.RWMutex
}

// NewCache creates a new cache instance
func NewCache() *SegmentCache {
	return &SegmentCache{
		data: make(map[string]bool),
	}
}

// IsCached returns true if the given CID is already cached.
func (c *SegmentCache) IsCached(cid string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.data[cid]
}

// MarkCached marks a CID as cached.
func (c *SegmentCache) MarkCached(cid string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[cid] = true
}

// Size returns the number of currently cached items.
func (c *SegmentCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}
