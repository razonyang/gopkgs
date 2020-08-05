package user

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"clevergo.tech/authmiddleware"
	"clevergo.tech/clevergo"
	"github.com/alexedwards/scs/v2"
	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

var (
	errIncorrectAccountOrPassword = errors.New("incorrect account or password")
)

type loginForm struct {
	web.Form
	db   *sqlx.DB
	user *models.User

	Account  string `json:"account" schema:"account"`
	Password string `json:"password" schema:"password"`
}

func (f *loginForm) Validate() error {
	f.Err = validation.ValidateStruct(f,
		validation.Field(&f.Account, validation.Required),
		validation.Field(&f.Password, validation.Required, validation.WithContext(f.validatePassword)),
	)
	return f.Err
}

func (f *loginForm) validatePassword(ctx context.Context, value interface{}) (err error) {
	var user models.User
	if govalidator.IsExistingEmail(f.Account) {
		err = models.FindUserByEmail(ctx, f.db, &user, f.Account)
	} else {
		err = models.FindUserByUsername(ctx, f.db, &user, f.Account)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return errIncorrectAccountOrPassword
		}
		return err
	}

	if err := user.ValidatePassword(value.(string)); err != nil {
		return errIncorrectAccountOrPassword
	}

	f.user = &user
	return nil
}

func (f *loginForm) login(ctx context.Context, sessionManager *scs.SessionManager) error {
	sessionManager.Put(ctx, "auth_user", f.user)
	return nil
}

func (h *Handler) login(c *clevergo.Context) error {
	ctx := c.Context()
	if user := authmiddleware.GetIdentity(ctx); user != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}

	form := &loginForm{db: h.DB}
	if c.IsPost() {
		if err := c.Decode(form); err == nil {
			if err := form.login(ctx, h.SessionManager); err == nil {
				h.AddAlert(ctx, alert.NewSuccess("Login successfully."))
				return c.Redirect(http.StatusFound, "/dashboard")
			}
		}
	}

	return c.Render(http.StatusOK, "user/login.tmpl", clevergo.Map{
		"form": form,
	})
}
