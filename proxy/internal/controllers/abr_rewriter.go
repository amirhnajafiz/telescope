package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/amirhnajafiz/telescope/internal/storage/cache"

	"github.com/hare1039/go-mpd"
	"go.uber.org/zap"
)

// AbrRewriter is a structure that rewrites the MPD file based on the current bandwidth
type AbrRewriter struct {
	Cache *cache.Cache
	Logr  *zap.Logger

	tn   float64 // ipfs network bandwidth
	tg   float64 // gateway network bandwidth
	lock sync.Mutex
}

// NewAbrRewriter creates a new instance of AbrRewriter
func NewAbrRewriter(cache *cache.Cache, logr *zap.Logger) *AbrRewriter {
	return &AbrRewriter{
		Cache: cache,
		Logr:  logr,
		tn:    0,
		tg:    0,
		lock:  sync.Mutex{},
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
func (p *AbrRewriter) RewriteMPD(original []byte, cid string, tc float64) ([]byte, error) {
	p.Logr.Info("rewriting MPD", zap.String("cid", cid), zap.Float64("bandwidth", tc))

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

				// Construct the full Media path
				if adapt.SegmentTemplate != nil && adapt.SegmentTemplate.Media != nil {
					mediaTemplate := *adapt.SegmentTemplate.Media
					fullPath := p.constructFullPath(mediaTemplate, *rep.ID, 4) // Example: Number = 4
					p.Logr.Info("Full Media Path", zap.String("path", fullPath))

					// Check if the item is cached
					if p.Cache.Exists(fullPath) {
						p.Logr.Info("Item is cached", zap.String("path", fullPath))
					} else {
						p.Logr.Info("Item is not cached", zap.String("path", fullPath))
					}
				}

				// adjust bandwidth based on cache status
				var newBw float64
				if p.Cache.Exists(*rep.ID) {
					newBw = tc - p.tg
				} else {
					newBw = tc - p.tn
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

// constructFullPath replaces placeholders in the Media template with actual values
func (p *AbrRewriter) constructFullPath(template, representationID string, number int) string {
	// Replace $RepresentationID$ with the actual representation ID
	path := strings.ReplaceAll(template, "$RepresentationID$", representationID)

	// Replace $Number%05d$ with the formatted number
	numberPlaceholder := regexp.MustCompile(`\$Number%0(\d+)d\$`)
	path = numberPlaceholder.ReplaceAllStringFunc(path, func(match string) string {
		width, _ := strconv.Atoi(match[8 : len(match)-2]) // Extract width from %05d
		return fmt.Sprintf("%0*d", width, number)
	})

	// Remove "/stream" from the path
	path = strings.ReplaceAll(path, "/stream", "")

	// Construct the relative path by trimming the "/api/" prefix
	relativePath := strings.TrimPrefix(path, "/api/")

	return relativePath
}
