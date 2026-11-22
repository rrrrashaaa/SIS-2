package domain

import "errors"

type TaskState string

const (
	StatePending    TaskState = "Pending"
	StateQueued     TaskState = "Queued"
	StateProcessing TaskState = "Processing"
	StateCompleted  TaskState = "Completed"
)

var (
	ErrInvalidInput      = errors.New("invalid task input")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrQueueFull         = errors.New("task queue is full")
	ErrAppStopped        = errors.New("application is stopped, cannot accept new tasks")
)

type TaskResult struct {
	WordCount int
	CharCount int
}

type Task struct {
	ID     string
	Input  string
	State  TaskState
	Result *TaskResult
}
