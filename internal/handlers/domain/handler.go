package domain

import (
	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/core"
	"pkg.razonyang.com/gopkgs/internal/middleware"
	"pkg.razonyang.com/gopkgs/internal/models"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/domain", h.index)
	router.Get("/domain/create", h.create)
	router.Post("/domain/create", h.create)
	router.Get("/domain/edit/:id", h.update)
	router.Post("/domain/edit/:id", h.update)
	router.Get("/domain/verify/:id", h.verify)
	router.Post("/domain/verify/:id", h.verify)
	router.Get("/domain/challenge-txt/:id", h.challengeTXT)
	router.Get("/.well-known/gopkgs-challenge/:token", h.challenge)
	router.Delete("/domain/:id", middleware.APIErrorHandler(h.delete))
}

func (h *Handler) findDomain(c *clevergo.Context) (*models.Domain, error) {
	ctx := c.Context()
	id, err := c.Params.Int64("id")
	if err != nil {
		return nil, err
	}
	userID := h.UserID(ctx)
	var domain models.Domain
	err = models.FindDomainByUser(ctx, h.DB, &domain, id, userID)
	return &domain, err
}
