package cache

import (
	"fmt"

	"github.com/amirhnajafiz/telescope/pkg/files"
)

// Cache is an interface for our cache system to store segment records by CID
type Cache struct {
	baseDir   string
	hitCount  int
	missCount int
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
	return files.Write(fmt.Sprintf("%s/%s", c.baseDir, cid), data)
}

// Retrieve retrieves the data from the cache, if it does not exist, it returns an error
func (c *Cache) Retrieve(cid string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s", c.baseDir, cid)
	if !files.Exists(path) {
		c.missCount++
		return nil, fmt.Errorf("file %s does not exist", cid)
	}

	c.hitCount++

	return files.Read(path)
}

// Size returns the number of files in the cache directory
func (c *Cache) Size() (int, error) {
	return files.CountInDir(c.baseDir)
}

// GetHitCounts returns the number of cache hits
func (c *Cache) GetHitCounts() int {
	return c.hitCount
}

// GetMissCounts returns the number of cache misses
func (c *Cache) GetMissCounts() int {
	return c.missCount
}
