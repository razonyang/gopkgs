package dashboard

import (
	"context"
	"net/http"
	"time"

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
	dailyReportCard, err := h.getDailyReportCard(ctx, userID)
	if err != nil {
		return err
	}
	monthlyReportCard, err := h.getMonthlyReportCard(ctx, userID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dashboard/index.tmpl", clevergo.Map{
		"page": web.NewPage("Dashboard"),
		"cards": []Card{
			domainCard,
			packageCard,
			dailyReportCard,
			monthlyReportCard,
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

func (h *Handler) getDailyReportCard(ctx context.Context, userID string) (card Card, err error) {
	card = NewCard("Daily Report", "cloud-download-alt", "secondary", "/report")
	return card, h.getReport(ctx, &card.Count, userID, time.Now())
}

func (h *Handler) getMonthlyReportCard(ctx context.Context, userID string) (card Card, err error) {
	card = NewCard("Monthly Report", "cloud-download-alt", "info", "/report")
	return card, h.getReport(ctx, &card.Count, userID, time.Now().AddDate(0, 0, -29))
}

func (h *Handler) getReport(ctx context.Context, count *int64, userID string, fromDate time.Time) error {
	query := `
SELECT COUNT(1) FROM actions 
LEFT JOIN packages ON packages.id = actions.package_id
LEFT JOIN domains ON domains.id = packages.domain_id
WHERE domains.user_id = ?
	AND actions.created_at >= ?
`
	return h.DB.GetContext(ctx, count, query, userID, fromDate.Format("2006-01-02"))
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
