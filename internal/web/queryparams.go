package web

import (
	"clevergo.tech/clevergo"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.IgnoreUnknownKeys(true)
}

func DecodeQueryParams(c *clevergo.Context, dest interface{}) error {
	return decoder.Decode(dest, c.QueryParams())
}
