package domain

import (
	"fmt"
	"net/http"

	"clevergo.tech/clevergo"
	"clevergo.tech/jsend"
	"pkg.razonyang.com/gopkgs/internal/models"
)

func (h *Handler) delete(c *clevergo.Context) error {
	domain, err := h.findDomain(c)
	if err != nil {
		return err
	}

	ctx := c.Context()
	var pkgCount int64
	if err = models.CountPackagesByDomainID(ctx, h.DB, &pkgCount, domain.ID); err != nil {
		return err
	}
	if pkgCount > 0 {
		return fmt.Errorf("Unable to delete domain %q, since there are %d packages belong to it.", domain.Name, pkgCount)
	}
	if err = domain.Delete(ctx, h.DB); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, jsend.New(nil))
}
