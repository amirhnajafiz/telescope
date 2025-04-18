package ipfs

// Client is an interface for IPFS client operations
type Client interface {
	Get(cid string) ([]byte, error)
	Put(data []byte) (string, error)
}

// NewClient creates a new IPFS gateway instance
func NewClient(gatewayURL string) Client {
	return &gateway{
		url: gatewayURL,
	}
}
