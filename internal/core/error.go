package core

import (
	"net/http"
	"strings"

	"clevergo.tech/clevergo"
)

func ErrorHandler(next clevergo.Handle) clevergo.Handle {
	return func(c *clevergo.Context) error {
		if err := next(c); err != nil {
			e, ok := err.(clevergo.Error)
			if !ok {
				e = clevergo.NewError(http.StatusInternalServerError, err)
			}
			return c.Render(http.StatusOK, "home/error.tmpl", clevergo.Map{
				"error":      e,
				"statusText": http.StatusText(e.Status()),
			})
		}

		return nil
	}
}

type MultiError []error

func (me *MultiError) Error() string {
	msgs := make([]string, len(*me))
	for i, err := range *me {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, ";")
}

func (me *MultiError) Add(err error) {
	*me = append(*me, err)
}
