package user

import (
	"database/sql"
	"net/http"

	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

func (h *Handler) verifyEmail(c *clevergo.Context) error {
	token := c.QueryParam("token")
	ctx := c.Context()
	var user models.User
	err := models.FindUserByVerificationToken(ctx, h.DB, &user, token)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if user.ID > 0 && !user.IsVerificationExpired() {
		user.EmailVerified = true
		_, err = h.DB.ExecContext(ctx, "UPDATE users SET email_verified = ?, verification_token=NULL WHERE id = ?", user.EmailVerified, user.ID)
		if err != nil {
			return err
		}

		// update session
		h.SessionManager.Put(ctx, "auth_user", user)

		h.AddAlert(ctx, alert.NewSuccess("Email verified."))

		return c.Redirect(http.StatusFound, "/login")
	}

	h.AddAlert(ctx, alert.NewDanger("Invalid token."))
	return c.Redirect(http.StatusFound, "/send-verification-email")
}
