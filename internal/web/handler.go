package web

import (
	"clevergo.tech/clevergo"
)

type Handler interface {
	Register(clevergo.Router)
}
