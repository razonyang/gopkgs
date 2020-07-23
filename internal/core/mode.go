package core

import "clevergo.tech/osenv"

const (
	ModeProd = "prod"
	ModeDev  = "dev"
)

func IsDevelopMode() bool {
	return ModeDev == osenv.Get("APP_MODE", ModeProd)
}
