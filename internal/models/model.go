package models

import (
	"strings"
	"time"
)

type Model struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Errors []error

func (errs *Errors) Error() string {
	msgs := make([]string, len(*errs))
	for i, err := range *errs {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, ";")
}

func (errs *Errors) Add(err error) {
	*errs = append(*errs, err)
}
