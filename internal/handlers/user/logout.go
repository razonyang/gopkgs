package user

import (
	"net/http"
	"net/url"
	"os"

	"clevergo.tech/clevergo"
)

func (h *Handler) logout(c *clevergo.Context) error {
	ctx := c.Context()

	// remove user information from session
	h.SessionManager.Remove(ctx, "auth_user")

	domain := os.Getenv("AUTH0_DOMAIN")

	logoutUrl, err := url.Parse("https://" + domain)

	if err != nil {
		return err
	}

	logoutUrl.Path += "/v2/logout"
	parameters := url.Values{}

	var scheme string
	if c.Request.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + c.Host())
	if err != nil {
		return err
	}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	return c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
