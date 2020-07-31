package cmd

import (
	"github.com/robfig/cron/v3"
	"pkg.razonyang.com/gopkgs/internal/jobs"
)

func startCrond() {
	crond := cron.New()
	calendarJob := jobs.NewCalendar(db)
	go func() {
		calendarJob.Run()
	}()
	crond.AddJob("0 0 * * *", calendarJob)
	crond.Start()
}
