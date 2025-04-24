package estimator

import "time"

const (
	alpha = 0.5 // default value for alpha
)

// Estimate updates the throughput estimate for a client based on the size of the downloaded content, duration of the download, and whether it was cached or not
func Estimate(
	size int,
	duration time.Duration,
	cached bool,
	cachedBW,
	uncachedBW,
	curBW float64,
) (float64, float64, float64) {
	bw := float64(size*8000000) / float64(duration.Microseconds()) // bits/sec

	// set default values if curBW is 0
	if curBW == 0 {
		cachedBW = bw   // default value
		uncachedBW = bw // default value
		curBW = bw      // default value
	}

	if cached {
		cachedBW = alpha*cachedBW + (1-alpha)*bw
	} else {
		uncachedBW = alpha*uncachedBW + (1-alpha)*bw
	}

	// blend CurBW based on ratio (could be adjusted later)
	curBW = (cachedBW + uncachedBW) / 2

	return cachedBW, uncachedBW, curBW
}
