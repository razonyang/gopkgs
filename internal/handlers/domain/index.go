package domain

import (
	"fmt"
	"net/http"

	"clevergo.tech/authmiddleware"
	"clevergo.tech/clevergo"
	"clevergo.tech/pagination"
	"github.com/Masterminds/squirrel"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

type QueryParams struct {
	Name string `json:"name"`
}

type Domain struct {
	models.Domain
	PackageCount int64 `json:"package_count" db:"package_count"`
}

func (d Domain) CanDelete() bool {
	return d.PackageCount == 0
}

func (h *Handler) index(c *clevergo.Context) error {
	ctx := c.Context()
	userID := authmiddleware.GetIdentity(ctx).GetID()

	var queryParams QueryParams
	if err := web.DecodeQueryParams(c, &queryParams); err != nil {
		return err
	}

	p := pagination.NewFromContext(c)
	query := squirrel.Select().From("domains").Where(squirrel.Eq{
		"user_id": userID,
	})
	if queryParams.Name != "" {
		query = query.Where(squirrel.Like{
			"name": fmt.Sprintf("%%%s%%", queryParams.Name),
		})
	}

	countQuery, countArgs, err := query.Columns("COUNT(1)").ToSql()
	if err != nil {
		return err
	}
	if err = h.DB.GetContext(ctx, &p.Total, countQuery, countArgs...); err != nil {
		return err
	}

	selectQuery, selectArgs, err := query.Columns("domains.*, COUNT(packages.id) as package_count").
		LeftJoin("packages ON packages.domain_id = domains.id").
		OrderBy("domains.name ASC").
		GroupBy("domains.id").
		Offset(p.UnsignedOffset()).
		Limit(p.UnsignedLimit()).
		ToSql()
	if err != nil {
		return err
	}
	var domains []Domain
	if err = h.DB.SelectContext(ctx, &domains, selectQuery, selectArgs...); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "domain/index.tmpl", clevergo.Map{
		"page":        web.NewPage("Domains"),
		"domains":     domains,
		"queryParams": queryParams,
		"pagination":  p,
	})
}
