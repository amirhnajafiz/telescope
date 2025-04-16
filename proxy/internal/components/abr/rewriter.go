package abr

import (
	"github.com/amirhnajafiz/telescope/internal/components/throughput"
	"github.com/amirhnajafiz/telescope/internal/storage/cache"

	"github.com/hare1039/go-mpd"
)

type CacheBasedPolicy struct {
	Estimator *throughput.Estimator
	Cache     *cache.Cache
}

func (p *CacheBasedPolicy) RewriteMPD(original []byte, clientID string) ([]byte, error) {
	tree := new(mpd.MPD)
	if err := tree.Decode(original); err != nil {
		return nil, err
	}

	// Tc := p.Estimator.GetCurBW(clientID)
	// Tg := p.Estimator.GetCached(clientID)
	// Tn := p.Estimator.GetUncached(clientID)

	for _, period := range tree.Period {
		for _, adapt := range period.AdaptationSets {
			for i := range adapt.Representations {
				rep := &adapt.Representations[i]
				bw := float64(*rep.Bandwidth)

				// Build a unique CID for each segment+quality (fake placeholder)
				// TODO parse segment base names from SegmentTemplate.Media or
				// track per-segment request progress (like Telescope originally did with LatestProgress per client)
				//cid := *rep.ID

				var adjustment float64

				// if p.Cache.IsCached(cid) {
				// 	adjustment = Tc - Tg
				// } else {
				// 	adjustment = Tc - Tn
				// }

				newBw := bw + adjustment
				if newBw < 1 {
					newBw = 1
				}

				tmp := uint64(newBw)
				rep.Bandwidth = &tmp
			}
		}
	}

	return tree.Encode()
}
