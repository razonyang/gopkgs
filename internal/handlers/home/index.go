package home

import (
	"context"
	"net/http"
	"time"

	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/helper"
	"pkg.razonyang.com/gopkgs/internal/models"
)

func (h *Handler) index(c *clevergo.Context) error {
	var domains, packages int64
	ctx := c.Context()
	if err := models.CountDomains(ctx, h.DB, &domains); err != nil {
		return err
	}
	if err := models.CountPackages(ctx, h.DB, &packages); err != nil {
		return err
	}
	downloads, err := h.getDownloads(ctx)
	if err != nil {
		return err
	}
	totalDownloads, err := h.getTotalDownloads(ctx)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "home/index.tmpl", clevergo.Map{
		"domains":        domains,
		"packages":       packages,
		"downloads":      downloads,
		"totalDownloads": totalDownloads,
	})
}

func (h *Handler) getTotalDownloads(ctx context.Context) (count int64, err error) {
	v, found := h.Cache.Get("index:downloads:total")
	if found {
		var ok bool
		count, ok = v.(int64)
		if ok {
			return
		}
	}
	err = models.CountActionsByKind(ctx, h.DB, &count, models.ActionGoGet)
	if err != nil {
		return
	}

	h.Cache.SetWithTTL("index:downloads:total", count, 0, 5*time.Minute)

	return
}

func (h *Handler) getDownloads(ctx context.Context) (count int64, err error) {
	v, found := h.Cache.Get("index:downloads")
	if found {
		var ok bool
		count, ok = v.(int64)
		if ok {
			return
		}
	}
	err = models.CountActionsByKindAndDate(ctx, h.DB, &count, models.ActionGoGet, helper.CurrentUTC().AddDate(0, 0, -29))
	if err != nil {
		return
	}

	h.Cache.SetWithTTL("index:downloads", count, 0, 5*time.Minute)

	return
}
