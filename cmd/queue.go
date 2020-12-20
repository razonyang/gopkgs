package cmd

import (
	"log"

	"pkg.razonyang.com/gopkgs/internal/tasks"
)

func startQueue() {
	packageTask := tasks.NewPackage(db)
	err := queue.RegisterTasks(map[string]interface{}{
		"package.action": packageTask.Action,
		"sendMail":       tasks.SendMail,
	})
	if err != nil {
		log.Fatal(err)
	}

	worker := queue.NewWorker("gopkgs_worker", 0)
	log.Fatal(worker.Launch())
}
