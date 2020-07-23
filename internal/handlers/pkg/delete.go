package pkg

import (
	"net/http"

	"clevergo.tech/clevergo"
	"clevergo.tech/jsend"
)

func (h *Handler) delete(c *clevergo.Context) error {
	pkg, err := h.findPackage(c)
	if err != nil {
		return err
	}
	ctx := c.Context()
	if err = pkg.Delete(ctx, h.DB); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jsend.New(nil))
}
