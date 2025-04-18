package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// getContent handles the download of content, used for getting mpd files
func (a *API) getContent(ctx *fiber.Ctx) error {
	cid := ctx.Params("cid")
	clientID := ctx.Get("X-Client-ID", "default")

	// Fetch original MPD file from IPFS
	originalMPD, err := a.IPFS.Get(cid)
	if err != nil {
		a.Metrics.ErrorCount.WithLabelValues("GET", "manifest").Inc()
		return ctx.Status(fiber.StatusBadGateway).SendString("failed to fetch .mpd")
	}

	// Rewrite MPD via ABR policy
	rewritten, err := a.ABRRewriter.RewriteMPD(originalMPD, clientID, cid)
	if err != nil {
		a.Metrics.ErrorCount.WithLabelValues("GET", "rewrite").Inc()
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("failed to rewrite manifest:\n %s", err))
	}

	a.Metrics.BytesTransferred.WithLabelValues("GET", "manifest").Add(float64(len(rewritten)))

	ctx.Set("Content-Type", "application/dash+xml")
	return ctx.Send(rewritten)
}

// listContents handles the listing of all contents
func (a *API) listContents(ctx *fiber.Ctx) error {
	return nil
}

// streamContent handles the streaming of content over DASH
func (a *API) streamContent(ctx *fiber.Ctx) error {
	cid := ctx.Params("cid")
	seg := ctx.Params("seg")

	filename := fmt.Sprintf("chunk%s.m4s", seg)

	// Mark metrics
	cacheKey := fmt.Sprintf("%s/%s", cid, filename)
	if _, err := a.Cache.Retrieve(cacheKey); err != nil {
		a.Metrics.CacheHits.Inc()
	} else {
		a.Metrics.CacheMisses.Inc()
	}

	// calculate cache ratio
	total := float64(a.Cache.GetHitCounts() + a.Cache.GetMissCounts())
	if total > 0 {
		ratio := float64(a.Cache.GetHitCounts()) / total
		a.Metrics.CacheRatio.Set(ratio)
	}

	// Fetch from IPFS
	start := time.Now()
	segment, err := a.IPFS.Get(fmt.Sprintf("%s/%s", cid, filename))
	if err != nil {
		a.Metrics.ErrorCount.WithLabelValues("GET", "stream").Inc()
		return ctx.Status(fiber.StatusBadGateway).SendString("fetch failed")
	}
	duration := time.Since(start)

	clientID := ctx.Get("X-Client-ID", "default")
	_, cached := a.Cache.Retrieve(cacheKey)
	a.ABRRewriter.Estimator.RecordDownload(clientID, len(segment), duration, cached == nil)

	a.Metrics.BytesTransferred.WithLabelValues("GET", "stream").Add(float64(len(segment)))
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
