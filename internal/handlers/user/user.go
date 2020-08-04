package user

import (
	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/core"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/login", h.login)
	router.Get("/callback", h.callback)
	router.Get("/logout", h.logout)
	router.Get("/signup", h.signup)
	router.Post("/signup", h.signup)
}
