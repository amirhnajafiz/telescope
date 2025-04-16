package cache

// Cache is an interface for our cache system to store segment records by CID
type Cache struct {
	baseDir string
}

// NewCache creates a new cache instance
func NewCache(bd string) *Cache {
	return &Cache{
		baseDir: bd,
	}
}

// IsCached returns true if the given CID is already cached.
func (c *Cache) IsCached(cid string) bool {
	return c.data[cid]
}

// MarkCached marks a CID as cached.
func (c *Cache) MarkCached(cid string) {
	c.data[cid] = true
}

// Size returns the number of currently cached items.
func (c *Cache) Size() int {
	return len(c.data)
}
