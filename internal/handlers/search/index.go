package search

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"clevergo.tech/clevergo"
	"github.com/Masterminds/squirrel"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

type Package struct {
	models.Package
	Downloads int64 `json:"downloads" db:"downloads"`
}

func (h *Handler) index(c *clevergo.Context) error {
	q := strings.TrimSpace(c.QueryParam("q"))
	if q == "" {
		return c.Render(http.StatusOK, "search/index.tmpl", clevergo.Map{
			"page": web.NewPage("Search"),
		})
	}
	ctx := c.Context()

	packages, totalCount, err := h.getPackages(ctx, q)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "search/result.tmpl", clevergo.Map{
		"page":       web.NewPage(fmt.Sprintf("%s - %s", q, "Search Result")),
		"q":          q,
		"packages":   packages,
		"totalCount": totalCount,
	})
}

func (h *Handler) getPackages(ctx context.Context, name string) (packages []Package, count int64, err error) {
	query := squirrel.Select().From("packages").
		LeftJoin("domains ON domains.id = packages.domain_id").
		Where(squirrel.Eq{
			"packages.private": 0,
		}).
		Where(squirrel.Like{
			"packages.path": fmt.Sprintf("%%%s%%", name),
		})

	countSql, countArgs, err := query.Columns("COUNT(1) as count").ToSql()
	if err != nil {
		return
	}
	if err = h.DB.GetContext(ctx, &count, countSql, countArgs...); err != nil {
		return
	}

	cloneQuery := query.
		Columns(
			"packages.*",
			`domains.id as "domain.id"`,
			`domains.name as "domain.name"`,
		).
		OrderBy("packages.id").
		Offset(0).
		Limit(20)
	sql, args, err := cloneQuery.ToSql()
	if err != nil {
		return
	}
	if err = h.DB.SelectContext(ctx, &packages, sql, args...); err != nil {
		return
	}

	return
}
