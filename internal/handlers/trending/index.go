package trending

import (
	"net/http"
	"time"

	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

type Packages struct {
	models.Package
	Downloads int64 `json:"downloads" db:"downloads"`
}

func (h *Handler) index(c *clevergo.Context) error {
	limit := 20
	packages := make([]Packages, 0, limit)
	fromDate := time.Now()

	interval := c.DefaultQuery("interval", "month")
	switch interval {
	case "day":
	case "week":
		fromDate = fromDate.AddDate(0, 0, -6)
	default:
		fromDate = fromDate.AddDate(0, 0, -29)
	}

	query := `
SELECT
	packages.*,
	domains.id as "domain.id",
	domains.name as "domain.name",
	actions.downloads
FROM packages
LEFT JOIN domains ON domains.id = packages.domain_id
LEFT JOIN (
	SELECT package_id, COUNT(1) AS downloads
	FROM actions
	WHERE created_at >= ?
	GROUP BY package_id
) actions ON actions.package_id = packages.id
WHERE packages.private = 0
	AND actions.package_id IS NOT NULL
ORDER BY actions.downloads DESC, domains.name ASC, packages.path ASC
LIMIT ?
`
	ctx := c.Context()
	if err := h.DB.SelectContext(ctx, &packages, query, fromDate.Format("2006-01-02"), limit); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "trending/index.tmpl", clevergo.Map{
		"page":     web.NewPage("Trending"),
		"interval": interval,
		"intervals": []Interval{
			{"Today", "day"},
			{"Last 7 days", "week"},
			{"Last 30 days", "month"},
		},
		"packages": packages,
	})
}

type Interval struct {
	Label string
	Value string
}
