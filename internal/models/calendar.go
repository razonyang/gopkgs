package models

import "time"

type Calendar struct {
	ID time.Time `db:"id" json:"id"`
}
