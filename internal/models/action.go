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
