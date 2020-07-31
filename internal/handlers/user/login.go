package user

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"clevergo.tech/authmiddleware"
	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/core"
)

func (h *Handler) login(c *clevergo.Context) error {
	ctx := c.Context()
	if user := authmiddleware.GetIdentity(ctx); user != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}

	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	state := base64.StdEncoding.EncodeToString(b)

	h.SessionManager.Put(ctx, "auth_state", state)
	authenticator, err := core.NewAuthenticator(ctx)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, authenticator.Config.AuthCodeURL(state))
}
