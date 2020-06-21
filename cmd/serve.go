package cmd

import (
	"clevergo.tech/clevergo"
	"clevergo.tech/jetpackr"
	"clevergo.tech/jetrenderer"
	"github.com/CloudyKit/jet/v3"
	"github.com/gobuffalo/packr/v2"
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
		box := packr.New("views", "../views")
		app.Renderer = jetrenderer.New(jet.NewHTMLSetLoader(jetpackr.New(box)))
		core.RegisterHandlers(app, db)
		return app.Run(cfg.Addr)
	},
}
