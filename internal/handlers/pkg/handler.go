package pkg

import (
	"fmt"

	"clevergo.tech/authmiddleware"
	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/core"
	"github.com/razonyang/gopkgs/internal/middleware"
	"github.com/razonyang/gopkgs/internal/models"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/package", h.index)
	router.Get("/package/create", h.create)
	router.Post("/package/create", h.create)
	router.Get("/package/edit/:id", h.update)
	router.Post("/package/edit/:id", h.update)
	router.Delete("/package/:id", middleware.APIErrorHandler(h.delete))
}

func (h *Handler) findPackage(c *clevergo.Context) (*models.Package, error) {
	id, err := c.Params.Int64("id")
	if err != nil {
		return nil, fmt.Errorf("invalid package ID: %s", err.Error())
	}
	ctx := c.Context()
	userID := authmiddleware.GetIdentity(ctx).GetID()
	var pkg models.Package
	if err = models.FindPackageByUser(ctx, h.DB, &pkg, id, userID); err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (h *Handler) findDomains(c *clevergo.Context) (domains []models.Domain, err error) {
	ctx := c.Context()
	query := "SELECT * FROM domains WHERE user_id = ? AND verified = 1 ORDER BY domains.name ASC"
	err = h.DB.SelectContext(ctx, &domains, query, h.UserID(ctx))
	return
}
