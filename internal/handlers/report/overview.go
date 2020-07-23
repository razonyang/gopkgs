package report

import (
	"net/http"
	"time"

	"clevergo.tech/clevergo"
	"clevergo.tech/jsend"
	"github.com/Masterminds/squirrel"
	"github.com/razonyang/gopkgs/internal/models"
	"github.com/razonyang/gopkgs/internal/web"
)

func (h *Handler) overview(c *clevergo.Context) error {
	ctx := c.Context()
	now := time.Now()
	query := squirrel.Select().
		Column(squirrel.Alias(squirrel.Expr("IFNULL(SUM(IF(DATE(actions.created_at)=?, 1, 0)), 0)", now.Format("2006-01-02")), "today")).
		Column(squirrel.Alias(squirrel.Expr("IFNULL(SUM(IF(DATE(actions.created_at)=?, 1, 0)), 0)", now.AddDate(0, 0, -1).Format("2006-01-02")), "yesterday")).
		Column(squirrel.Alias(squirrel.Expr("IFNULL(SUM(IF(DATE(actions.created_at)=?, 1, 0)), 0)", now.AddDate(0, 0, -6).Format("2006-01-02")), "last_seven_days")).
		Column(squirrel.Alias(squirrel.Expr("COUNT(1)"), "last_thirty_days")).
		From("actions").
		LeftJoin("packages ON packages.id = actions.package_id").
		LeftJoin("domains ON domains.id = packages.domain_id").
		Where(squirrel.Eq{
			"actions.kind":    models.ActionGoGet,
			"domains.user_id": h.UserID(ctx),
		}).
		Where(squirrel.GtOrEq{
			"actions.created_at": now.AddDate(0, 0, -29).Format("2006-01-02"),
		})

	var queryParams QueryParams
	if err := web.DecodeQueryParams(c, &queryParams); err != nil {
		return err
	}
	if queryParams.DomainID != 0 {
		query = query.Where(squirrel.Eq{
			"packages.domain_id": queryParams.DomainID,
		})
	}
	if queryParams.PackageID != 0 {
		query = query.Where(squirrel.Eq{
			"actions.package_id": queryParams.PackageID,
		})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	var overview Overview
	if err := h.DB.GetContext(ctx, &overview, sql, args...); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jsend.New(overview))
}

type Overview struct {
	Today          int64 `db:"today" json:"today"`
	Yesterday      int64 `db:"yesterday" json:"yesterday"`
	LastSevenDays  int64 `db:"last_seven_days" json:"last_seven_days"`
	LastThirtyDays int64 `db:"last_thirty_days" json:"last_thirty_days"`
}
