package tasks

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
)

type Package struct {
	db *sqlx.DB
}

func NewPackage(db *sqlx.DB) *Package {
	return &Package{db: db}
}

func (pkg *Package) Action(kind string, packageID int64, createdAt int64) error {
	action := models.NewAction(kind, packageID)
	action.CreatedAt = time.Unix(createdAt, 0)
	return action.Save(context.Background(), pkg.db)
}
