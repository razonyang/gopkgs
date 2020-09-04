package search

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"clevergo.tech/clevergo"
	"github.com/Masterminds/squirrel"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

type Package struct {
	models.Package
	Downloads int64  `json:"downloads" db:"downloads"`
	FullPath  string `json:"full_path" db:"full_path"`
}

func (h *Handler) index(c *clevergo.Context) error {
	q := strings.TrimSpace(c.QueryParam("q"))
	if q == "" {
		return c.Render(http.StatusOK, "search/index.tmpl", clevergo.Map{
			"page": web.NewPage("Search"),
		})
	}

	pageNum, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}
	limit := int64(10)

	ctx := c.Context()

	packages, totalCount, err := h.getPackages(ctx, q, pageNum, limit)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "search/result.tmpl", clevergo.Map{
		"page":       web.NewPage(fmt.Sprintf("%s - %s", q, "Search Result")),
		"pageNum":    pageNum,
		"limit":      limit,
		"q":          q,
		"packages":   packages,
		"totalCount": totalCount,
	})
}

func (h *Handler) getPackages(ctx context.Context, name string, page, limit int64) (packages []Package, count int64, err error) {
	query := squirrel.Select().From("packages").
		Columns(
			"packages.*",
			`domains.id as "domain.id"`,
			`domains.name as "domain.name"`,
			`CONCAT(domains.name, "/", packages.path) as  full_path`,
			`IFNULL(actions.downloads, 0) as downloads`,
		).
		LeftJoin("domains ON domains.id = packages.domain_id").
		LeftJoin("(SELECT package_id, COUNT(1) as downloads FROM actions GROUP BY package_id) actions ON actions.package_id = packages.id").
		Where(squirrel.Eq{
			"packages.private": 0,
		}).
		Having(squirrel.Like{
			"full_path": fmt.Sprintf("%%%s%%", name),
		})

	countSql, countArgs, err := squirrel.Select("COUNT(1) as count").FromSelect(query, "t").ToSql()
	if err != nil {
		return
	}
	if err = h.DB.GetContext(ctx, &count, countSql, countArgs...); err != nil {
		return
	}

	cloneQuery := query.
		OrderBy("actions.downloads DESC", "packages.path ASC", "packages.id ASC").
		Offset(uint64((page - 1) * limit)).
		Limit(uint64(limit))
	sql, args, err := cloneQuery.ToSql()
	if err != nil {
		return
	}
	if err = h.DB.SelectContext(ctx, &packages, sql, args...); err != nil {
		return
	}

	return
}
