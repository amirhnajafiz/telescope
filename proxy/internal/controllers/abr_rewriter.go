package controllers

import (
	"fmt"
	"log"

	"github.com/amirhnajafiz/telescope/internal/storage/cache"

	"github.com/hare1039/go-mpd"
)

type AbrRewriter struct {
	Estimator *Estimator
	Cache     *cache.Cache
}

func (p *AbrRewriter) RewriteMPD(original []byte, clientID string, cid string) ([]byte, error) {
	log.Printf("Rewriting MPD for client %s with CID %s", clientID, cid)

	tree := new(mpd.MPD)
	if err := tree.Decode(original); err != nil {
		return nil, err
	}

	Tc := p.Estimator.GetCurBW(clientID)
	Tg := p.Estimator.GetCached(clientID)
	Tn := p.Estimator.GetUncached(clientID)

	initPath := fmt.Sprintf("/api/contents/%s/init/stream", cid)
	mediaPath := fmt.Sprintf("/api/contents/%s/$Number$/stream", cid)

	for _, period := range tree.Period {
		for _, adapt := range period.AdaptationSets {
			// Rewrite SegmentTemplate paths
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
