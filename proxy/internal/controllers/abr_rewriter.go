package controllers

import (
	"fmt"

	"github.com/amirhnajafiz/telescope/internal/storage/cache"

	"github.com/hare1039/go-mpd"
	"go.uber.org/zap"
)

// AbrRewriter is a structure that rewrites the MPD file based on the current bandwidth
type AbrRewriter struct {
	Cache *cache.Cache
	Logr  *zap.Logger
}

// NewAbrRewriter creates a new instance of AbrRewriter
func NewAbrRewriter(cache *cache.Cache, logr *zap.Logger) *AbrRewriter {
	return &AbrRewriter{
		Cache: cache,
		Logr:  logr,
	}
}

// RewriteMPD rewrites the MPD file based on the current bandwidth and cache status
func (p *AbrRewriter) RewriteMPD(
	original []byte,
	clientID string,
	cid string,
	curBW,
	cachedBW,
	uncachedBW float64,
) ([]byte, error) {
	p.Logr.Info("rewriting MPD", zap.String("clientId", clientID), zap.String("cid", cid))

	tree := new(mpd.MPD)
	if err := tree.Decode(original); err != nil {
		return nil, err
	}

	Tc := curBW
	Tg := cachedBW
	Tn := uncachedBW

	initPath := fmt.Sprintf("/api/contents/%s/init/stream", cid)
	mediaPath := fmt.Sprintf("/api/contents/%s/$Number$/stream", cid)

	for _, period := range tree.Period {
		for _, adapt := range period.AdaptationSets {
			// rewrite SegmentTemplate paths
			if adapt.SegmentTemplate != nil {
				adapt.SegmentTemplate.Media = &mediaPath
				adapt.SegmentTemplate.Initialization = &initPath
			}

			for i := range adapt.Representations {
				rep := &adapt.Representations[i]
				if rep.SegmentTemplate != nil {
					rep.SegmentTemplate.Media = &mediaPath
					rep.SegmentTemplate.Initialization = &initPath
				}
				bw := float64(*rep.Bandwidth)

				var adjustment float64
				if _, err := p.Cache.Retrieve(*rep.ID); err == nil {
					adjustment = Tc - Tg
				} else {
					adjustment = Tc - Tn
				}

				newBw := max(bw+adjustment, 1)
				tmp := uint64(newBw)
				rep.Bandwidth = &tmp
			}
		}
	}

	return tree.Encode()
}
