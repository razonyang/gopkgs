package pkg

import (
	"net/http"

	"clevergo.tech/clevergo"
)

func (h *Handler) view(c *clevergo.Context) error {
	pkg, err := h.findPackage(c)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "package/view.tmpl", clevergo.Map{
		"pkg": pkg,
	})
}
