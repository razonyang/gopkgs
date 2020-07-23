package home

import (
	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/core"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/", h.index)
}
