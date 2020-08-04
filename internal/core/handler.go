package core

import (
	"context"

	"clevergo.tech/authmiddleware"
	"github.com/RichardKnop/machinery/v2"
	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

type Handler struct {
	DB             *sqlx.DB
	SessionManager *scs.SessionManager
	Queue          *machinery.Server
}

func NewHandler(db *sqlx.DB, sessionManager *scs.SessionManager, queue *machinery.Server) Handler {
	return Handler{
		DB:             db,
		SessionManager: sessionManager,
		Queue:          queue,
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
