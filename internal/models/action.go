package models

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	ActionGoGet = "go-get"
)

type Action struct {
	ID        string    `db:"id" json:"id"`
	Kind      string    `db:"kind" json:"kind"`
	PackageID int64     `db:"package_id" json:"package_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewAction(kind string, packageID int64) *Action {
	return &Action{
		ID:        strings.ReplaceAll(uuid.New().String(), "-", ""),
		Kind:      kind,
		PackageID: packageID,
		CreatedAt: time.Now(),
	}
}

func (act *Action) Save(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx,
		"INSERT INTO actions(id, kind, package_id, created_at) VALUES(?, ?, ?, ?)",
		act.ID, act.Kind, act.PackageID, act.CreatedAt,
	)
	return err
}

func CountActionsByKind(ctx context.Context, db *sqlx.DB, count *int64, kind string) error {
	query := "SELECT COUNT(1) FROM actions WHERE kind = ?"
	return db.GetContext(ctx, count, query, kind)
}

func CountActionsByKindAndDate(ctx context.Context, db *sqlx.DB, count *int64, kind string, fromDate time.Time) error {
	query := "SELECT COUNT(1) FROM actions WHERE kind = ? AND created_at >= ?"
	return db.GetContext(ctx, count, query, kind, fromDate.Format("2006-01-02"))
}
