package user

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"clevergo.tech/clevergo"
	"clevergo.tech/osenv"
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/tasks"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

type forgotPasswordForm struct {
	web.Form
	db    *sqlx.DB
	user  *models.User
	Email string `json:"email" schema:"email"`
}

func (f *forgotPasswordForm) Validate() error {
	f.Err = validation.ValidateStruct(f,
		validation.Field(&f.Email, validation.Required, is.Email, validation.WithContext(f.validateEmail)),
	)

	return f.Err
}

func (f *forgotPasswordForm) validateEmail(ctx context.Context, value interface{}) error {
	var user models.User
	err := models.FindUserByEmail(ctx, f.db, &user, value.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("account doesn't exists")
		}
		return err
	}

	f.user = &user
	return nil
}

var tmplResetPasswordEmail = template.Must(template.New("reset-password-email").Parse(`
<p>Dear {{ .user.Username }}</p>
   <p>Please follow the link to reset your password:</p>
   </p><a href="{{ .link }}">{{ .link }}</a>.</p>
`))

func sendResetPasswordEmail(queue *machinery.Server, user *models.User) error {
	var buf bytes.Buffer
	err := tmplResetPasswordEmail.Execute(&buf, clevergo.Map{
		"user": user,
		"link": fmt.Sprintf("%s/reset-password?token=%s", osenv.Get("APP_URL", "http://localhost:8080"), user.PasswordResetToken.String),
	})
	if err != nil {
		return err
	}

	_, err = queue.SendTaskWithContext(context.Background(), &tasks.Signature{
		UUID: uuid.New().String(),
		Name: "sendMail",
		Args: []tasks.Arg{
			tasks.Arg{"to", "[]string", []string{user.Email}},
			tasks.Arg{"subject", "string", "Reset password."},
			tasks.Arg{"html", "string", buf.String()},
		},
	})
	return err
}

func (f *forgotPasswordForm) reset(ctx context.Context, queue *machinery.Server) error {
	f.user.PasswordResetToken = models.GeneratePasswordResetToken()
	_, err := f.db.ExecContext(ctx, "UPDATE users SET password_reset_token = ? WHERE id = ?", f.user.PasswordResetToken, f.user.ID)
	if err != nil {
		return err
	}

	err = sendResetPasswordEmail(queue, f.user)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) forgotPassword(c *clevergo.Context) error {
	form := &forgotPasswordForm{db: h.DB}

	if c.IsPost() {
		ctx := c.Context()
		if err := c.Decode(form); err == nil {
			if err := form.reset(ctx, h.Queue); err != nil {
				return err
			}

			h.AddAlert(ctx, alert.NewSuccess("Reset password email has sent to your mailbox."))
			return c.Redirect(http.StatusFound, "/forgot-password")
		}
	}

	return c.Render(http.StatusOK, "user/forgot-password.tmpl", clevergo.Map{
		"form": form,
	})
}
