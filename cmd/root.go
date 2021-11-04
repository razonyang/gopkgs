package cmd

import (
	"log"
	"os"
	"strconv"

	"clevergo.tech/osenv"
	"github.com/RichardKnop/machinery/v2"
	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	"github.com/RichardKnop/machinery/v2/config"
	locksiface "github.com/RichardKnop/machinery/v2/locks/iface"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	db    *sqlx.DB
	queue *machinery.Server
	app   = &cli.App{
		EnableBashCompletion: true,
		Name:                 "gopkgs",
		Usage:                "Go Packages",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "configs/prod",
			},
		},
		Before: func(c *cli.Context) error {
			if err := godotenv.Load(c.String("config")); err != nil {
				return err
			}

			if err := initialize(); err != nil {
				return err
			}

			return nil
		},
	}
)

// Execute executes commands.
func Execute() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initialize() (err error) {
	db, err = sqlx.Open("mysql", osenv.Get("MYSQL_DNS"))
	if err != nil {
		return
	}

	queueCfg, err := config.NewFromEnvironment()
	if err != nil {
		return
	}
	queueCfg.DefaultQueue = "gopkgs_tasks"
	queueCfg.TLSConfig = nil
	queueCfg.Redis = &config.RedisConfig{}

	redisHost := osenv.Get("REDIS_ADDR", "localhost:6379")
	redisPassword := osenv.Get("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(osenv.Get("REDIS_DATABASE", "0"))
	broker := redisbroker.New(queueCfg, redisHost, redisPassword, "", redisDB)
	backend := redisbackend.New(queueCfg, redisHost, redisPassword, "", redisDB)
	var lock locksiface.Lock
	queue = machinery.NewServer(queueCfg, broker, backend, lock)

	return nil
}
