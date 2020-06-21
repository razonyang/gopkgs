package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, showCmd)
}

var showCmd = &cli.Command{
	Name:      "show",
	Usage:     "show package details",
	UsageText: `gopkgs show <prefix>`,
	Action: func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		pkg, err := findPackage(ctx, c)
		if err != nil {
			return err
		}
		fmt.Printf("%s\nvcs : %s\nroot: %s\ndocs: %s\n", pkg.Prefix, pkg.VCS, pkg.Root, pkg.DocsURL())
		return nil
	},
}
