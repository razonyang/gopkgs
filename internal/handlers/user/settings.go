package user

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"clevergo.tech/clevergo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

type settingsForm struct {
	web.Form
	db       *sqlx.DB
	user     *models.User
	Timezone string `json:"timezone" schema:"timezone"`
}

func (f *settingsForm) Validate() error {
	f.Err = validation.ValidateStruct(f,
		validation.Field(&f.Timezone, validation.Required, validation.WithContext(f.validateTimezone)),
	)

	return f.Err
}

func (f *settingsForm) validateTimezone(ctx context.Context, value interface{}) error {
	var tz models.Timezone
	if err := models.FindTimezone(ctx, f.db, &tz, value.(string)); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Invalid timezone %s", value)
		}

		return err
	}
	return nil
}

func (f *settingsForm) update(ctx context.Context) error {
	query := `UPDATE users SET timezone = ? WHERE id = ?`
	_, err := f.db.ExecContext(ctx, query, f.Timezone, f.user.ID)
	return err
}

func (h *Handler) settings(c *clevergo.Context) error {
	ctx := c.Context()
	user := h.User(ctx)
	form := &settingsForm{db: h.DB, Timezone: user.Timezone, user: user}

	if c.IsPost() {
		var err error
		if err = c.Decode(form); err == nil {
			if err = form.update(ctx); err == nil {
				return c.Redirect(http.StatusFound, "/settings")
			}

		}
		h.AddAlert(ctx, alert.NewDanger(err.Error()))
	}

	var timezones []models.Timezone
	if err := models.FindAllTimezones(ctx, h.DB, &timezones); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "user/settings.tmpl", clevergo.Map{
		"form":      form,
		"timezones": timezones,
	})
}
