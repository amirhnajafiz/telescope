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

// serveFile handles the streaming of content over DASH
func (a *API) serveFile(
	ctx *fiber.Ctx,
	cid,
	filename,
	cacheKey,
	clientId string,
) error {
	// start the tracing span
	_, span := a.Tracer.Start(ctx.Context(), ctx.Path(), trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

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
		segment, err = a.IPFS.Get(fmt.Sprintf("%s/%s", cid, filename))
		if err != nil {
			a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), ctx.Path()).Inc()
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
		a.Metrics.Bandwidth.WithLabelValues(ctx.Method(), ctx.Path()).Add(float64(len(segment) / int(duration.Seconds())))
		a.Metrics.RoundTripTime.WithLabelValues(ctx.Method(), ctx.Path()).Observe(duration.Seconds())
	}

	a.Metrics.BytesTransferred.WithLabelValues(ctx.Method(), ctx.Path()).Add(float64(len(segment)))

	ctx.Set("Content-Type", "video/mp4")
	ctx.Set("X-Server-BW", parser.ParseMapToHeader(headerMap))

	return ctx.Send(segment)
}
