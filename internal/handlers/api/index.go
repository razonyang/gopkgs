package api

import (
	"net/http"

	"clevergo.tech/clevergo"
	"clevergo.tech/jsend"
	"github.com/Masterminds/squirrel"
	"github.com/razonyang/gopkgs/internal/core"
	"github.com/razonyang/gopkgs/internal/models"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/api/domains", h.domains)
	router.Get("/api/packages", h.packages)
}

func (h *Handler) domains(c *clevergo.Context) error {
	query := "SELECT * FROM domains WHERE user_id = ? ORDER BY name ASC"
	var domains []models.Domain
	ctx := c.Context()
	if err := h.DB.SelectContext(ctx, &domains, query, h.UserID(ctx)); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, jsend.New(domains))
}

func (h *Handler) packages(c *clevergo.Context) error {
	ctx := c.Context()
	query := squirrel.Select("packages.*").From("packages").LeftJoin("domains on domains.id = packages.domain_id").
		Where(squirrel.Eq{
			"domains.user_id": h.UserID(ctx),
		})
	if domainID := c.QueryParam("domain_id"); domainID != "" {
		query = query.Where(squirrel.Eq{
			"packages.domain_id": domainID,
		})
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	packages := []models.Package{}
	if err := h.DB.SelectContext(ctx, &packages, sql, args...); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, jsend.New(packages))
}
