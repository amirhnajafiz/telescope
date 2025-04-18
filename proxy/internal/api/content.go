package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// getContent handles the download of content, used for getting mpd files
func (a *API) getContent(ctx *fiber.Ctx) error {
	// get the cid from the URL
	cid := ctx.Params("cid")

	// get clientId from the request header
	clientID := ctx.Get("X-Client-ID", "default")

	// fetch MPD file from IPFS
	mpd, err := a.IPFS.Get(cid)
	if err != nil {
		a.Logr.Error("failed to fetch mpd", zap.String("cid", cid), zap.Error(err))

		a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), ctx.Path()).Inc()

		return ctx.Status(fiber.StatusBadGateway).SendString("failed to fetch .mpd")
	}

	// rewrite MPD via ABR policy
	rewritten, err := a.ABRRewriter.RewriteMPD(mpd, clientID, cid)
	if err != nil {
		a.Logr.Error("failed to rewrite mpd", zap.String("cid", cid), zap.Error(err))

		a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), ctx.Path()).Inc()

		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("failed to rewrite manifest:\n %s", err))
	}

	a.Metrics.BytesTransferred.WithLabelValues(ctx.Method(), ctx.Path()).Add(float64(len(rewritten)))

	ctx.Set("Content-Type", "application/dash+xml")

	return ctx.Send(rewritten)
}

// listContents handles the listing of all contents
func (a *API) listContents(ctx *fiber.Ctx) error {
	return nil
}

// streamContent handles the streaming of content over DASH
func (a *API) streamContent(ctx *fiber.Ctx) error {
	// get the cid and segment from the URL
	cid := ctx.Params("cid")
	seg := ctx.Params("seg")

	// get clientId from the request header
	clientID := ctx.Get("X-Client-ID", "default")

	filename := fmt.Sprintf("chunk%s.m4s", seg)
	cacheKey := fmt.Sprintf("%s/%s", cid, filename)

	// check if the segment is cached
	var cached bool
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

	a.ABRRewriter.Estimator.RecordDownload(clientID, len(segment), duration, cached)

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

func (a *API) streamInit(ctx *fiber.Ctx) error {
	cid := ctx.Params("cid")
	//TODO handle multiple bitrate + use dynamic init files per quality level (Multi-client evaluation / prefetching)
	filename := "init.mp4"

	cacheKey := fmt.Sprintf("%s/%s", cid, filename)
	if _, err := a.Cache.Retrieve(cacheKey); err == nil {
		a.Metrics.CacheHits.Inc()
	} else {
		a.Metrics.CacheMisses.Inc()
	}

	start := time.Now()
	data, err := a.IPFS.Get(fmt.Sprintf("%s/%s", cid, filename))
	if err != nil {
		a.Metrics.ErrorCount.WithLabelValues("GET", "init").Inc()
		return ctx.Status(fiber.StatusBadGateway).SendString("init fetch failed")
	}
	duration := time.Since(start)

	clientID := ctx.Get("X-Client-ID", "default")
	_, cached := a.Cache.Retrieve(cacheKey)
	a.ABRRewriter.Estimator.RecordDownload(clientID, len(data), duration, cached == nil)

	a.Metrics.BytesTransferred.WithLabelValues("GET", "init").Add(float64(len(data)))
	ctx.Set("Content-Type", "video/mp4")
	return ctx.Send(data)
}
