package ipfs

// Client is an interface for IPFS client operations
type Client interface {
	Get(cid string) ([]byte, int64, error)
}

// NewClient creates a new IPFS gateway instance
func NewClient(gatewayURL string) Client {
	if gatewayURL == "mock" {
		return &mock{}
	}

	return &gateway{
		url: gatewayURL,
	}
}
