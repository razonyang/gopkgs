package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"clevergo.tech/plugins"
	"github.com/razonyang/gopkgs/internal/core"
	"github.com/razonyang/gopkgs/internal/models"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

const version = "v0.1.0"

var (
	cfg           = &core.Config{}
	db            *gorm.DB
	pluginManager *plugins.Manager
	app           = &cli.App{
		EnableBashCompletion: true,
		Version:              version,
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

			pluginManager = plugins.New(cfg.Plugins)

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

func initDB(debug bool) error {
	symbol, err := pluginManager.Lookup(cfg.DB.Driver+".so", "NewDB")
	if err != nil {
		return err
	}
	f := symbol.(func(string) (*gorm.DB, error))
	db, err = f(cfg.DB.DSN)
	if err != nil {
		return err
	}
	if debug {
		db = db.Debug()
	}
	return db.AutoMigrate(models.Package{})
}
