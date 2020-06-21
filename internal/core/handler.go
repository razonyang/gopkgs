package core

import (
	"net/http"

	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/models"
	"gorm.io/gorm"
)

// RegisterHandlers registers handlers.
func RegisterHandlers(r clevergo.Router, db *gorm.DB) {
	r.Get("/*path", func(c *clevergo.Context) error {
		path := c.Params.String("path")
		if path == "/" {
			return c.NotFound()
		}
		path = c.Host() + path
		pkg, err := models.FindPackage(c.Context(), db, path)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.NotFound()
			}
			return err
		}

		return c.Render(http.StatusOK, "package.tmpl", clevergo.Map{
			"pkg": pkg,
		})
	})
}
