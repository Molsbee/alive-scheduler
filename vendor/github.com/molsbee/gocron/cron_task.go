package gocron

import (
	"time"

	"github.com/molsbee/gocron/parser"
)

// CronTask struct defining an operation that needs to be executed on a schedule
type CronTask struct {
	task     Runnable
	schedule *parser.Schedule
	nextRun  time.Time
}

// NewCronTask constructor for creating a new CronTask from a cron spec string
func NewCronTask(cron string, fn func()) Task {
	return &CronTask{
		task:     RunnableFunc(fn),
		schedule: parser.Parse(cron),
	}
}

func (c *CronTask) run() {
	c.task.Run()
	c.scheduleNextExecution()
}

func (c *CronTask) shouldRun() bool {
	if (c.nextRun == time.Time{}) {
		c.nextRun = c.schedule.Next(time.Now())
	}

	return time.Now().After(c.nextRun)
}

func (c *CronTask) scheduleNextExecution() {
	c.nextRun = c.schedule.Next(time.Now())
}
