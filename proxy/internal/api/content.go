package api

import (
	"fmt"
	"time"

	"github.com/amirhnajafiz/telescope/pkg/estimator"
	"github.com/amirhnajafiz/telescope/pkg/parser"

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
	mpd, _, err := a.IPFS.Get(fmt.Sprintf("%s/stream.mpd", cid))
	if err != nil {
		a.Logr.Error("failed to fetch mpd", zap.String("cid", cid), zap.Error(err))

		a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), "/api/content").Inc()

		return ctx.Status(fiber.StatusBadGateway).SendString("failed to fetch .mpd")
	}

	// get the header map from the request
	headerMap := ctx.Locals("headerMap").(map[string]float64)

	// rewrite MPD via ABR policy
	rewritten, err := a.ABRRewriter.RewriteMPD(
		mpd,
		clientID,
		cid,
		headerMap["cached"],
		headerMap["uncached"],
		headerMap["current"],
	)
	if err != nil {
		a.Logr.Error("failed to rewrite mpd", zap.String("cid", cid), zap.Error(err))

		a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), "/api/content").Inc()

		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("failed to rewrite manifest:\n %s", err))
	}

	a.Metrics.BytesTransferred.WithLabelValues(ctx.Method(), "/api/content").Add(float64(len(rewritten)))

	ctx.Set("Content-Type", "application/dash+xml")
	ctx.Set("X-Server-BW", parser.ParseMapToHeader(headerMap))

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
		a.Logr.Warn("cache miss", zap.String("cid", cid), zap.String("filename", filename), zap.Error(err))
		a.Metrics.CacheMisses.Inc()
	} else {
		cached = true
		a.Logr.Info("cache hit", zap.String("cid", cid), zap.String("filename", filename))
		a.Metrics.CacheHits.Inc()
	}

	// set the span attributes
	span.SetAttributes(
		attribute.String("cid", cid),
		attribute.String("filename", filename),
		attribute.String("clientId", clientId),
		attribute.Bool("cached", cached),
	)

	// calculate cache ratio
	total := float64(a.Cache.GetHitCounts() + a.Cache.GetMissCounts())
	if total > 0 {
		ratio := float64(a.Cache.GetHitCounts()) / total
		a.Metrics.CacheRatio.Set(ratio)
	}

	// fetch the segment from IPFS if not cached
	start := time.Now()
	if !cached {
		segment, _, err = a.IPFS.Get(fmt.Sprintf("%s/%s", cid, filename))
		if err != nil {
			a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), "/api/stream").Inc()
			return ctx.Status(fiber.StatusBadGateway).SendString("fetch failed")
		}
	}
	duration := time.Since(start)

	// get the header map from the request
	headerMap := ctx.Locals("headerMap").(map[string]float64)

	// record the download in the ABR rewriter
	xsb, xsc, xsu := estimator.Estimate(
		len(segment),
		duration,
		cached,
		headerMap["cached"],
		headerMap["uncached"],
		headerMap["current"],
	)

	// update the header map
	headerMap["cached"] = xsc
	headerMap["uncached"] = xsu
	headerMap["current"] = xsb

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

		a.Metrics.LocalStorageSize.Add(float64(len(segment)))
		a.Metrics.Bandwidth.WithLabelValues(ctx.Method(), "/api/stream").Add(float64(len(segment) / int(duration.Microseconds())))
		a.Metrics.RoundTripTime.WithLabelValues(ctx.Method(), "/api/stream").Observe(float64(duration.Microseconds()))
	}

	a.Metrics.BytesTransferred.WithLabelValues(ctx.Method(), "/api/stream").Add(float64(len(segment)))

	ctx.Set("Content-Type", "video/mp4")
	ctx.Set("X-Server-BW", parser.ParseMapToHeader(headerMap))

	return ctx.Send(segment)
}
