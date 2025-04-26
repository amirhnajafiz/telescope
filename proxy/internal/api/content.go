package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// getContent handles the download of content, used for getting mpd files
func (a *API) getContent(ctx *fiber.Ctx) error {
	// start the tracing span
	_, span := a.Tracer.Start(ctx.Context(), "/api/content", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// get the cid from the URL
	cid := ctx.Params("cid")

	// get clientId from the request header
	clientID := ctx.Get("X-Client-ID", "default")

	// set the span attributes
	span.SetAttributes(attribute.String("cid", cid))
	span.SetAttributes(attribute.String("clientId", clientID))

	// fetch MPD file from IPFS
	mpd, rtt, err := a.IPFS.Get(fmt.Sprintf("%s/stream.mpd", cid))
	if err != nil {
		a.Logr.Error("failed to fetch mpd", zap.String("cid", cid), zap.Error(err))
		a.Metrics.SysErrorCount.WithLabelValues("/api/content").Inc()
		return ctx.Status(fiber.StatusBadGateway).SendString("failed to fetch .mpd")
	}

	a.Metrics.IPFSRTT.Observe(float64(rtt))
	a.Metrics.IPFSBandwidth.Set(float64(len(mpd)) / float64(rtt))

	// get client bandwidth from the request header
	clientBandwidth, _ := strconv.ParseFloat(ctx.Get("X-Bandwidth", "0"), 64)

	// build MPD
	modifiedMPD, err := a.MPDBuilder.Build(mpd, cid)
	if err != nil {
		a.Logr.Error("failed to build mpd", zap.String("cid", cid), zap.Error(err))
		a.Metrics.SysErrorCount.WithLabelValues("/api/content").Inc()
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("failed to build manifest:\n %s", err))
	}

	// rewrite MPD via ABR policy
	rewritten, err := a.ABRRewriter.RewriteMPD(modifiedMPD, cid, clientBandwidth)
	if err != nil {
		a.Logr.Error("failed to rewrite mpd", zap.String("cid", cid), zap.Error(err))
		a.Metrics.SysErrorCount.WithLabelValues("/api/content").Inc()
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("failed to rewrite manifest:\n %s", err))
	}

	a.Metrics.SysBytesTransferred.WithLabelValues("/api/content").Add(float64(len(rewritten)))

	ctx.Set("Content-Type", "application/dash+xml")

	return ctx.Send(rewritten)
}

// streamContent handles the streaming of content over DASH
func (a *API) streamContent(ctx *fiber.Ctx) error {
	// start the tracing span
	_, span := a.Tracer.Start(ctx.Context(), "/api/stream", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// get the cid and segment from the URL
	cid := ctx.Params("cid")
	seg := ctx.Params("seg")

	// get clientId from the request header
	clientId := ctx.Get("X-Client-ID", "default")

	// set the span attributes
	span.SetAttributes(attribute.String("cid", cid))
	span.SetAttributes(attribute.String("clientId", clientId))
	span.SetAttributes(attribute.String("seg", seg))

	// build the filename and cache key
	filename := seg
	cacheKey := fmt.Sprintf("%s/%s", cid, filename)

	var cached bool

	// check if the segment is cached
	segment, err := a.Cache.Retrieve(cacheKey)
	if err != nil {
		cached = false
		a.Metrics.SysCacheMisses.Inc()
	} else {
		cached = true
		a.Metrics.SysCacheHits.Inc()
	}

	// set the span attributes
	span.SetAttributes(
		attribute.String("cid", cid),
		attribute.String("clientId", clientId),
	)

	// calculate cache ratio
	total := float64(a.Cache.GetHitCounts() + a.Cache.GetMissCounts())
	if total > 0 {
		ratio := float64(a.Cache.GetHitCounts()) / total
		a.Metrics.SysCacheRatio.Set(ratio)
	}

	// fetch the segment from IPFS if not cached
	start := time.Now()
	if !cached {
		segment, _, err = a.IPFS.Get(fmt.Sprintf("%s/%s", cid, filename))
		if err != nil {
			a.Logr.Error("failed to fetch segment", zap.String("cid", cid), zap.String("filename", filename), zap.Error(err))
			a.Metrics.SysErrorCount.WithLabelValues("/api/stream").Inc()
			return ctx.Status(fiber.StatusBadGateway).SendString("fetch failed")
		}
	}
	duration := time.Since(start)

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

		a.Metrics.IPFSBandwidth.Set(float64(len(segment)) / float64(duration.Microseconds()))
		a.Metrics.IPFSRTT.Observe(float64(duration.Microseconds()))
	}

	a.Metrics.SysBytesTransferred.WithLabelValues("/api/stream").Add(float64(len(segment)))

	ctx.Set("Content-Type", "video/mp4")

	return ctx.Send(segment)
}
