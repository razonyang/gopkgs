package core

import (
	"clevergo.tech/osenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB opens a database.
func NewDB() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(osenv.Get("MYSQL_DNS")), &gorm.Config{})
}
