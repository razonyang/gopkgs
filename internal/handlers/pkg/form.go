package pkg

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
)

type Form struct {
	db            *sqlx.DB
	pkg           *models.Package
	userID        string
	DomainID      int64  `json:"domain_id" schema:"domain_id"`
	Path          string `json:"path" schema:"path"`
	VCS           string `json:"vcs" schema:"vcs"`
	Root          string `json:"root" schema:"root"`
	Docs          string `json:"docs" schema:"docs"`
	Private       bool   `json:"-" schema:"-"`
	PrivateSwitch string `json:"private" schema:"private"`
	Description   string `json:"description" schema:"description"`
	Homepage      string `json:"homepage" schema:"homepage"`
	License       string `json:"license" schema:"license"`
}

func newForm(db *sqlx.DB, userID string) *Form {
	return &Form{
		db:     db,
		userID: userID,
	}
}

func newFormPkg(db *sqlx.DB, userID string, pkg *models.Package) *Form {
	f := newForm(db, userID)
	f.pkg = pkg
	f.DomainID = f.pkg.DomainID
	f.Path = f.pkg.Path
	f.Root = f.pkg.Root
	f.VCS = f.pkg.VCS
	f.Docs = f.pkg.Docs
	f.Private = f.pkg.Private
	f.Description = f.pkg.Description
	f.Homepage = f.pkg.Homepage
	f.License = f.pkg.License
	return f
}

// Validate validates package.
func (f *Form) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.DomainID, validation.Required, validation.WithContext(f.validateDomain)),
		validation.Field(&f.Path, validation.Required, validation.WithContext(f.validatePath)),
		validation.Field(&f.VCS, validation.Required, validation.In(models.VCSGit, models.VCSSvn, models.VCSBzr, models.VCSHg, models.VCSFossil)),
		validation.Field(&f.Root, validation.Required),
		validation.Field(&f.Docs, is.URL),
		validation.Field(&f.Homepage, is.URL),
		validation.Field(&f.License, is.URL),
	)
}

func (f *Form) validateDomain(ctx context.Context, value interface{}) error {
	query := "SELECT * FROM domains WHERE id = ? AND user_id = ?"
	var domain models.Domain
	if err := f.db.GetContext(ctx, &domain, query, value.(int64), f.userID); err != nil {
		return err
	}
	if !domain.Verified {
		return errors.New("domain has not been verified")
	}
	return nil
}

func (f *Form) validatePath(ctx context.Context, value interface{}) error {
	query := "SELECT COUNT(1) FROM packages WHERE domain_id = ? AND path = ?"
	args := []interface{}{f.DomainID, value.(string)}
	if f.pkg != nil {
		query += " AND id != ?"
		args = append(args, f.pkg.ID)
	}
	var count int64
	if err := f.db.GetContext(ctx, &count, query, args...); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("package already exists")
	}
	return nil
}

func (f *Form) Create(ctx context.Context) (pkg *models.Package, err error) {
	pkg = models.NewPackage(f.DomainID, f.Path, f.VCS, f.Root)
	pkg.Docs = f.Docs
	pkg.Description = f.Description
	pkg.Homepage = f.Homepage
	pkg.License = f.License
	pkg.Private = f.PrivateValue()
	err = pkg.Insert(ctx, f.db)
	return
}

func (f *Form) Update(ctx context.Context) (err error) {
	f.pkg.DomainID = f.DomainID
	f.pkg.Path = f.Path
	f.pkg.Root = f.Root
	f.pkg.VCS = f.VCS
	f.pkg.Docs = f.Docs
	f.pkg.Description = f.Description
	f.pkg.Homepage = f.Homepage
	f.pkg.License = f.License
	f.pkg.Private = f.PrivateValue()
	return f.pkg.Update(ctx, f.db)
}

func (f *Form) PrivateValue() bool {
	if f.PrivateSwitch != "" {
		return true
	}
	return false
}
