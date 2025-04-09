package throughput

import (
	"testing"
	"time"
)

func TestEstimator(t *testing.T) {
	e := NewEstimator()

	client := "client-test"

	// Simulate 1 MB download in 1 sec (8 Mbps)
	e.RecordDownload(client, 1024*1024, time.Second, false)

	bw := e.GetBandwidth(client)

	if bw < 7.5e6 || bw > 8.5e6 {
		t.Errorf("expected ~8Mbps, got %.2f", bw)
	}
}
