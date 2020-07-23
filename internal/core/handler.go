package core

import (
	"context"

	"clevergo.tech/authmiddleware"
	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"github.com/razonyang/gopkgs/internal/web/alert"
)

type Handler struct {
	DB             *sqlx.DB
	SessionManager *scs.SessionManager
}

func NewHandler(db *sqlx.DB, sessionManager *scs.SessionManager) Handler {
	return Handler{
		DB:             db,
		SessionManager: sessionManager,
	}
}

func (h Handler) UserID(ctx context.Context) string {
	return authmiddleware.GetIdentity(ctx).GetID()
}

func (h Handler) AddAlert(ctx context.Context, a alert.Alert) {
	h.SessionManager.Put(ctx, "alert", a)
}

func (h Handler) AddErrorAlert(ctx context.Context, err error) {
	h.AddAlert(ctx, alert.NewDanger(err.Error()))
}
