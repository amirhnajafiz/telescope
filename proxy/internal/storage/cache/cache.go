package cache

import (
	"fmt"
	"sync/atomic"
)

// Cache is an interface for our cache system to store segment records by CID
type Cache struct {
	baseDir   string
	hitCount  int64
	missCount int64
}

// NewCache creates a new cache instance
func NewCache(bd string) *Cache {
	return &Cache{
		baseDir:   bd,
		hitCount:  0,
		missCount: 0,
	}
}

// Store stores the data in the cache under the given CID
func (c *Cache) Store(cid string, data []byte) error {
	return write(fmt.Sprintf("%s/%s", c.baseDir, cid), data)
}

// Retrieve retrieves the data from the cache, if it does not exist, it returns an error
func (c *Cache) Retrieve(cid string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s", c.baseDir, cid)
	if !exists(path) {
		atomic.AddInt64(&c.missCount, 1)
		return nil, fmt.Errorf("file %s does not exist", cid)
	}

	atomic.AddInt64(&c.hitCount, 1)

	return read(path)
}

// Exists checks if the data exists in the cache
func (c *Cache) Exists(path string) bool {
	return exists(path)
}

// GetHitCounts returns the number of cache hits
func (c *Cache) GetHitCounts() int {
	return int(atomic.LoadInt64(&c.hitCount))
}

// GetMissCounts returns the number of cache misses
func (c *Cache) GetMissCounts() int {
	return int(atomic.LoadInt64(&c.missCount))
}
