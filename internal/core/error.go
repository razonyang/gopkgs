package core

import (
	"net/http"
	"strings"

	"clevergo.tech/clevergo"
)

func ErrorHandler(next clevergo.Handle) clevergo.Handle {
	return func(c *clevergo.Context) error {
		return next(c)
		if err := next(c); err != nil {
			return c.Render(http.StatusOK, "error.tmpl", clevergo.Map{
				"error": err,
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
