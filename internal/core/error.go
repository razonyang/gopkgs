package core

import (
	"net/http"

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
