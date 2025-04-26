package ipfs

import "time"

// Client is an interface for IPFS client operations
type Client interface {
	Get(cid string) ([]byte, time.Duration, error)
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
