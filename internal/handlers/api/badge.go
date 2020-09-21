package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"clevergo.tech/clevergo"
	"clevergo.tech/shields"
	"pkg.razonyang.com/gopkgs/internal/helper"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/stringhelper"
)

func (h *Handler) download(c *clevergo.Context) error {
	interval := c.Params.String("interval")

	path := strings.Split(strings.TrimPrefix(c.Params.String("path"), "/"), "/")
	if len(path) < 2 {
		return h.notFound(c)
	}
	ctx := c.Context()
	var pkg models.Package
	err := models.FindPackageByDomainAndPath(ctx, h.DB, &pkg, path[0], strings.Join(path[1:], "/"))
	if err != nil {
		if err == sql.ErrNoRows {
			return h.notFound(c)
		}
		return err
	}

	count, err := h.getDownloads(ctx, interval, pkg.ID)
	if err != nil {
		return err
	}

	var downloads string
	if interval == "total" {
		downloads = stringhelper.ShortScale(count)
	} else {
		downloads = fmt.Sprintf("%s/%s", stringhelper.ShortScale(count), interval)
	}

	badge := shields.New("downloads", downloads)
	badge.Color = shields.ColorBrightGreen
	if err := badge.ParseRequest(c.Request); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, badge)
}

func (h *Handler) notFound(c *clevergo.Context) error {
	badge := shields.New("downloads", "package not found")
	badge.Color = shields.ColorRed
	return c.JSON(http.StatusOK, badge)
}

func (h *Handler) getDownloads(ctx context.Context, interval string, id int64) (int64, error) {
	fromDate := helper.CurrentUTC()
	switch interval {
	case "day":
	case "week":
		fromDate = fromDate.AddDate(0, 0, -6)
	case "month":
		fromDate = fromDate.AddDate(0, 0, -29)
	case "total":
		fromDate = time.Time{}
	default:
		return 0, fmt.Errorf("invalid interval parameter")
	}

	key := fmt.Sprintf("badge:downloads:%s:%d", interval, id)
	v, found := h.Cache.Get(key)
	if found {
		if count, ok := v.(int64); ok {
			return count, nil
		}
	}

	query := `
SELECT COUNT(1) FROM actions
WHERE package_id = ? 
	AND created_at >= ?
`
	var count int64
	err := h.DB.GetContext(ctx, &count, query, id, fromDate.Format("2006-01-02"))
	if err != nil {
		return 0, err
	}

	h.Cache.SetWithTTL(key, count, 0, 5*time.Minute)

	return count, err
}
