package core

import (
	"context"

	"clevergo.tech/authmiddleware"
	"github.com/RichardKnop/machinery/v2"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/ristretto"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web/alert"
)

type Handler struct {
	DB             *sqlx.DB
	SessionManager *scs.SessionManager
	Queue          *machinery.Server
	Cache          *ristretto.Cache
}

func NewHandler(db *sqlx.DB, sessionManager *scs.SessionManager, queue *machinery.Server, cache *ristretto.Cache) Handler {
	return Handler{
		DB:             db,
		SessionManager: sessionManager,
		Queue:          queue,
		Cache:          cache,
	}
}

func (h Handler) User(ctx context.Context) *models.User {
	return authmiddleware.GetIdentity(ctx).(*models.User)
}

func (h Handler) UserID(ctx context.Context) int64 {
	return h.User(ctx).ID
}

func (h Handler) AddAlert(ctx context.Context, a alert.Alert) {
	h.SessionManager.Put(ctx, "alert", a)
}

func (h Handler) AddErrorAlert(ctx context.Context, err error) {
	h.AddAlert(ctx, alert.NewDanger(err.Error()))
}
