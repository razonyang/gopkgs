package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Version control system constants.
const (
	VCSGit    = "git"
	VCSSvn    = "svn"
	VCSBzr    = "bzr"
	VCSHg     = "hg"
	VCSFossil = "fossil"
)

var VCSSet = []string{VCSGit, VCSSvn, VCSBzr, VCSHg, VCSFossil}

// Package is a model that mapping to table "packages".
type Package struct {
	Model
	DomainID    int64  `db:"domain_id" json:"domain_id" schema:"domain_id"`
	Private     bool   `db:"private" json:"private" schema:"private"`
	Path        string `db:"path" json:"path" schema:"path"`
	VCS         string `db:"vcs" json:"vcs" schema:"vcs"`
	Root        string `db:"root" json:"root" schema:"root"`
	Docs        string `db:"docs" json:"docs" schema:"docs"`
	Description string `db:"description" json:"description" schema:"description"`
	Homepage    string `db:"homepage" json:"homepage" schema:"homepage"`
	License     string `db:"license" json:"license" schema:"license"`

	Domain Domain `db:"domain,prefix=domain."`
}

func NewPackage(domainID int64, path, vcs, root string) *Package {
	return &Package{
		DomainID: domainID,
		Path:     path,
		VCS:      vcs,
		Root:     root,
	}
}

// Insert saves package into database.
func (pkg *Package) Insert(ctx context.Context, db *sqlx.DB) error {
	now := time.Now()
	pkg.CreatedAt = now
	pkg.UpdatedAt = now
	query := `
INSERT INTO packages(domain_id, private, path, vcs, root, docs, description, homepage, license, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	res, err := db.ExecContext(ctx, query, pkg.DomainID, pkg.Private, pkg.Path, pkg.VCS, pkg.Root, pkg.Docs, pkg.Description, pkg.Homepage, pkg.License, now, now)
	if err != nil {
		return err
	}
	pkg.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (pkg *Package) Update(ctx context.Context, db *sqlx.DB) error {
	pkg.UpdatedAt = time.Now()
	query := "UPDATE packages SET domain_id = ?, private = ?, path = ?, vcs = ?, root = ?, docs = ?, description = ?, homepage = ?, license = ?, updated_at = ? WHERE id = ?"
	_, err := db.ExecContext(ctx, query, pkg.DomainID, pkg.Private, pkg.Path, pkg.VCS, pkg.Root, pkg.Docs, pkg.Description, pkg.Homepage, pkg.License, pkg.UpdatedAt, pkg.ID)
	if err != nil {
		return err
	}
	return nil
}

func FindPackageByUser(ctx context.Context, db *sqlx.DB, pkg interface{}, id int64, userID string) error {
	query := `
SELECT packages.*
FROM packages
LEFT JOIN domains ON domains.id = packages.domain_id
WHERE packages.id = ? AND domains.user_id = ?
`
	return db.GetContext(ctx, pkg, query, id, userID)
}

// ImportMeta returns go-import meta value.
func (pkg Package) ImportMeta() string {
	return fmt.Sprintf("%s %s %s", pkg.Prefix(), pkg.VCS, pkg.Root)
}

// DocsURL returns the URL of docs.
func (pkg Package) DocsURL() string {
	if pkg.Docs != "" {
		return pkg.Docs
	}

	return fmt.Sprintf("https://pkg.go.dev/%s?tab=doc", pkg.Prefix())
}

// Prefix returns the prefix of package.
func (pkg Package) Prefix() string {
	return fmt.Sprintf("%s/%s", pkg.Domain.Name, pkg.Path)
}

func (pkg Package) Delete(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, "DELETE FROM packages WHERE id = ?", pkg.ID)
	return err
}

func CountPackages(ctx context.Context, db *sqlx.DB, count *int64) error {
	query := "SELECT COUNT(1) FROM packages"
	return db.GetContext(ctx, count, query)
}

func CountPackagesByUser(ctx context.Context, db *sqlx.DB, count *int64, userID int64) error {
	query := "SELECT COUNT(p.id) FROM packages p LEFT JOIN domains d ON d.id = p.domain_id WHERE d.user_id = ?"
	return db.GetContext(ctx, count, query, userID)
}

func CountPackagesByDomainID(ctx context.Context, db *sqlx.DB, count *int64, domainID int64) error {
	query := "SELECT COUNT(id) FROM packages WHERE domain_id = ?"
	return db.GetContext(ctx, count, query, domainID)
}

func FindPackageByDomainAndPath(ctx context.Context, db *sqlx.DB, pkg interface{}, domain, path string) error {
	query := `
SELECT 
	packages.*,
	domains.id as "domain.id",
	domains.name as "domain.name"
FROM packages 
LEFT JOIN domains ON domains.id = packages.domain_id
WHERE domains.name = ? AND packages.path = ?
`
	return db.GetContext(ctx, pkg, query, domain, path)
}
