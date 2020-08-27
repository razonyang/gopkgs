package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"clevergo.tech/clevergo"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/RichardKnop/machinery/v2"
	"github.com/dgraph-io/ristretto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
		if !c.IsGet() || c.QueryParam("go-get") != "1" {
			return next(c)
		}

		ctx := c.Context()
		pkg, err := gg.getPackage(ctx, c.Host(), c.Request.URL.Path[1:])
		if err != nil {
			return err
		}

		if c.QueryParam("preview") == "" {
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
