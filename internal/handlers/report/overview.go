package report

import (
	"net/http"

	"clevergo.tech/clevergo"
	"clevergo.tech/jsend"
	"github.com/razonyang/gopkgs/internal/models"
)

func (h *Handler) overview(c *clevergo.Context) error {
	query := `
SELECT
	IFNULL(SUM(IF(DATE(actions.created_at)=CURRENT_DATE(), 1, 0)), 0) as today,
	IFNULL(SUM(IF(DATE(actions.created_at)=CURRENT_DATE() - INTERVAL 1 DAY, 1, 0)), 0) as yesterday,
	IFNULL(SUM(IF(DATE(actions.created_at)>=CURRENT_DATE() - INTERVAL 6 DAY, 1, 0)), 0) as last_seven_days,
	COUNT(1) as last_thirty_days
FROM actions
LEFT JOIN packages ON packages.id = actions.package_id
LEFT JOIN domains ON domains.id = packages.domain_id
WHERE actions.kind = ?
	AND domains.user_id = ?
	AND actions.created_at >= CURRENT_DATE() - INTERVAL 29 DAY
`
	ctx := c.Context()
	var overview Overview
	if err := h.DB.GetContext(ctx, &overview, query, models.ActionGoGet, h.UserID(ctx)); err != nil {
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
