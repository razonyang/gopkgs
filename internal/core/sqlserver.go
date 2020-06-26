// +build sqlserver

package core

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// NewDB opens a database.
func NewDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
}
