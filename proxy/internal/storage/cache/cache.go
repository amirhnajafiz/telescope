package cache

import "fmt"

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

// Store stores the data in the cache under the given CID
func (c *Cache) Store(cid string, data []byte) error {
	return writeToFile(fmt.Sprintf("%s/%s", c.baseDir, cid), data)
}

// Retrieve retrieves the data from the cache, if it does not exist, it returns an error
func (c *Cache) Retrieve(cid string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s", c.baseDir, cid)
	if !isFileExists(path) {
		return nil, fmt.Errorf("file %s does not exist", cid)
	}

	return readFromFile(path)
}

// Size returns the number of files in the cache directory
func (c *Cache) Size() (int, error) {
	return countFilesInDir(c.baseDir)
}
