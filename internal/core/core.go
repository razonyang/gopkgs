package core

import (
	"clevergo.tech/clevergo"
	_ "github.com/go-sql-driver/mysql"
)

const Version = "v0.3.0"

// GetHost returns the request host.
func GetHost(c *clevergo.Context) string {
	if host := c.Request.Header.Get("X-Forwarded-Host"); host != "" {
		return host
	}

	return c.Host()
}
