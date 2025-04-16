package ipfs

type Client interface {
	FetchSegment(cid string) ([]byte, error)
	FetchMPD(cid string) ([]byte, error)
}

type GatewayClient struct {
	BaseURL string
}

// HTTP GET to Kubo gateway
func (g *GatewayClient) FetchSegment(cid string) ([]byte, error) {
	// placeholder for HTTP GET to Kubo gateway
	return nil, nil
}

// Fetching mpd
func (g *GatewayClient) FetchMPD(cid string) ([]byte, error) {
	return nil, nil
}
