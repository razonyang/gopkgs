package validators

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Domain struct {
	db *sqlx.DB
}

func NewDomain(db *sqlx.DB) *Domain {
	return &Domain{
		db: db,
	}
}

// ValidateName validates the domain name is valid.
func (d *Domain) ValidateName(ctx context.Context, value interface{}) error {
	name := value.(string)
	query := "SELECT COUNT(*) FROM domains WHERE name = ? AND verified = 1"
	var count int64
	err := d.db.GetContext(ctx, &count, query, name)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("domain %q has been taken", name)
	}
	return nil
}
