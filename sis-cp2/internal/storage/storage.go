package storage

import (
	"sync"
	"sis-cp2/internal/domain"
)

type Repository interface {
	Save(task domain.Task) error
	Get(id string) (domain.Task, error)
	Update(task domain.Task) error
}

type MemoryRepo struct {
	mu    sync.RWMutex
	tasks map[string]domain.Task
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{tasks: make(map[string]domain.Task)}
}

func (r *MemoryRepo) Save(task domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryRepo) Get(id string) (domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists {
		return domain.Task{}, domain.ErrTaskNotFound
	}
	return task, nil
}

func (r *MemoryRepo) Update(task domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; !exists {
		return domain.ErrTaskNotFound
	}
	r.tasks[task.ID] = task
	return nil
}
