package domain

import (
	"net/http"

	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/core"
	"pkg.razonyang.com/gopkgs/internal/models"
)

func (h *Handler) challenge(c *clevergo.Context) error {
	ctx := c.Context()
	var domain models.Domain
	if err := models.FindDomainByChallengeTXT(ctx, h.DB, &domain, core.GetHost(c), c.Params.String("token")); err != nil {
		c.Logger().Error(err)
		return c.NotFound()
	}

	return c.String(http.StatusOK, domain.ChallengeTXT)
}
