package queue

import (
	"fmt"

	"github.com/grafana/sobek"
)

// represents a ready-to-execute task
type Task interface {
	Run(vm *sobek.Runtime) error
}

// A FIFO queue for Tasks
type TaskQueue struct {
	queue []Task
}

func NewQueue() *TaskQueue {
	return &TaskQueue{
		queue: []Task{},
	}
}

// Whether we have anything to run right now
func (e *TaskQueue) HasAvailableTasks() bool {
	return len(e.queue) > 0
}

// DrainQueue runs all tasks in the queue until it is empty
func (e *TaskQueue) DrainQueue(vm *sobek.Runtime) {
	for _, task := range e.queue {
		if err := task.Run(vm); err != nil {
			fmt.Printf("[TaskQueue] error in task: %v\n", err)
		}
	}

	e.queue = []Task{}
}

// Adds a task to the queue. It does not run the task yet.
func (e *TaskQueue) ScheduleTask(t Task) {
	e.queue = append(e.queue, t)
}

type funcWrapperTask struct {
	f func(*sobek.Runtime) error
}

func (w funcWrapperTask) Run(vm *sobek.Runtime) error {
	return w.f(vm)
}

func (e *TaskQueue) ScheduleFunc(fn func(*sobek.Runtime) error) {
	e.ScheduleTask(funcWrapperTask{f: fn})
}
