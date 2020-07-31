package badge

import (
	"fmt"
	"net/http"
	"net/url"

	"clevergo.tech/clevergo"
	"clevergo.tech/osenv"
	"pkg.razonyang.com/gopkgs/internal/core"
	"pkg.razonyang.com/gopkgs/internal/web"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/badges", h.index)
	router.Get("/badges/downloads/:interval/*path", h.download)
}

func (h *Handler) index(c *clevergo.Context) error {
	return c.Render(http.StatusOK, "badge/index.tmpl", clevergo.Map{
		"page": web.NewPage("Badges"),
	})
}

func (h *Handler) download(c *clevergo.Context) error {
	u := &url.URL{
		Scheme: "https",
		Host:   "img.shields.io",
		Path:   "/endpoint",
	}
	query := c.QueryParams()
	endpoint := fmt.Sprintf("https://%s/api/badges/downloads/%s%s", osenv.MustGet("APP_HOST"), c.Params.String("interval"), c.Params.String("path"))
	u.RawQuery = fmt.Sprintf("url=%s&%s", endpoint, query.Encode())
	return c.Redirect(http.StatusFound, u.String())
}
