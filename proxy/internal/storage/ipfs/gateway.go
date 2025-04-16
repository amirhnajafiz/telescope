package ipfs

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
	return "", nil
}
