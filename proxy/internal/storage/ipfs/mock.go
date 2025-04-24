package ipfs

import (
	"fmt"
	"os"
)

const dirPath = "bp/idp"

// mock is a struct that implements the interface for mock-IPFS operations
type mock struct{}

// Get retrieves data from IPFS using the provided CID
func (m *mock) Get(cid string) ([]byte, error) {
	// build the file path
	path := fmt.Sprintf("%s/%s", dirPath, cid)

	// read the data from the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
