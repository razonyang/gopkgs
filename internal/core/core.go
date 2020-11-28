package core

import (
	"clevergo.tech/clevergo"
	_ "github.com/go-sql-driver/mysql"
)

// GetHost returns the request host.
func GetHost(c *clevergo.Context) string {
	if host := c.Request.Header.Get("X-Forwarded-Host"); host != "" {
		return host
	}
	if host := c.Request.Header.Get("X-Original-Host"); host != "" {
		return host
	}

	return c.Host()
}
