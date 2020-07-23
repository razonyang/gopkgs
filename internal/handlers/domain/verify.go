package domain

import (
	"net/http"

	"clevergo.tech/clevergo"
	"clevergo.tech/jsend"
	"clevergo.tech/osenv"
)

func (h *Handler) verify(c *clevergo.Context) error {
	domain, err := h.findDomain(c)
	if err != nil {
		return err
	}

	if c.IsAJAX() {
		if !domain.Verified {
			if err := domain.Challenge(c.Context(), h.DB); err != nil {
				return c.JSON(http.StatusOK, jsend.NewError(err.Error(), 0, nil))
			}
		}

		return c.JSON(http.StatusOK, jsend.New(clevergo.Map{
			"verified": domain.Verified,
		}))
	}

	return c.Render(http.StatusOK, "domain/verify.tmpl", clevergo.Map{
		"domain": domain,
		"host":   osenv.MustGet("APP_HOST"),
	})
}
