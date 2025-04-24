package ipfs

import (
	"io"

	shell "github.com/ipfs/go-ipfs-api"
)

// gateway is a struct that implements the interface for IPFS operations
type gateway struct {
	url string
}

// Get retrieves data from IPFS using the provided CID
func (g *gateway) Get(cid string) ([]byte, error) {
	// connect to the IPFS node
	sh := shell.NewShell(g.url)

	// get the data from IPFS
	reader, err := sh.Cat(cid)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// read the data from the reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
