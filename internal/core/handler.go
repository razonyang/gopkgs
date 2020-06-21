package core

import (
	"html/template"

	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/models"
	"gorm.io/gorm"
)

var (
	tmplHTML = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{ .pkg.ImportMeta }}">
<meta http-equiv="refresh" content="0; url={{ .pkg.DocsURL }}">
<title>Package {{ .pkg.Prefix }}</title>
</head>
<body>
Nothing to see here; <a href="{{ .pkg.DocsURL }}">move along</a>.
</body>
</html>
`
	tmpl = template.Must(template.New("pacakge").Parse(tmplHTML))
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

		c.SetContentTypeHTML()
		return tmpl.Execute(c.Response, clevergo.Map{
			"pkg": pkg,
		})
	})
}
