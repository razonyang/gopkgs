package tasks

import (
	"context"

	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
)

type Package struct {
	db *sqlx.DB
}

func NewPackage(db *sqlx.DB) *Package {
	return &Package{db: db}
}

func (pkg *Package) AddAction(kind string, packageID int64) error {
	action := models.NewAction(kind, packageID)
	return action.Save(context.Background(), pkg.db)
}
