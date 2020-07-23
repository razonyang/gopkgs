package pkg

import (
	"net/http"

	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/models"
	"github.com/razonyang/gopkgs/internal/web/alert"
)

func (h *Handler) update(c *clevergo.Context) error {
	ctx := c.Context()
	pkg, err := h.findPackage(c)
	if err != nil {
		return err
	}

	form := newFormPkg(h.DB, h.UserID(ctx), pkg)
	if c.IsPost() {
		if err := c.Decode(form); err != nil {
			h.AddAlert(ctx, alert.NewDanger(err.Error()))
		} else if err := form.Update(ctx); err != nil {
			h.AddAlert(ctx, alert.NewDanger(err.Error()))
		} else {
			return c.Redirect(http.StatusFound, "/package")
		}
	}

	domains, err := h.findDomains(c)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "package/update.tmpl", clevergo.Map{
		"domains": domains,
		"vcs":     models.VCSSet,
		"form":    form,
	})
}
