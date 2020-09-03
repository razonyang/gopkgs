package models

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Timezone struct {
	Name string `db:"Name"`
}

func FindAllTimezones(ctx context.Context, db *sqlx.DB, dest interface{}) error {
	query := "SELECT Name FROM mysql.time_zone_name ORDER BY Name ASC"
	return db.SelectContext(ctx, dest, query)
}

func FindTimezone(ctx context.Context, db *sqlx.DB, dest interface{}, name string) error {
	query := "SELECT Name FROM mysql.time_zone_name WHERE Name = ?"
	return db.GetContext(ctx, dest, query, name)
}
