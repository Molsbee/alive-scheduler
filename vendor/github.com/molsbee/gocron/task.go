package gocron

// Task - interface for defining generic work contract used by scheduler
type Task interface {
	run()
	shouldRun() bool
}

// Runnable contract for a function that is runnable by cron scheduler
type Runnable interface {
	Run()
}

// RunnableFunc type is an adapter to allow the use of a ordinary functions
// as argument to cron Run.
type RunnableFunc func()

// Run method executes function with no arguments
func (r RunnableFunc) Run() {
	r()
}
