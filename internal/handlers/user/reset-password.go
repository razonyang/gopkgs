package user

import (
	"context"
	"database/sql"
	"net/http"

	"clevergo.tech/clevergo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

type resetPasswordForm struct {
	web.Form
	db       *sqlx.DB
	user     *models.User
	Password string `json:"password" schema:"password"`
}

func (f *resetPasswordForm) Validate() error {
	f.Err = validation.ValidateStruct(f,
		validation.Field(&f.Password, validation.Required, validation.RuneLength(6, 0),
			validation.Match(regPasswordLowercase).Error("password must contain at least one lowercase letter"),
			validation.Match(regPasswordUppercase).Error("password must contain at least one uppercase letter"),
			validation.Match(regPasswordDigit).Error("password must contain at least one digit"),
		),
	)

	return f.Err
}

func (f *resetPasswordForm) reset(ctx context.Context) error {
	err := f.user.SetPassword(ctx, f.db, f.Password)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) resetPassword(c *clevergo.Context) error {
	token := c.QueryParam("token")
	ctx := c.Context()
	var user models.User
	err := models.FindUserByPasswordResetToken(ctx, h.DB, &user, token)
	if err != nil {
		if err == sql.ErrNoRows {
			h.AddAlert(ctx, alert.NewDanger("invalid token"))
			return c.Redirect(http.StatusFound, "/forgot-password")
		}
		return err
	}
	form := &resetPasswordForm{db: h.DB, user: &user}

	if c.IsPost() {
		if err := c.Decode(form); err == nil {
			if err := form.reset(ctx); err != nil {
				return err
			}

			h.AddAlert(ctx, alert.NewSuccess("Reset password successfully."))
			return c.Redirect(http.StatusFound, "/login")
		}
	}

	return c.Render(http.StatusOK, "user/reset-password.tmpl", clevergo.Map{
		"token": token,
		"form":  form,
	})
}
