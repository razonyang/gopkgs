package user

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"clevergo.tech/clevergo"
	"clevergo.tech/osenv"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/RichardKnop/machinery/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

var (
	regUsername          = regexp.MustCompile(`[a-zA-Z0-9]+`)
	regPasswordLowercase = regexp.MustCompile(`[a-z]+`)
	regPasswordUppercase = regexp.MustCompile(`[A-Z]+`)
	regPasswordDigit     = regexp.MustCompile(`[0-9]+`)
)

type signupForm struct {
	web.Form
	db       *sqlx.DB
	Username string `json:"username" schema:"username"`
	Email    string `json:"email" schema:"email"`
	Password string `json:"password" schema:"password"`
}

func (f *signupForm) Validate() error {
	f.Form.Err = validation.ValidateStruct(f,
		validation.Field(&f.Username, validation.Required, validation.RuneLength(5, 0), validation.Match(regUsername),
			validation.WithContext(f.isUsernameTaken),
		),
		validation.Field(&f.Email, validation.Required, is.Email, validation.WithContext(f.isEmailTaken)),
		validation.Field(&f.Password, validation.Required, validation.RuneLength(6, 0),
			validation.Match(regPasswordLowercase).Error("password must contain at least one lowercase letter"),
			validation.Match(regPasswordUppercase).Error("password must contain at least one uppercase letter"),
			validation.Match(regPasswordDigit).Error("password must contain at least one digit"),
		),
	)

	return f.Form.Err
}

func (f *signupForm) isUsernameTaken(ctx context.Context, value interface{}) error {
	var count int64
	if err := f.db.GetContext(ctx, &count, "SELECT COUNT(1) FROM users WHERE username = ?", value); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("username %s has been taken", value)
	}
	return nil
}

func (f *signupForm) isEmailTaken(ctx context.Context, value interface{}) error {
	var count int64
	if err := f.db.GetContext(ctx, &count, "SELECT COUNT(1) FROM users WHERE email = ?", value); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("email %s has been taken", value)
	}
	return nil
}

func (f *signupForm) signup(ctx context.Context) (user *models.User, err error) {
	user, err = models.NewUser(f.Username, f.Email, f.Password)
	if err != nil {
		return
	}
	if err = user.Insert(ctx, f.db); err != nil {
		return
	}

	return user, nil
}

func newSignupForm(db *sqlx.DB) *signupForm {
	return &signupForm{
		db: db,
	}
}

var tmplVerificationMail = template.Must(template.New("verification-email").Parse(`
<p>Dear {{ .user.Username }}</p>
   <p>Please follow the link to activate your account:</p>
   </p><a href="{{ .link }}">{{ .link }}</a>.</p>
`))

func (h *Handler) signup(c *clevergo.Context) error {
	form := newSignupForm(h.DB)
	if c.IsPost() {
		ctx := c.Context()
		if err := c.Decode(form); err != nil {
		} else if user, err := form.signup(ctx); err != nil {
		} else if sendVerificationMail(h.Queue, user); err != nil {
			log.Println(err)
		} else {
			return c.Redirect(http.StatusFound, "/login")
		}
	}

	return c.Render(http.StatusOK, "user/signup.tmpl", clevergo.Map{
		"form": form,
	})
}

func sendVerificationMail(queue *machinery.Server, user *models.User) error {
	var buf bytes.Buffer
	err := tmplVerificationMail.Execute(&buf, clevergo.Map{
		"user": user,
		"link": fmt.Sprintf("%s/verify-email?token=%s", osenv.Get("APP_URL", "http://localhost:8080"), user.VerificationToken.String),
	})
	if err != nil {
		return err
	}

	_, err = queue.SendTaskWithContext(context.Background(), &tasks.Signature{
		UUID: uuid.New().String(),
		Name: "sendMail",
		Args: []tasks.Arg{
			tasks.Arg{"to", "[]string", []string{user.Email}},
			tasks.Arg{"subject", "string", "Please verify your account."},
			tasks.Arg{"html", "string", buf.String()},
		},
	})
	return err
}
