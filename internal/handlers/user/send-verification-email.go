package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"clevergo.tech/clevergo"
	"github.com/RichardKnop/machinery/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

type sendVerificationEmailForm struct {
	web.Form
	db    *sqlx.DB
	user  *models.User
	Email string `json:"email" schema:"email"`
}

func (f *sendVerificationEmailForm) Validate() error {
	f.Err = validation.ValidateStruct(f,
		validation.Field(&f.Email, validation.Required, is.Email, validation.WithContext(f.validateEmail)),
	)
	return f.Err
}

func (f *sendVerificationEmailForm) validateEmail(ctx context.Context, value interface{}) error {
	var user models.User
	if err := models.FindUserByEmail(ctx, f.db, &user, value.(string)); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("account doesn't exists")
		}
		return err
	}
	if user.EmailVerified {
		return errors.New("email has already been verified")
	}
	f.user = &user
	return nil
}

func (f *sendVerificationEmailForm) send(ctx context.Context, queue *machinery.Server) error {
	// regenerate token.
	f.user.VerificationToken = models.GenerateVerificationToken()
	query := "UPDATE users SET verification_token = ? WHERE id = ?"
	if _, err := f.db.ExecContext(ctx, query, f.user.VerificationToken, f.user.ID); err != nil {
		return err
	}

	return sendVerificationMail(queue, f.user)
}

func (h *Handler) sendVerificationEmail(c *clevergo.Context) error {
	form := &sendVerificationEmailForm{db: h.DB}

	if c.IsPost() {
		ctx := c.Context()
		if err := c.Decode(form); err != nil {
		} else if err := form.send(ctx, h.Queue); err != nil {
			return err
		} else {
			h.AddAlert(ctx, alert.NewSuccess("Verification email has been sent, please check your mailbox."))
			return c.Redirect(http.StatusFound, "/send-verification-email")
		}
	}

	return c.Render(http.StatusFound, "user/send-verification-email.tmpl", clevergo.Map{
		"form": form,
	})
}
