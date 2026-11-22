package app

import (
	"fmt"
	"sync"
	"sis-cp2/internal/analyzer"
	"sis-cp2/internal/domain"
	"sis-cp2/internal/storage"
)

type App struct {
	repo      storage.Repository
	processor analyzer.Processor
	
	nextID int
	idMu   sync.Mutex

	taskQueue chan string
	wg        sync.WaitGroup
	stopped   bool
	stopMu    sync.RWMutex
}

func NewApp(repo storage.Repository, processor analyzer.Processor) *App {
	return &App{
		repo:      repo,
		processor: processor,
		nextID:    1,
	}
}

func (a *App) Start(workers int, queueSize int) {
	a.taskQueue = make(chan string, queueSize)
	a.stopped = false

	for i := 0; i < workers; i++ {
		a.wg.Add(1)
		go a.worker()
	}
}

func (a *App) Stop() {
	a.stopMu.Lock()
	if a.stopped {
		a.stopMu.Unlock()
		return
	}
	a.stopped = true
	a.stopMu.Unlock()

	close(a.taskQueue)
	a.wg.Wait()
}

func (a *App) CreateTask(input string) (string, error) {
	if input == "" {
		return "", domain.ErrInvalidInput
	}
	
	a.idMu.Lock()
	id := fmt.Sprintf("task-%d", a.nextID)
	a.nextID++
	a.idMu.Unlock()
	
	task := domain.Task{
		ID:    id,
		Input: input,
		State: domain.StatePending,
	}
	
	err := a.repo.Save(task)
	return id, err
}

func (a *App) GetTask(id string) (domain.Task, error) {
	return a.repo.Get(id)
}

func (a *App) SubmitTask(id string) error {
	a.stopMu.RLock()
	defer a.stopMu.RUnlock()
	
	if a.stopped {
		return domain.ErrAppStopped
	}

	task, err := a.repo.Get(id)
	if err != nil {
		return err
	}
	
	if task.State != domain.StatePending {
		return domain.ErrInvalidTransition
	}
	
	task.State = domain.StateQueued
	err = a.repo.Update(task)
	if err != nil {
		return err
	}
	
	select {
	case a.taskQueue <- id:
		return nil
	default:
		task.State = domain.StatePending // revert state if queue is full
		a.repo.Update(task)
		return domain.ErrQueueFull
	}
}

func (a *App) worker() {
	defer a.wg.Done()
	for id := range a.taskQueue {
		a.processTask(id)
	}
}

func (a *App) processTask(id string) {
	task, err := a.repo.Get(id)
	if err != nil {
		return
	}
	
	if task.State != domain.StateQueued {
		return 
	}
	
	task.State = domain.StateProcessing
	a.repo.Update(task)
	
	result, err := a.processor.Process(task.Input)
	if err == nil {
		task.Result = &result
		task.State = domain.StateCompleted
	} else {
		// keeping it simple, revert to pending or error state (not explicitly required)
		task.State = domain.StatePending 
	}
	a.repo.Update(task)
}
