package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB opens a database.
func NewDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func main() {
}
