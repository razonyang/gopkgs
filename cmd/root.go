package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/razonyang/gopkgs/internal/core"
	"github.com/razonyang/gopkgs/internal/models"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

var (
	cfg = &core.Config{}
	db  *gorm.DB
	app = &cli.App{
		EnableBashCompletion: true,
		Version:              core.Version,
		Name:                 "gopkgs",
		Usage:                "custom and manage your package import path",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.json",
				Usage:   "load configuration from `FILE`",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "enable debug mode",
			},
		},
		Before: func(c *cli.Context) error {
			if err := parseConfig(c.String("config")); err != nil {
				return err
			}

			if err := initDB(c.Bool("debug")); err != nil {
				return err
			}

			return nil
		},
	}
)

func parseConfig(cfgFile string) error {
	cfgData, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(cfgData, cfg); err != nil {
		return err
	}
	return nil
}

// Execute executes commands.
func Execute() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initDB(debug bool) (err error) {
	db, err = core.NewDB(cfg.DB.DSN)
	if err != nil {
		return
	}
	if debug {
		db = db.Debug()
	}
	return db.AutoMigrate(models.Package{})
}
