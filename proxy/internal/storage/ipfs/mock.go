package ipfs

import (
	"fmt"
	"os"
	"time"
)

const dirPath = "bp/idp"

// mock is a struct that implements the interface for mock-IPFS operations
type mock struct{}

// Get retrieves data from IPFS using the provided CID
func (m *mock) Get(cid string) ([]byte, int64, error) {
	// build the file path
	path := fmt.Sprintf("%s/%s", dirPath, cid)

	start := time.Now()

	// read the data from the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}

	return data, time.Since(start).Milliseconds(), nil
}
