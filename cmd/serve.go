package cmd

import (
	"encoding/gob"
	"io"
	"io/ioutil"
	stdlog "log"
	"reflect"
	"strconv"
	"time"

	"clevergo.tech/authmiddleware"
	"clevergo.tech/clevergo"
	"clevergo.tech/form"
	"clevergo.tech/i18nmiddleware"
	"clevergo.tech/jetpackr"
	"clevergo.tech/jetrenderer"
	"clevergo.tech/jetsprig"
	"clevergo.tech/log"
	"clevergo.tech/osenv"
	"clevergo.tech/pprof"
	"github.com/BurntSushi/toml"
	"github.com/CloudyKit/jet/v5"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/ristretto"
	"github.com/gobuffalo/packr/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/justinas/nosurf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"pkg.razonyang.com/gopkgs/internal/core"
	"pkg.razonyang.com/gopkgs/internal/handlers/api"
	"pkg.razonyang.com/gopkgs/internal/handlers/badge"
	"pkg.razonyang.com/gopkgs/internal/handlers/dashboard"
	"pkg.razonyang.com/gopkgs/internal/handlers/domain"
	"pkg.razonyang.com/gopkgs/internal/handlers/home"
	"pkg.razonyang.com/gopkgs/internal/handlers/pkg"
	"pkg.razonyang.com/gopkgs/internal/handlers/report"
	"pkg.razonyang.com/gopkgs/internal/handlers/search"
	"pkg.razonyang.com/gopkgs/internal/handlers/trending"
	"pkg.razonyang.com/gopkgs/internal/handlers/user"
	"pkg.razonyang.com/gopkgs/internal/middleware"
	"pkg.razonyang.com/gopkgs/internal/stringhelper"
	"pkg.razonyang.com/gopkgs/internal/web"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

func init() {
	app.Commands = append(app.Commands, serveCmd)
}

var serveCmd = &cli.Command{
	Name:  "serve",
	Usage: "start a HTTP server",
	Action: func(c *cli.Context) error {
		db.DB.SetMaxIdleConns(10)
		db.DB.SetMaxOpenConns(100)
		startCrond()
		go startQueue()

		logger, err := provideLogger()
		if err != nil {
			return err
		}
		app := clevergo.Pure()
		app.Logger = logger
		app.Use(clevergo.Logging(clevergo.LoggingLogger(logger)))

		cache := provideCache()

		app.Decoder = form.New()
		sessionManager := provideSessionManager()
		app.Use(
			clevergo.Recovery(),
			clevergo.ServerHeader("CleverGo"),
			clevergo.WrapHH(sessionManager.LoadAndSave),
			core.ErrorHandler,
			// clevergo.WrapHH(middleware.Minify()),
			authmiddleware.New(core.NewSessionAuthenticator(sessionManager, db)),
			i18nmiddleware.New(provideI18N()),
			middleware.GoGet(db, queue, cache),
			middleware.Host(osenv.MustGet("APP_HOST"), clevergo.PathSkipper("/assets/*", "/.well-known/*")),
			middleware.IsAuthenticated("/login", clevergo.PathSkipper(
				"/", "/login", "/assets/*", "/.well-known/*", "/api/badges/*", "/badges/*",
				"/trending", "/debug/pprof/*",
				"/signup", "/verify-email", "/send-verification-email", "/forgot-password", "/reset-password",
			)),
			clevergo.WrapHH(nosurf.NewPure),
		)
		app.Renderer = provideRenderer(sessionManager)
		app.ServeFiles("/assets", packr.New("public", "../public"))

		pprof.RegisterHandler(app)

		basicHandler := core.NewHandler(db, sessionManager, queue, cache)
		handlers := []web.Handler{
			&home.Handler{basicHandler},
			&dashboard.Handler{basicHandler},
			&user.Handler{basicHandler},
			&pkg.Handler{basicHandler},
			&domain.Handler{basicHandler},
			&report.Handler{basicHandler},
			&api.Handler{basicHandler},
			&badge.Handler{basicHandler},
			&trending.Handler{basicHandler},
			&search.Handler{basicHandler},
		}
		for _, handler := range handlers {
			handler.Register(app)
		}

		return app.Run(osenv.Get("HTTP_ADDR", ":8080"))
	},
}

func provideI18N() *i18n.Bundle {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	box := packr.New("locales", "../locales")
	for _, lang := range []string{"en", "zh-CN", "zh-TW"} {
		f, err := box.Open(lang + ".toml")
		if err != nil {
			stdlog.Fatal(err)
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			stdlog.Fatal(err)
		}
		_, err = b.ParseMessageFileBytes(buf, lang+".toml")
		if err != nil {
			stdlog.Fatal(err)
		}
	}
	return b
}

func provideRenderer(sessionManager *scs.SessionManager) clevergo.Renderer {
	box := packr.New("views", "../views")
	set := jet.NewHTMLSetLoader(jetpackr.New(box))
	set.SetDevelopmentMode(core.IsDevelopMode())
	set.AddGlobalFunc("shortScale", func(args jet.Arguments) reflect.Value {
		args.RequireNumOfArguments("shortScale", 1, 1)
		return reflect.ValueOf(stringhelper.ShortScale(args.Get(0).Int()))
	})
	set.AddGlobal("siteURL", osenv.MustGet("APP_URL"))
	jetsprig.GenericFuncMap().AttachTo(set)
	renderer := jetrenderer.New(set)
	renderer.SetBeforeRender(func(w io.Writer, name string, vars jet.VarMap, data interface{}, c *clevergo.Context) error {
		ctx := c.Context()
		localizer := i18nmiddleware.Localizer(c)
		vars.SetFunc("i18n", func(args jet.Arguments) reflect.Value {
			args.RequireNumOfArguments("i18n", 1, 2)
			msgID := args.Get(0).String()
			var tmplData interface{}
			var pc interface{}
			if args.NumOfArguments() > 1 {
				data := args.Get(1)
				switch data.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					reflect.Float32, reflect.Float64:
					pc = data.Interface()
					tmplData = clevergo.Map{
						"Count": pc,
					}
				default:
					tmplData = data.Interface()
				}
			}
			v, _, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
				MessageID:    msgID,
				TemplateData: tmplData,
				PluralCount:  pc,
			})
			if err != nil {
				c.Logger().Warnf("failed to translate %s: %s", msgID, err)
			}
			return reflect.ValueOf(v)
		})
		vars.Set("user", authmiddleware.GetIdentity(ctx))
		vars.Set("csrf", nosurf.Token(c.Request))
		vars.Set("alert", sessionManager.Pop(ctx, "alert"))
		vars.SetFunc("date", func(args jet.Arguments) reflect.Value {
			args.RequireNumOfArguments("date", 1, 1)
			date := args.Get(0).Interface().(time.Time)
			return reflect.ValueOf(date.Format(osenv.Get("DATE_FORMAT", "2006-01-02 15:04:05")))
		})
		return nil
	})
	return renderer
}

func provideSessionManager() *scs.SessionManager {
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			db, _ := strconv.Atoi(osenv.Get("REDIS_DATABASE", "0"))
			opts := []redis.DialOption{redis.DialDatabase(db)}
			if password := osenv.Get("REDIS_PASSWORD"); password != "" {
				opts = append(opts, redis.DialPassword(password))
			}
			return redis.Dial("tcp", osenv.Get("REDIS_ADDR", "localhost:6379"), opts...)
		},
	}

	_, err := pool.Get().Do("PING")
	if err != nil {
		panic(err)
	}

	gob.Register(map[string]interface{}{})
	m := scs.New()
	m.Store = redisstore.New(pool)
	m.Lifetime = 24 * time.Hour
	m.Cookie.HttpOnly = false
	return m
}

func provideLogger() (log.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	if logFile := osenv.Get("LOG_FILE"); logFile != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, logFile)
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func provideCache() *ristretto.Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 27, // maximum cost of cache (128M).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return cache
}
