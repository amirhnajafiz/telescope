package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// serveFile handles the streaming of content over DASH
func (a *API) serveFile(
	ctx *fiber.Ctx,
	cid,
	filename,
	cacheKey,
	clientId string,
) error {
	var cached bool

	// check if the segment is cached
	segment, err := a.Cache.Retrieve(cacheKey)
	if err != nil {
		cached = false
		a.Logr.Warn("cache miss", zap.String("cid", cid), zap.String("filename", filename), zap.Error(err))
		a.Metrics.CacheMisses.Inc()
	} else {
		cached = true
		a.Logr.Info("cache hit", zap.String("cid", cid), zap.String("filename", filename))
		a.Metrics.CacheHits.Inc()
	}

	// calculate cache ratio
	total := float64(a.Cache.GetHitCounts() + a.Cache.GetMissCounts())
	if total > 0 {
		ratio := float64(a.Cache.GetHitCounts()) / total
		a.Metrics.CacheRatio.Set(ratio)
	}

	// fetch the segment from IPFS if not cached
	start := time.Now()
	if segment == nil {
		segment, err = a.IPFS.Get(fmt.Sprintf("%s/%s", cid, filename))
		if err != nil {
			a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), ctx.Path()).Inc()
			return ctx.Status(fiber.StatusBadGateway).SendString("fetch failed")
		}
	}
	duration := time.Since(start)

	// record the download in the ABR rewriter
	a.ABRRewriter.Estimator.RecordDownload(clientId, len(segment), duration, cached)

	// cache the segment
	if !cached {
		if err := a.Cache.Store(cacheKey, segment); err != nil {
			a.Logr.Error(
				"failed to store segment in cache",
				zap.String("cid", cid),
				zap.String("filename", filename),
				zap.Error(err),
			)
		}
	}

	a.Metrics.BytesTransferred.WithLabelValues(ctx.Method(), ctx.Path()).Add(float64(len(segment)))

	ctx.Set("Content-Type", "video/mp4")

	return ctx.Send(segment)
}
