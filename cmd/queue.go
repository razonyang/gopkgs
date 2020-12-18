package cmd

import (
	"pkg.razonyang.com/gopkgs/internal/tasks"
)

func startQueue() error {
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
}
