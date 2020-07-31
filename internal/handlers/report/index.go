package report

import (
	"net/http"

	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/web"
)

func (h *Handler) index(c *clevergo.Context) error {
	return c.Render(http.StatusOK, "report/index.tmpl", clevergo.Map{
		"page": web.NewPage("Report"),
	})
}
