package home

import (
	"net/http"

	"clevergo.tech/clevergo"
)

func (h *Handler) index(c *clevergo.Context) error {
	return c.Render(http.StatusOK, "home/index.tmpl", nil)
}
