package controllers

import (
	"fmt"
	"sync"

	"github.com/amirhnajafiz/telescope/internal/storage/cache"

	"github.com/hare1039/go-mpd"
	"go.uber.org/zap"
)

// AbrRewriter is a structure that rewrites the MPD file based on the current bandwidth
type AbrRewriter struct {
	Cache *cache.Cache
	Logr  *zap.Logger

	ipfsBW  float64
	ipfsRTT float64
	lock    sync.Mutex
}

// NewAbrRewriter creates a new instance of AbrRewriter
func NewAbrRewriter(cache *cache.Cache, logr *zap.Logger) *AbrRewriter {
	return &AbrRewriter{
		Cache:   cache,
		Logr:    logr,
		ipfsBW:  0,
		ipfsRTT: 0,
		lock:    sync.Mutex{},
	}
}

// SetIpfsBandwidth sets the IPFS bandwidth
func (p *AbrRewriter) SetIpfsBandwidth(bw float64) {
	p.lock.Lock()
	p.ipfsBW = bw
	p.lock.Unlock()
}

// GetIpfsBandwidth returns the IPFS RTT
func (p *AbrRewriter) SetIpfsRTT(rtt float64) {
	p.lock.Lock()
	p.ipfsRTT = rtt
	p.lock.Unlock()
}

// RewriteMPD rewrites the MPD file based on the current bandwidth and cache status
func (p *AbrRewriter) RewriteMPD(original []byte, cid string, ebw float64) ([]byte, error) {
	p.Logr.Info("rewriting MPD", zap.String("cid", cid), zap.Float64("bandwidth", ebw))

	// create a copy of the original MPD
	tree := new(mpd.MPD)
	if err := tree.Decode(original); err != nil {
		return nil, err
	}

	// calculate the bandwidth
	for _, period := range tree.Period {
		for _, adapt := range period.AdaptationSets {
			for i := range adapt.Representations {
				rep := &adapt.Representations[i]
				bw := float64(*rep.Bandwidth)

				// adjust bandwidth based on cache status
				var newBw float64

				fmt.Println("Cache status for representation ID:", *rep.ID)

				if p.Cache.Exists(*rep.ID) {
					// if cached, set bandwidth to the minimum of actual bandwidth and estimated bandwidth
					newBw = min(bw, ebw)
				} else {
					// if not cached, involve IPFS bandwidth and RTT
					if p.ipfsRTT > 0 {
						ipfsAdjustedBw := p.ipfsBW / p.ipfsRTT
						newBw = max(bw, ipfsAdjustedBw)
					} else {
						newBw = bw // fallback to actual bandwidth if RTT is zero
					}
				}

				// ensure bandwidth is at least 1
				newBw = max(newBw, 1)
				tmp := uint64(newBw)
				rep.Bandwidth = &tmp
			}
		}
	}

	return tree.Encode()
}
