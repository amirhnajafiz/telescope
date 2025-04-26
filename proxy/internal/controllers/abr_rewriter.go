package controllers

import (
	"context"
	"sync"

	"github.com/amirhnajafiz/telescope/internal/storage/cache"
	"go.opentelemetry.io/otel/trace"

	"github.com/hare1039/go-mpd"
	"go.uber.org/zap"
)

const (
	alphaFactor       = 0.5 // adjustment factor for real and estimated bandwidth
	replicationFactor = 0.3 // replication factor for bandwidth adjustment based on replication
)

// AbrRewriter is a structure that rewrites the MPD file based on the current bandwidth
type AbrRewriter struct {
	Cache  *cache.Cache
	Logr   *zap.Logger
	Tracer trace.Tracer

	tn   float64 // ipfs network bandwidth
	tg   float64 // gateway network bandwidth
	lock sync.Mutex
}

// NewAbrRewriter creates a new instance of AbrRewriter
func NewAbrRewriter(cache *cache.Cache, logr *zap.Logger, trc trace.Tracer) *AbrRewriter {
	return &AbrRewriter{
		Cache:  cache,
		Logr:   logr,
		Tracer: trc,
		tn:     1,
		tg:     1,
		lock:   sync.Mutex{},
	}
}

// SetIpfsBandwidth sets the IPFS bandwidth
func (p *AbrRewriter) SetIpfsBandwidth(bw float64) {
	p.lock.Lock()
	p.tn = bw
	p.lock.Unlock()
}

// SetGatewayBandwidth sets the gateway bandwidth
func (p *AbrRewriter) SetGatewayBandwidth(bw float64) {
	p.lock.Lock()
	p.tg = bw
	p.lock.Unlock()
}

// RewriteMPD rewrites the MPD file based on the current bandwidth and cache status
func (p *AbrRewriter) RewriteMPD(
	ctx context.Context,
	original []byte,
	cid string,
	tc float64,
) ([]byte, error) {
	_, span := p.Tracer.Start(ctx, "abr-rewriter")
	defer span.End()

	p.Logr.Info("rewriting MPD", zap.String("cid", cid), zap.Float64("bandwidth", tc))

	// create a copy of the original MPD
	tree := new(mpd.MPD)
	if err := tree.Decode(original); err != nil {
		return nil, err
	}

	// change the bandwidth of the representations based on the cache status
	for _, period := range tree.Period {
		for _, adapt := range period.AdaptationSets {
			for i := range adapt.Representations {
				rep := &adapt.Representations[i]

				// construct the full Media path
				mediaTemplate := *adapt.SegmentTemplate.Media
				fullPath := constructFullPath(mediaTemplate, *rep.ID, i+1)

				// adjust bandwidth based on cache status and replication factor
				var (
					bw    = float64(*rep.Bandwidth)
					newBw float64
				)
				if p.Cache.Exists(fullPath) {
					p.Logr.Info("segment is cached", zap.String("path", fullPath))
					newBw = (tc-p.tg)*alphaFactor + bw*(1-alphaFactor)
				} else {
					p.Logr.Info("segment is not cached", zap.String("path", fullPath))
					newBw = (tc-p.tn)*alphaFactor + bw*(1-alphaFactor)
				}

				// apply probabilistic adjustment based on replication factor
				if shouldSelectReplica(replicationFactor) {
					p.Logr.Info("next replica selected for adjustment", zap.String("path", fullPath))
					newBw *= (1 + replicationFactor) // increase bandwidth slightly
				} else {
					p.Logr.Info("next replica not selected for adjustment", zap.String("path", fullPath))
					newBw *= (1 - replicationFactor) // decrease bandwidth slightly
				}

				// ensure bandwidth is at least 1 Mbps
				newBw = max(newBw, 1)

				tmp := uint64(newBw)
				rep.Bandwidth = &tmp
			}
		}
	}

	return tree.Encode()
}
