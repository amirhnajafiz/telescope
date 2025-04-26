package controllers

import (
	"fmt"

	"github.com/hare1039/go-mpd"
	"go.uber.org/zap"
)

// MPDBuilder is a structure that handles the building of MPD files
type MPDBuilder struct {
	Logr *zap.Logger
}

// NewMPDBuilder creates a new instance of MPDBuilder
func NewMPDBuilder(logr *zap.Logger) *MPDBuilder {
	return &MPDBuilder{
		Logr: logr,
	}
}

// Build creates a new MPD file with the given CID based on the original MPD
func (p *MPDBuilder) Build(original []byte, cid string) ([]byte, error) {
	p.Logr.Info("building MPD", zap.String("cid", cid))

	// create a copy of the original MPD
	tree := new(mpd.MPD)
	if err := tree.Decode(original); err != nil {
		return nil, err
	}

	// set the media and initialization paths
	initPath := fmt.Sprintf("/api/%s/stream/init-stream$RepresentationID$.m4s", cid)
	mediaPath := fmt.Sprintf("/api/%s/stream/chunk-stream$RepresentationID$-$Number%%05d$.m4s", cid)

	// for each period, rewrite the SegmentTemplate and Representation paths
	for _, period := range tree.Period {
		for _, adapt := range period.AdaptationSets {
			// rewrite SegmentTemplate paths
			if adapt.SegmentTemplate != nil {
				adapt.SegmentTemplate.Media = &mediaPath
				adapt.SegmentTemplate.Initialization = &initPath
			} else {
				adapt.SegmentTemplate = &mpd.SegmentTemplate{
					Media:          &mediaPath,
					Initialization: &initPath,
				}
			}

			// rewrite Representation paths
			for i := range adapt.Representations {
				rep := &adapt.Representations[i]
				if rep.SegmentTemplate != nil {
					rep.SegmentTemplate.Media = &mediaPath
					rep.SegmentTemplate.Initialization = &initPath
				} else {
					rep.SegmentTemplate = &mpd.SegmentTemplate{
						Media:          &mediaPath,
						Initialization: &initPath,
					}
				}
			}
		}
	}

	return tree.Encode()
}
