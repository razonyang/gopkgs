package middleware

import (
	"net/http"

	"clevergo.tech/clevergo"
)

func Host(host string, skipper clevergo.Skipper) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(c *clevergo.Context) error {
			if host != c.Host() && !skipper(c) {
				url := c.Request.URL
				url.Host = host
				url.Scheme = "http"
				if c.Request.TLS != nil {
					url.Scheme = "https"
				}
				return c.Redirect(http.StatusMovedPermanently, url.String())
			}
			return next(c)
		}
	}
}
