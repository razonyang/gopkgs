package middleware

import (
	"context"
	"net/http"

	"clevergo.tech/clevergo"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/RichardKnop/machinery/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/models"
	"pkg.razonyang.com/gopkgs/internal/web"
)

func GoGet(db *sqlx.DB, queue *machinery.Server) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(c *clevergo.Context) error {
			if !c.IsGet() || c.QueryParam("go-get") != "1" {
				return next(c)
			}

			var pkg models.Package
			ctx := c.Context()
			if err := models.FindPackageByDomainAndPath(ctx, db, &pkg, c.Host(), c.Request.URL.Path[1:]); err != nil {
				return err
			}

			go func() {
				_, err := queue.SendTaskWithContext(context.Background(), &tasks.Signature{
					UUID: uuid.New().String(),
					Name: "package.action",
					Args: []tasks.Arg{
						{Name: "kind", Type: "string", Value: models.ActionGoGet},
						{Name: "packageID", Type: "int64", Value: pkg.ID},
					},
					RetryCount: 3,
				})

				if err != nil {
					c.Logger().Errorf("failed to enqueue a task: %s", err.Error())
				}
			}()

			return c.Render(http.StatusOK, "package/view.tmpl", clevergo.Map{
				"page": web.NewPage(pkg.Prefix()),
				"pkg":  pkg,
			})
		}
	}
}
