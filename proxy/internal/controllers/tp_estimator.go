package controllers

import (
	"sync"
	"time"

	"github.com/amirhnajafiz/telescope/pkg/models"
)

type Estimator struct {
	clients map[string]*models.Estimane
	mu      sync.RWMutex
}

// NewEstimator returns a fresh estimator
func NewEstimator() *Estimator {
	return &Estimator{
		clients: make(map[string]*models.Estimane),
	}
}

// RecordDownload updates the throughput estimate for a client
func (e *Estimator) RecordDownload(clientID string, size int, duration time.Duration, cached bool) {
	bw := float64(size*8) / duration.Seconds() // bits/sec

	e.mu.Lock()
	defer e.mu.Unlock()

	est, exists := e.clients[clientID]
	if !exists {
		est = &models.Estimane{
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

	// blend CurBW based on ratio (could be adjusted later)
	est.CurBW = (est.Cached + est.Uncached) / 2
}
