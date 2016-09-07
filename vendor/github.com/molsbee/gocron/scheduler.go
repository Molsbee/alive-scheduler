package gocron

import "time"

// Scheduler contains multiple tasks for execution and manages those tasks
// execution
type Scheduler interface {
	ScheduleSimpleTask() *SimpleTask
	ScheduleCronTask(schedule string, fn func())
	Start() chan bool
}

type scheduler struct {
	tasks       []Task
	workChannel chan Task
}

// NewScheduler creates a new instance of a Scheduler which consists of multiple
// task for execution.  Along with a pool of workers defined for processing those
// tasks.  The pool is defined through this constractor and it's recomended to
// minimize the amount of time intensive tasks as they will cause the workers
// to delay execution.  Size tasks and worker pool appropriately
func NewScheduler(poolSize int) Scheduler {
	scheduler := &scheduler{
		tasks:       []Task{},
		workChannel: make(chan Task, 10),
	}

	for i := 0; i <= poolSize; i++ {
		go scheduler.processTask()
	}

	return scheduler
}

// ScheduleSimpleTask is used as a builder pattern to return a task that can then be
// added to the Scheduler through a fluent style api
func (s *scheduler) ScheduleSimpleTask() *SimpleTask {
	task := NewSimpleTask()
	s.tasks = append(s.tasks, task)
	return task
}

func (s *scheduler) ScheduleCronTask(cron string, fn func()) {
	task := NewCronTask(cron, fn)
	s.tasks = append(s.tasks, task)
}

// Start must be called as it starts the polling process for task execution
// returns a channel which you can push a message onto to stop the scheduler
func (s *scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.runPending()
			case <-stopped:
				return
			}
		}
	}()

	return stopped
}

// Iterates through all tasks and sends task that need running to work channel
func (s *scheduler) runPending() {
	for i := 0; i < len(s.tasks); i++ {
		if s.tasks[i].shouldRun() {
			s.workChannel <- s.tasks[i]
		}
	}
}

// Processes task passed over the work channel and calls run on the task
func (s *scheduler) processTask() {
	for task := range s.workChannel {
		task.run()
	}
}
