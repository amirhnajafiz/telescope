package ipfs

import (
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	FetchSegment(cid string) ([]byte, error)
	FetchMPD(cid string) ([]byte, error)
}

// GatewayClient implements the IPFS Client interface using HTTP Gateway
type GatewayClient struct {
	BaseURL string //e.g., http://127.0.0.1:8081/ipfs
}

// HTTP GET to Kubo gateway
func (g *GatewayClient) FetchSegment(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", g.BaseURL, path)
	fmt.Printf("URL Segment: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Fetching mpd
func (g *GatewayClient) FetchMPD(cid string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/bunny.mpd", g.BaseURL, cid)
	fmt.Printf("MPD URL: %s/%s/bunny.mpd\ncid: %s\n", g.BaseURL, cid, cid)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to fetch MPD from IPFS: %s\nBody:\n%s\n", resp.Status, string(body))
		return nil, fmt.Errorf("non-200 from IPFS: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
