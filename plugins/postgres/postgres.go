package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB opens a database.
func NewDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func main() {
}
