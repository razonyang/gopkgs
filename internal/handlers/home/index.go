package home

import (
	"net/http"
	"time"

	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/models"
)

func (h *Handler) index(c *clevergo.Context) error {
	var domains, packages, downloads int64
	ctx := c.Context()
	if err := models.CountDomains(ctx, h.DB, &domains); err != nil {
		return err
	}
	if err := models.CountPackages(ctx, h.DB, &packages); err != nil {
		return err
	}
	if err := models.CountActionsByKindAndDate(ctx, h.DB, &downloads, models.ActionGoGet, time.Now().AddDate(0, 0, -29)); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "home/index.tmpl", clevergo.Map{
		"domains":   domains,
		"packages":  packages,
		"downloads": downloads,
	})
}
