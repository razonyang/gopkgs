package domain

import (
	"fmt"
	"net/http"

	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/models"
)

func (h *Handler) create(c *clevergo.Context) error {
	ctx := c.Context()
	user := h.User(ctx).(*models.User)
	if !user.EmailVerified {
		return c.Redirect(http.StatusFound, "/send-verification-email")
	}
	form := NewForm(h.DB, user.GetID())
	if c.IsPost() {
		if err := c.Decode(form); err != nil {
			h.AddErrorAlert(ctx, err)
		} else if domain, err := form.Create(ctx); err != nil {
			h.AddErrorAlert(ctx, err)
		} else {
			return c.Redirect(http.StatusFound, fmt.Sprintf("/domain/verify/%d", domain.ID))
		}
	}

	return c.Render(http.StatusOK, "domain/create.tmpl", clevergo.Map{
		"form": form,
	})
}
