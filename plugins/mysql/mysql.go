package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB opens a database.
func NewDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func main() {
}
