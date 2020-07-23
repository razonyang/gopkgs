package cmd

import (
	"github.com/razonyang/gopkgs/internal/jobs"
	"github.com/robfig/cron/v3"
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
