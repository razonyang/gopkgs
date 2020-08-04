package cmd

import (
	"github.com/urfave/cli/v2"
	"pkg.razonyang.com/gopkgs/internal/tasks"
)

func init() {
	app.Commands = append(app.Commands, queueCmd)
}

var queueCmd = &cli.Command{
	Name:  "queue",
	Usage: "start a tasks queue",
	Action: func(c *cli.Context) error {
		packageTask := tasks.NewPackage(db)
		err := queue.RegisterTasks(map[string]interface{}{
			"package.action": packageTask.Action,
			"sendMail":       tasks.SendMail,
		})
		if err != nil {
			return err
		}

		worker := queue.NewWorker("gopkgs_worker", 0)
		return worker.Launch()
	},
}
