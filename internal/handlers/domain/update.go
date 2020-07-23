package domain

import (
	"fmt"
	"net/http"

	"clevergo.tech/clevergo"
)

func (h *Handler) update(c *clevergo.Context) error {
	ctx := c.Context()
	domain, err := h.findDomain(c)
	if err != nil {
		return err
	}
	form := NewForm(h.DB, h.UserID(ctx))
	form.domain = domain
	form.Name = domain.Name
	if c.IsPost() {
		if err := c.Decode(form); err != nil {
			h.AddErrorAlert(ctx, err)
		} else if err := form.Update(ctx); err != nil {
			h.AddErrorAlert(ctx, err)
		} else {
			return c.Redirect(http.StatusFound, fmt.Sprintf("/domain/verify/%d", domain.ID))
		}
	}

	return c.Render(http.StatusOK, "domain/edit.tmpl", clevergo.Map{
		"domain": domain,
		"form":   form,
	})
}
