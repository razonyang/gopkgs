package middleware

import (
	"net/http"

	"clevergo.tech/clevergo"
	"clevergo.tech/jsend"
)

func APIErrorHandler(next clevergo.Handle) clevergo.Handle {
	return func(c *clevergo.Context) error {
		if err := next(c); err != nil {
			return c.JSON(http.StatusOK, jsend.NewError(err.Error(), 0, nil))
		}
		return nil
	}
}
