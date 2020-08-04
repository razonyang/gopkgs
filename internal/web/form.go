package web

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Form struct {
	Err error
}

func (f Form) HasError(field string) (found bool) {
	errs, ok := f.Err.(validation.Errors)
	if ok {
		_, found = errs[field]
	}

	return
}

func (f Form) GetError(field string) error {
	errs, ok := f.Err.(validation.Errors)
	if ok {
		if err, ok := errs[field]; ok {
			return err
		}
	}

	return nil
}
