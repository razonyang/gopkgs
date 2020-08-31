package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"clevergo.tech/clevergo"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/RichardKnop/machinery/v2"
	"github.com/dgraph-io/ristretto"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/core"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

type goGet struct {
	db    *sqlx.DB
	queue *machinery.Server
	cache *ristretto.Cache
}

func (gg *goGet) getCacheKey(host, path string) string {
	return fmt.Sprintf("pkg:%s/%s", host, path)
}

func (gg *goGet) getPackage(ctx context.Context, host, path string) (*models.Package, error) {
	pkg, err := gg.getPackageFromCache(host, path)
	if err == nil {
		return pkg, nil
	}

	pkg, err = gg.getPackageFromDB(ctx, host, path)
	if err != nil {
		return nil, err
	}

	// cache
	if !gg.cache.SetWithTTL(gg.getCacheKey(host, path), pkg, 0, 5*time.Minute) {
		log.Println("unable to cache package")
	}

	return pkg, nil
}

func (gg *goGet) getPackageFromCache(host, path string) (*models.Package, error) {
	key := gg.getCacheKey(host, path)
	value, found := gg.cache.Get(key)
	if found {
		if pkg, ok := value.(*models.Package); ok {
			return pkg, nil
		}
	}
	return nil, fmt.Errorf("cache %s no found", key)
}

func (gg *goGet) getPackageFromDB(ctx context.Context, host, path string) (*models.Package, error) {
	var pkg models.Package
	err := models.FindPackageByDomainAndPath(ctx, gg.db, &pkg, host, path)
	return &pkg, err
}

func (gg *goGet) middleware(next clevergo.Handle) clevergo.Handle {
	return func(c *clevergo.Context) error {
		if !c.IsGet() {
			return next(c)
		}

		var host, path string
		goGet := c.QueryParam("go-get") == "1"
		if goGet {
			host = core.GetHost(c)
			path = c.Request.URL.Path[1:]
		} else {
			parts := strings.Split(c.Request.URL.Path, "/")
			if len(parts) < 3 {
				return next(c)
			}
			if err := validation.Validate(parts[1], validation.Required, is.Domain); err != nil {
				return next(c)
			}

			host = parts[1]
			path = strings.Join(parts[2:], "/")
		}

		ctx := c.Context()
		pkg, err := gg.getPackage(ctx, host, path)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.NotFound()
			}
			return err
		}

		if goGet {
			go func() {
				_, err := gg.queue.SendTaskWithContext(context.Background(), &tasks.Signature{
					UUID: uuid.New().String(),
					Name: "package.action",
					Args: []tasks.Arg{
						{Name: "kind", Type: "string", Value: models.ActionGoGet},
						{Name: "packageID", Type: "int64", Value: pkg.ID},
						{Name: "createdAt", Type: "int64", Value: time.Now().Unix()},
					},
					RetryCount: 3,
				})

				if err != nil {
					c.Logger().Errorf("failed to enqueue a task: %s", err.Error())
				}
			}()
		}

		return c.Render(http.StatusOK, "package/view.tmpl", clevergo.Map{
			"page": web.NewPage(pkg.Prefix()),
			"pkg":  pkg,
		})
	}

}

func GoGet(db *sqlx.DB, queue *machinery.Server, cache *ristretto.Cache) clevergo.MiddlewareFunc {
	m := &goGet{db, queue, cache}
	return m.middleware
}
