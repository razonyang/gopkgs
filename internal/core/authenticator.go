package core

import (
	"errors"
	"net/http"

	"clevergo.tech/auth"
	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
)

type SessionAuthenticator struct {
	sessionManager *scs.SessionManager
	db             *sqlx.DB
}

func NewSessionAuthenticator(sessionManager *scs.SessionManager, db *sqlx.DB) *SessionAuthenticator {
	return &SessionAuthenticator{
		sessionManager: sessionManager,
		db:             db,
	}
}

// Authenticates the current user.
func (a *SessionAuthenticator) Authenticate(r *http.Request, w http.ResponseWriter) (auth.Identity, error) {
	ctx := r.Context()
	authKey := a.sessionManager.GetString(ctx, "auth_key")
	if authKey == "" {
		return nil, errors.New("no logged user")
	}
	var user models.User
	if err := models.FindUserByAuthKey(ctx, a.db, &user, authKey); err != nil {
		return nil, err
	}
	return &user, nil
}

// Challenge generates challenges upon authentication failure.
func (a *SessionAuthenticator) Challenge(*http.Request, http.ResponseWriter) {
}
