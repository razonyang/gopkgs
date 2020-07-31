package pkg

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
	DomainID int64  `schema:"domain_id"`
	Path     string `schema:"path"`
	VCS      string `schema:"vcs"`
}

func (h *Handler) index(c *clevergo.Context) error {
	p := pagination.NewFromContext(c)

	ctx := c.Context()
	userID := authmiddleware.GetIdentity(ctx).GetID()

	var queryParams QueryParams
	if err := web.DecodeQueryParams(c, &queryParams); err != nil {
		return err
	}

	query := squirrel.Select().From("packages").LeftJoin("domains ON domains.id = packages.domain_id").
		Where(squirrel.Eq{
			"domains.user_id": userID,
		})
	if queryParams.DomainID != 0 {
		query = query.Where(squirrel.Eq{
			"packages.domain_id": queryParams.DomainID,
		})
	}
	if queryParams.Path != "" {
		query = query.Where(squirrel.Like{
			"packages.path": fmt.Sprintf("%%%s%%", queryParams.Path),
		})
	}
	if queryParams.VCS != "" {
		query = query.Where(squirrel.Eq{
			"packages.vcs": queryParams.VCS,
		})
	}

	countQuery, countArgs, err := query.Columns("COUNT(DISTINCT(packages.id))").ToSql()
	if err != nil {
		return err
	}
	if err = h.DB.GetContext(ctx, &p.Total, countQuery, countArgs...); err != nil {
		return err
	}

	var packages []models.Package
	selectQuery, selectArgs, err := query.Columns("packages.*", `domains.id as "domain.id"`, `domains.name as "domain.name"`).
		Limit(p.UnsignedLimit()).
		Offset(p.UnsignedOffset()).
		OrderBy("packages.id DESC").
		ToSql()
	if err != nil {
		return err
	}
	if err = h.DB.SelectContext(ctx, &packages, selectQuery, selectArgs...); err != nil {
		return err
	}

	var domains []models.Domain
	if err = models.FindDomainsByUser(ctx, h.DB, &domains, userID); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "package/index.tmpl", clevergo.Map{
		"page":        web.NewPage("Packages"),
		"domains":     domains,
		"packages":    packages,
		"pagination":  p,
		"queryParams": queryParams,
		"vcs":         models.VCSSet,
	})
}
