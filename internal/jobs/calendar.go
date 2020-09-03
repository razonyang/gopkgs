package jobs

import (
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/helper"
)

type Calendar struct {
	db *sqlx.DB
}

func NewCalendar(db *sqlx.DB) *Calendar {
	return &Calendar{
		db: db,
	}
}

func (c *Calendar) Run() {
	query := "INSERT IGNORE INTO calendars(id) VALUES "
	query = strings.TrimSuffix(query+strings.Repeat("(?),", 30), ",")
	args := make([]interface{}, 30)
	now := helper.CurrentUTC().AddDate(0, 0, -29)
	for i := 0; i < 30; i++ {
		args[i] = now.AddDate(0, 0, i)
	}
	if _, err := c.db.Exec(query, args...); err != nil {
		log.Println(err)
	}
}
