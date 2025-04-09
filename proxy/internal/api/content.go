package api

import "github.com/gofiber/fiber/v2"

// newContent handles the upload of content, used for uploading
func (a *API) newContent(ctx *fiber.Ctx) error {
	return nil
}

// getContent handles the download of content, used for getting mpd files
func (a *API) getContent(ctx *fiber.Ctx) error {
	cid := ctx.Params("cid")
	clientID := ctx.Get("X-Client-ID", "default")

	// Fetch original MPD file from IPFS
	originalMPD, err := a.IPFS.FetchMPD(cid)
	if err != nil {
		a.Metrics.ErrorCount.WithLabelValues("GET", "manifest").Inc()
		return ctx.Status(fiber.StatusBadGateway).SendString("failed to fetch .mpd")
	}

	// Rewrite MPD via ABR policy
	rewritten, err := a.ABR.RewriteMPD(originalMPD, clientID)
	if err != nil {
		a.Metrics.ErrorCount.WithLabelValues("GET", "rewrite").Inc()
		return ctx.Status(fiber.StatusInternalServerError).SendString("failed to rewrite manifest")
	}

	a.Metrics.BytesTransferred.WithLabelValues("GET", "manifest").Add(float64(len(rewritten)))

	ctx.Type("application/dash+xml")
	return ctx.Send(rewritten)
}

// listContents handles the listing of all contents
func (a *API) listContents(ctx *fiber.Ctx) error {
	return nil
}

// streamContent handles the streaming of content over DASH
func (a *API) streamContent(ctx *fiber.Ctx) error {
	cid := ctx.Params("cid")

	segment, err := a.IPFS.FetchSegment(cid)
	if err != nil {
		a.Metrics.ErrorCount.WithLabelValues("GET", "stream").Inc()
		return ctx.Status(fiber.StatusBadGateway).SendString("fetch failed")
	}

	a.Metrics.BytesTransferred.WithLabelValues("GET", "stream").Add(float64(len(segment)))
	return ctx.Send(segment)
}
