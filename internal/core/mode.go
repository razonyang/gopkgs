package core

import "clevergo.tech/osenv"

func IsDevelopMode() bool {
	return osenv.Get("DEBUG") != ""
}
