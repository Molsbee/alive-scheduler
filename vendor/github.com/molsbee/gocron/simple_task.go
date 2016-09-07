package gocron

import (
	"log"
	"time"
)

// SimpleTask struct defining an operation that needs to be execute by gocron
type SimpleTask struct {
	task   Runnable
	hour   int
	minute int
	second int

	lastRun time.Time
	nextRun time.Time
	period  time.Duration
}

// NewSimpleTask constructor for creating a unpopulated SimpleTask struct
func NewSimpleTask() *SimpleTask {
	return &SimpleTask{}
}

// Hour every x number of hours the cron task will be executed
func (t *SimpleTask) Hour(h int) *SimpleTask {
	if h < 0 || h > 25 {
		log.Panic("Hour time interval outside of bounds [0 - 24]")
	}
	t.hour = h
	return t
}

// Minute every x number of minutes the cron task will be executed
func (t *SimpleTask) Minute(m int) *SimpleTask {
	if m < 0 || m > 60 {
		log.Panic("Minute time interval outside of bounds [0 - 60]")
	}
	t.minute = m
	return t
}

// Second every x number of seconds the cron task will be executed
func (t *SimpleTask) Second(s int) *SimpleTask {
	if s < 0 || s > 60 {
		log.Panic("Second time interval outside of bounds [0-60]")
	}
	t.second = s
	return t
}

// Run specifies the function that should be called every time the task Scheduler
// executes.  Currently the only supported function will have to implement the
// runnable contract which is a no argument function with no return value.
// Most tasks being ran on cron won't have anyone to parse results so no return
// value is expected.  Also a function that is runnable should be self contained
// at this point should alow avoiding reflection
func (t *SimpleTask) Run(fn func()) {
	t.task = RunnableFunc(fn)
	t.scheduleNextExecution()
}

func (t *SimpleTask) run() {
	t.task.Run()
	t.lastRun = time.Now()
	t.scheduleNextExecution()
}

// ShouldRun used to deteremine if a task is ready for execution
func (t *SimpleTask) shouldRun() bool {
	return time.Now().After(t.nextRun)
}

func (t *SimpleTask) scheduleNextExecution() {
	if t.lastRun == time.Unix(0, 0) {
		t.lastRun = time.Now()
	}

	if t.period == 0 {
		second := time.Duration(t.second)
		minute := time.Duration(t.minute * 60)
		hour := time.Duration(t.hour * 60 * 60)
		t.period = second + minute + hour
	}
	t.nextRun = t.lastRun.Add(t.period * time.Second)
}
