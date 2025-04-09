package abr

import (
	"testing"
	"time"

	"github.com/amirhnajafiz/telescope/internal/cache"
	"github.com/amirhnajafiz/telescope/internal/throughput"
)

func TestRewriteMPD(t *testing.T) {
	sampleMPD := `
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" type="static">
	<Period>
		<AdaptationSet>
			<Representation id="480p" bandwidth="2000000">
				<SegmentTemplate media="video_480_$Number$.m4s" duration="2" timescale="1"/>
			</Representation>
			<Representation id="720p" bandwidth="4000000">
				<SegmentTemplate media="video_720_$Number$.m4s" duration="2" timescale="1"/>
			</Representation>
		</AdaptationSet>
	</Period>
</MPD>`

	cache := cache.NewCache()
	cache.MarkCached("480p")

	estimator := throughput.NewEstimator()
	client := "client1"
	// Assume CurBW = 3 Mbps, Cached = 2.5 Mbps, Uncached = 4 Mbps
	estimator.RecordDownload(client, 300000, 800*time.Millisecond, true)   // cached
	estimator.RecordDownload(client, 500000, 1000*time.Millisecond, false) // uncached

	policy := &CacheBasedPolicy{
		Estimator: estimator,
		Cache:     cache,
	}

	out, err := policy.RewriteMPD([]byte(sampleMPD), client)
	if err != nil {
		t.Fatalf("rewrite failed: %v", err)
	}

	if string(out) == sampleMPD {
		t.Errorf("expected MPD to be rewritten, but it wasn't")
	}
}
