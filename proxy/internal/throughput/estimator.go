package throughput

import (
	"sync"
	"time"
)

type ClientEstimate struct {
	Uncached float64
	Cached   float64
	CurBW    float64
	Alpha    float64 // smoothing factor
}

type Estimator struct {
	clients map[string]*ClientEstimate
	mu      sync.RWMutex
}

// NewEstimator returns a fresh estimator
func NewEstimator() *Estimator {
	return &Estimator{
		clients: make(map[string]*ClientEstimate),
	}
}

// RecordDownload updates the throughput estimate for a client
func (e *Estimator) RecordDownload(clientID string, size int, duration time.Duration, cached bool) {
	bw := float64(size*8) / duration.Seconds() // bits/sec

	e.mu.Lock()
	defer e.mu.Unlock()

	est, exists := e.clients[clientID]
	if !exists {
		est = &ClientEstimate{
			Cached:   bw,
			Uncached: bw,
			CurBW:    bw,
			Alpha:    0.5,
		}
		e.clients[clientID] = est
	}

	if cached {
		est.Cached = est.Alpha*est.Cached + (1-est.Alpha)*bw
	} else {
		est.Uncached = est.Alpha*est.Uncached + (1-est.Alpha)*bw
	}

	// Blend CurBW based on ratio (could be adjusted later)
	est.CurBW = (est.Cached + est.Uncached) / 2
}

// GetBandwidth returns the estimated blended bandwidth
func (e *Estimator) GetBandwidth(clientID string) float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if est, ok := e.clients[clientID]; ok {
		return est.CurBW
	}
	return 0.0
}

// Get cached throughput Tg
// How fast the gateway delivers segments
func (e *Estimator) GetUncached(clientID string) float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if est, ok := e.clients[clientID]; ok {
		return est.Uncached
	}
	return 0
}

// Get uncached throughput Tc
// How fast IPFS providers deliver segments
func (e *Estimator) GetCached(clientID string) float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if est, ok := e.clients[clientID]; ok {
		return est.Cached
	}
	return 0
}
