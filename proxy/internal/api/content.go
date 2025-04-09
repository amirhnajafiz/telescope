package api

import "github.com/gofiber/fiber/v2"

// newContent handles the upload of content, used for uploading
func (a *API) newContent(ctx *fiber.Ctx) error {
	return nil
}

// getContent handles the download of content, used for getting mpd files
func (a *API) getContent(ctx *fiber.Ctx) error {
	return nil
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
