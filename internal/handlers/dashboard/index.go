package dashboard

import (
	"context"
	"net/http"

	"clevergo.tech/authmiddleware"
	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/core"
	"github.com/razonyang/gopkgs/internal/models"
	"github.com/razonyang/gopkgs/internal/web"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/dashboard", h.index)
}

func (h *Handler) index(c *clevergo.Context) error {
	ctx := c.Context()
	userID := authmiddleware.GetIdentity(ctx).GetID()
	domainCard, err := h.getDomainCard(ctx, userID)
	if err != nil {
		return err
	}
	packageCard, err := h.getPackageCard(ctx, userID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dashboard/index.tmpl", clevergo.Map{
		"page": web.NewPage("Dashboard"),
		"cards": []Card{
			domainCard,
			packageCard,
		},
	})
}

func (h *Handler) getDomainCard(ctx context.Context, userID string) (card Card, err error) {
	card = NewCard("DOMAINS", "globe", "primary", "/domain")
	err = models.CountDomainsByUser(ctx, h.DB, &card.Count, userID)
	return
}

func (h *Handler) getPackageCard(ctx context.Context, userID string) (card Card, err error) {
	card = NewCard("PACKAGES", "cubes", "success", "/package")
	err = models.CountPackagesByUser(ctx, h.DB, &card.Count, userID)
	return
}

type Card struct {
	Background string
	Icon       string
	Title      string
	Count      int64
	Link       string
}

func NewCard(title, icon, background, link string) Card {
	return Card{
		Background: background,
		Icon:       icon,
		Title:      title,
		Link:       link,
	}
}
