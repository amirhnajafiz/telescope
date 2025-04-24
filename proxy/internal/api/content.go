package api

import (
	"fmt"

	"github.com/amirhnajafiz/telescope/pkg/parser"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// getContent handles the download of content, used for getting mpd files
func (a *API) getContent(ctx *fiber.Ctx) error {
	// start the tracing span
	_, span := a.Tracer.Start(ctx.Context(), ctx.Path(), trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// get the cid from the URL
	cid := ctx.Params("cid")

	// get clientId from the request header
	clientID := ctx.Get("X-Client-ID", "default")

	// set the span attributes
	span.SetAttributes(attribute.String("cid", cid))
	span.SetAttributes(attribute.String("clientId", clientID))

	// fetch MPD file from IPFS
	mpd, err := a.IPFS.Get(fmt.Sprintf("%s/stream.mpd", cid))
	if err != nil {
		a.Logr.Error("failed to fetch mpd", zap.String("cid", cid), zap.Error(err))

		a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), ctx.Path()).Inc()

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

		a.Metrics.ErrorCount.WithLabelValues(ctx.Method(), ctx.Path()).Inc()

		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("failed to rewrite manifest:\n %s", err))
	}

	a.Metrics.BytesTransferred.WithLabelValues(ctx.Method(), ctx.Path()).Add(float64(len(rewritten)))

	ctx.Set("Content-Type", "application/dash+xml")
	ctx.Set("X-Server-BW", parser.ParseMapToHeader(headerMap))

	return ctx.Send(rewritten)
}

// streamContent handles the streaming of content over DASH
func (a *API) streamContent(ctx *fiber.Ctx) error {
	// start the tracing span
	_, span := a.Tracer.Start(ctx.Context(), ctx.Path(), trace.WithSpanKind(trace.SpanKindServer))
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

	return a.serveFile(ctx, cid, filename, cacheKey, clientId)
}
