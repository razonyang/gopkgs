package user

import (
	"net/http"

	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

func (h *Handler) logout(c *clevergo.Context) error {
	ctx := c.Context()
	// remove user information from session
	h.SessionManager.Remove(ctx, "auth_user")
	h.AddAlert(ctx, alert.NewSuccess("Logout successfully."))
	return c.Redirect(http.StatusFound, "/")
}
