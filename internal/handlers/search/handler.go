package search

import (
	"clevergo.tech/clevergo"
	"pkg.razonyang.com/gopkgs/internal/core"
)

type Handler struct {
	core.Handler
}

func (h *Handler) Register(router clevergo.Router) {
	router.Get("/search", h.index)
}
