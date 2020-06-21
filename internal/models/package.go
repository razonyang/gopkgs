package models

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
)

// Version control system constants.
const (
	VCSGit    = "git"
	VCSSvn    = "svn"
	VCSBzr    = "bzr"
	VCSHg     = "hg"
	VCSFossil = "fossil"
)

// Package is a model that mapping to table "packages".
type Package struct {
	Prefix string `gorm:"type:varchar(64);PRIMARY_KEY" json:"prefix"`
	VCS    string `gorm:"type:varchar(6);NOT NULL" json:"vcs"`
	Root   string `gorm:"type:varchar(256);NOT NULL" json:"root"`
	Docs   string `gorm:"type:varchar(256);NOT NULL" json:"docs"`
}

// Validate validates package.
func (pkg *Package) Validate() error {
	return validation.ValidateStruct(pkg,
		validation.Field(&pkg.Prefix, validation.Required),
		validation.Field(&pkg.VCS, validation.Required, validation.In(VCSGit, VCSSvn, VCSBzr, VCSHg, VCSFossil)),
		validation.Field(&pkg.Root, validation.Required),
		validation.Field(&pkg.Docs, is.URL),
	)
}

// Save saves package into database.
func (pkg *Package) Save(ctx context.Context, db *gorm.DB) error {
	if err := pkg.Validate(); err != nil {
		return err
	}
	return db.WithContext(ctx).Save(pkg).Error
}

// IsPackagePrefixTaken verifies whether the prefix is already taken.
func IsPackagePrefixTaken(ctx context.Context, db *gorm.DB, prefix string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&Package{}).Where("prefix=?", prefix).Count(&count).Error
	return count > 0, err
}

// FindPackage finds package by prefix.
func FindPackage(ctx context.Context, db *gorm.DB, prefix string) (pkg *Package, err error) {
	pkg = &Package{}
	err = db.WithContext(ctx).Where("prefix=?", prefix).First(&pkg).Error
	return
}

// FindPackageByPath finds package by path.
func FindPackageByPath(ctx context.Context, db *gorm.DB, path string) (pkg *Package, err error) {
	var prefixes []string
	parts := strings.Split(path, "/")
	for i := 2; i <= len(parts); i++ {
		prefixes = append(prefixes, strings.Join(parts[:i], "/"))
	}
	pkg = &Package{}
	err = db.WithContext(ctx).Where("prefix IN ?", prefixes).Order("prefix DESC").First(&pkg).Error
	return
}

// ImportMeta returns go-import meta value.
func (pkg *Package) ImportMeta() string {
	return fmt.Sprintf("%s %s %s", pkg.Prefix, pkg.VCS, pkg.Root)
}

// DocsURL returns the URL of docs.
func (pkg *Package) DocsURL() string {
	if pkg.Docs != "" {
		return pkg.Docs
	}

	return fmt.Sprintf("https://pkg.go.dev/%s?tab=doc", pkg.Prefix)
}
