package domain

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/razonyang/gopkgs/internal/models"
)

type Form struct {
	domain *models.Domain
	db     *sqlx.DB
	userID string
	Name   string `json:"name" schema:"name"`
}

func NewForm(db *sqlx.DB, userID string) *Form {
	return &Form{
		db:     db,
		userID: userID,
	}
}

func (f *Form) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Name, validation.Required, is.Domain,
			validation.WithContext(f.isDomainExists), validation.WithContext(f.isDomainTaken),
			validation.WithContext(f.isLimitExceeded),
		),
	)
}

func (f *Form) isLimitExceeded(ctx context.Context, value interface{}) error {
	if f.domain != nil {
		return nil
	}

	var count int64
	query := "SELECT COUNT(1) FROM domains WHERE user_id = ?"
	err := f.db.GetContext(ctx, &count, query, f.userID)
	if err != nil {
		return err
	}
	max := int64(5)
	if count >= max {
		return fmt.Errorf("the number of domain names exceeds the limit(maximum: %d)", max)
	}
	return nil
}

func (f *Form) isDomainExists(ctx context.Context, value interface{}) error {
	name := value.(string)
	var count int64
	query := "SELECT COUNT(1) FROM domains WHERE name = ? AND user_id = ?"
	args := []interface{}{name, f.userID}
	if f.domain != nil {
		query += " AND id != ?"
		args = append(args, f.domain.ID)
	}
	err := f.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("domain %s already exists", name)
	}
	return nil
}

func (f *Form) isDomainTaken(ctx context.Context, value interface{}) error {
	name := value.(string)
	var count int64
	query := "SELECT COUNT(*) FROM domains WHERE name = ? AND verified = 1"
	args := []interface{}{name}
	if f.domain != nil {
		query += " AND id != ?"
		args = append(args, f.domain.ID)
	}
	err := f.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("domain %q has been taken", name)
	}
	return nil
}

func (f *Form) Create(ctx context.Context) (*models.Domain, error) {
	domain := models.NewDomain(f.Name, f.userID)
	err := domain.Insert(ctx, f.db)
	return domain, err
}

func (f *Form) Update(ctx context.Context) error {
	if f.domain.Name != f.Name {
		f.domain.Verified = false
		f.domain.ChallengeTXT = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	f.domain.Name = f.Name
	err := f.domain.Update(ctx, f.db)
	return err
}
