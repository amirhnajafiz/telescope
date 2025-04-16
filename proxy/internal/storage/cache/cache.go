package cache

// SegmentCache holds a thread-safe in-memory record of cached segment CIDs.
type SegmentCache struct {
	data map[string]bool
}

// NewCache creates a new cache instance
func NewCache() *SegmentCache {
	return &SegmentCache{
		data: make(map[string]bool),
	}
}

// IsCached returns true if the given CID is already cached.
func (c *SegmentCache) IsCached(cid string) bool {
	return c.data[cid]
}

// MarkCached marks a CID as cached.
func (c *SegmentCache) MarkCached(cid string) {
	c.data[cid] = true
}

// Size returns the number of currently cached items.
func (c *SegmentCache) Size() int {
	return len(c.data)
}
