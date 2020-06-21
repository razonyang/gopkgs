package cmd

import (
	"context"

	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, removeCmd)
}

var removeCmd = &cli.Command{
	Name:      "remove",
	Usage:     "remove a package",
	UsageText: "gopkgs remove <prefix>",
	Action: func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		pkg, err := findPackage(ctx, c)
		if err != nil {
			return err
		}
		return db.Delete(&pkg).Error
	},
}
