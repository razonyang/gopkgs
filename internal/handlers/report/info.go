package report

import (
	"fmt"
	"net/http"
	"time"

	"clevergo.tech/clevergo"
	"clevergo.tech/jsend"
	"github.com/Masterminds/squirrel"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

type QueryParams struct {
	DomainID  int64 `json:"domain_id" schema:"domain_id"`
	PackageID int64 `json:"package_id" schema:"package_id"`
}

func (h *Handler) info(c *clevergo.Context) error {
	ctx := c.Context()

	var queryParams QueryParams
	if err := web.DecodeQueryParams(c, &queryParams); err != nil {
		return err
	}

	fromDate := time.Now().AddDate(0, 0, -29).Format("2006-01-02")
	actionsQuery := squirrel.Select("DATE(actions.created_at) as date", "COUNT(1) as count").
		From("actions").
		LeftJoin("packages ON packages.id = actions.package_id").
		LeftJoin("domains ON domains.id = packages.domain_id").
		Where(squirrel.Eq{
			"actions.kind":    models.ActionGoGet,
			"domains.user_id": h.UserID(ctx),
		}).
		Where(squirrel.Gt{
			"actions.created_at": fromDate,
		}).
		GroupBy("DATE(actions.created_at)")
	if queryParams.DomainID != 0 {
		actionsQuery = actionsQuery.Where(squirrel.Eq{
			"packages.domain_id": queryParams.DomainID,
		})
	}
	if queryParams.PackageID != 0 {
		actionsQuery = actionsQuery.Where(squirrel.Eq{
			"actions.package_id": queryParams.PackageID,
		})
	}

	actionsSQL, actionsArgs, err := actionsQuery.ToSql()
	if err != nil {
		return err
	}

	query, args, err := squirrel.Select("c.id as date", "IFNULL(a.count, 0) as count").From("calendars c").
		LeftJoin(fmt.Sprintf("(%s) a ON c.id = a.date", actionsSQL), actionsArgs...).
		Where(squirrel.Gt{
			"c.id": fromDate,
		}).
		OrderBy("c.id ASC").
		ToSql()
	if err != nil {
		return err
	}

	var lines []Line
	if err = h.DB.SelectContext(ctx, &lines, query, args...); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jsend.New(lines))
}

type Line struct {
	Date  time.Time `db:"date" json:"date"`
	Count int64     `db:"count" json:"count"`
}
