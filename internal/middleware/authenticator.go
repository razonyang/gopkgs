package middleware

import (
	"net/http"

	"clevergo.tech/authmidware"
	"clevergo.tech/clevergo"
)

// IsAuthenticated redirects to the given URL if user isn't login.
func IsAuthenticated(loginURL string, skipper clevergo.Skipper) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(c *clevergo.Context) error {
			if user := authmidware.GetIdentity(c.Context()); user == nil && !skipper(c) {
				return c.Redirect(http.StatusTemporaryRedirect, loginURL)
			}
			return next(c)
		}
	}
}

type CompositeSkipper struct {
	skippers []clevergo.Skipper
}

func (s *CompositeSkipper) CanSkip(c *clevergo.Context) bool {
	for _, skipper := range s.skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

func NewCompositeSkipper(skippers ...clevergo.Skipper) clevergo.Skipper {
	cs := &CompositeSkipper{skippers: skippers}
	return cs.CanSkip
}
