package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/razonyang/gopkgs/internal/models"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func init() {
	app.Commands = append(app.Commands, setVCSCmd)
	app.Commands = append(app.Commands, setRootCmd)
	app.Commands = append(app.Commands, setDocsCmd)
}

var setVCSCmd = &cli.Command{
	Name:      "set-vcs",
	Usage:     "modify the VCS of package",
	UsageText: "gopkgs set-vcs <prefix> <vcs>",
	Action: func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		pkg, err := findPackage(ctx, c)
		if err != nil {
			return err
		}
		pkg.VCS = c.Args().Get(1)
		return pkg.Save(ctx, db)
	},
}

var setRootCmd = &cli.Command{
	Name:      "set-root",
	Usage:     "modify root of repository",
	UsageText: "gopkgs set-root <prefix> <repo-root>",
	Action: func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		pkg, err := findPackage(ctx, c)
		if err != nil {
			return err
		}
		pkg.Root = c.Args().Get(1)
		return pkg.Save(ctx, db)
	},
}

var setDocsCmd = &cli.Command{
	Name:      "set-docs",
	Usage:     "modify the URL of documentations",
	UsageText: "gopkgs set-docs <prefix> <docs-url>",
	Action: func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		pkg, err := findPackage(ctx, c)
		if err != nil {
			return err
		}
		pkg.Docs = c.Args().Get(1)
		return pkg.Save(ctx, db)
	},
}

var errNoPackage = errors.New("no pakcage prefix provided")

func findPackage(ctx context.Context, c *cli.Context) (*models.Package, error) {
	if c.NArg() == 0 {
		return nil, errNoPackage
	}
	prefix := c.Args().Get(0)
	pkg, err := models.FindPackage(ctx, db, prefix)
	if err != nil && err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("package %q doesn't exists", prefix)
	}
	return pkg, err
}
