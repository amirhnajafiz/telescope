package ipfs

import (
	"bytes"

	shell "github.com/ipfs/go-ipfs-api"
)

// gateway is a struct that implements the interface for IPFS operations
type gateway struct {
	url string
}

// Get retrieves data from IPFS using the provided CID
func (g *gateway) Get(cid string) ([]byte, error) {
	return nil, nil
}

// Put uploads data to IPFS and returns the CID
func (g *gateway) Put(data []byte) (string, error) {
	// connect to the IPFS node
	sh := shell.NewShell(g.url)

	// add the data to IPFS
	cid, err := sh.Add(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	return cid, nil
}
