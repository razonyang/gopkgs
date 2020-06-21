package cmd

import (
	"context"
	"fmt"

	"github.com/razonyang/gopkgs/internal/models"
	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, addCmd)
}

var addCmd = &cli.Command{
	Name:      "add",
	Usage:     "add a new package",
	UsageText: "gopkgs add <prefix> <vcs> <repo-root> [<docs-url>]",
	Description: `prefix: the prefix of import path.
	 vcs: bzr, fossil, git, hg, svn.
	 repo-root: the location of your repository.
	 docs-url: URL of documentations, optional.

	 gopkgs add clevergo.tech/clevergo git https://github.com/clevergo/clevergo
	`,
	Action: func(c *cli.Context) error {
		args := c.Args()
		prefix := args.Get(0)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		taken, err := models.IsPackagePrefixTaken(ctx, db, prefix)
		if err != nil {
			return err
		}
		if taken {
			return fmt.Errorf("package %q already exists", prefix)
		}

		pkg := models.Package{
			Prefix: prefix,
			VCS:    args.Get(1),
			Root:   args.Get(2),
			Docs:   args.Get(3),
		}
		if err := pkg.Validate(); err != nil {
			return err
		}
		return db.WithContext(ctx).Create(&pkg).Error
	},
}
