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

func TestRecordDownload(t *testing.T) {
	e := NewEstimator()
	client := "client1"

	// 1MB in 1s = 8 Mbps
	e.RecordDownload(client, 1024*1024, time.Second, false)

	if uncached := e.GetUncached(client); uncached < 7.5e6 || uncached > 8.5e6 {
		t.Errorf("unexpected uncached bw: %.2f", uncached)
	}

	if cur := e.GetCurBW(client); cur < 7.5e6 || cur > 8.5e6 {
		t.Errorf("unexpected CurBW: %.2f", cur)
	}

	// Add cached download
	e.RecordDownload(client, 512*1024, 500*time.Millisecond, true)

	if cached := e.GetCached(client); cached == 0 {
		t.Errorf("expected non-zero cached bw")
	}
}
