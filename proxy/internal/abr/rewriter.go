package abr

import (
	"github.com/amirhnajafiz/telescope/internal/cache"
	"github.com/amirhnajafiz/telescope/internal/throughput"
	"github.com/hare1039/go-mpd"
)

type CacheBasedPolicy struct {
	Estimator *throughput.Estimator
	Cache     *cache.SegmentCache
}

func (p *CacheBasedPolicy) RewriteMPD(original []byte, clientID string) ([]byte, error) {
	// Parse original MPD XML
	tree := new(mpd.MPD)
	if err := tree.Decode(original); err != nil {
		return nil, err
	}

	Tc := p.Estimator.GetUncached(clientID)
	Tg := p.Estimator.GetCached(clientID)

	adjustment := Tc - Tg

	for _, period := range tree.Period {
		for _, adapt := range period.AdaptationSets {
			for i := range adapt.Representations {
				rep := &adapt.Representations[i]
				bw := float64(*rep.Bandwidth)

				// Apply adjustment
				newBw := bw
				if adjustment > 0 {
					newBw = bw + adjustment // penalize uncached
				} else {
					newBw = bw + adjustment // reward cached
				}

				if newBw < 1 {
					newBw = 1
				}

				tmp := uint64(newBw)
				rep.Bandwidth = &tmp
			}
		}
	}

	// Encode modified MPD back to bytes
	return tree.Encode()
}
