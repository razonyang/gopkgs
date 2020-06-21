package cmd

import (
	"clevergo.tech/clevergo"
	"github.com/razonyang/gopkgs/internal/core"
	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, serveCmd)
}

var serveCmd = &cli.Command{
	Name:  "serve",
	Usage: "start a HTTP server",
	Action: func(c *cli.Context) error {
		app := clevergo.New()
		core.RegisterHandlers(app, db)
		return app.Run(cfg.Addr)
	},
}
