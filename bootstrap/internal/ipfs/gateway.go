package ipfs

import (
	"fmt"

	shell "github.com/ipfs/go-ipfs-api"
)

// gateway is a struct that implements the interface for IPFS operations
type gateway struct {
	url string
}

// PutDIR uploads a directory to IPFS and returns the CID
func (g *gateway) PutDIR(path string) (string, error) {
	// connect to the IPFS node
	sh := shell.NewShell(g.url)

	// upload the directory
	cid, err := sh.AddDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to upload directory: %w", err)
	}

	return cid, nil
}
