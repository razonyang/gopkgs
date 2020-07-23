package report

import (
	"net/http"

	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/web"
)

func (h *Handler) index(c *clevergo.Context) error {
	return c.Render(http.StatusOK, "report/index.tmpl", clevergo.Map{
		"page": web.NewPage("Report"),
	})
}
