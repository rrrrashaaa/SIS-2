package app

import (
	"sync"
	"testing"
	"sis-cp2/internal/analyzer"
	"sis-cp2/internal/domain"
	"sis-cp2/internal/storage"
)

func setupApp() *App {
	repo := storage.NewMemoryRepo()
	proc := &analyzer.TextProcessor{}
	return NewApp(repo, proc)
}

func TestConcurrentProcessing(t *testing.T) {
	app := setupApp()
	app.Start(3, 10)
	defer app.Stop()

	var wg sync.WaitGroup
	var ids []string
	
	// Create tasks
	for i := 0; i < 5; i++ {
		id, _ := app.CreateTask("hello concurrent world")
		ids = append(ids, id)
	}

	// Submit concurrently
	for _, id := range ids {
		wg.Add(1)
		go func(taskID string) {
			defer wg.Done()
			err := app.SubmitTask(taskID)
			if err != nil {
				t.Errorf("failed to submit: %v", err)
			}
		}(id)
	}
	wg.Wait()

	app.Stop() // wait for processing to finish

	for _, id := range ids {
		task, _ := app.GetTask(id)
		if task.State != domain.StateCompleted {
			t.Errorf("task %s state is %s, want %s", id, task.State, domain.StateCompleted)
		}
		if task.Result == nil || task.Result.WordCount != 3 {
			t.Errorf("task result incorrect")
		}
	}
}

func TestQueueFull(t *testing.T) {
	app := setupApp()
	app.Start(1, 1) // 1 worker, capacity 1
	
	// Stop worker from picking up fast enough (app.Stop will wait, so we test differently)
	// We will fill the queue up quickly.
	app.stopMu.Lock()
	app.stopped = false
	app.stopMu.Unlock()
	
	id1, _ := app.CreateTask("task 1")
	id2, _ := app.CreateTask("task 2")
	id3, _ := app.CreateTask("task 3")
	
	app.SubmitTask(id1)
	app.SubmitTask(id2)
	err := app.SubmitTask(id3)
	
	// Due to race conditions, it might be 2 or 3 that fails, but at least one should fail if we submit fast enough
	// Wait, actually because worker picks up immediately, let's just test that the channel can be full
	// if we don't start the app.
	
	app2 := setupApp()
	app2.taskQueue = make(chan string, 1)
	
	id, _ := app2.CreateTask("dummy")
	app2.SubmitTask(id) // fills the queue
	
	idNew, _ := app2.CreateTask("dummy2")
	errNew := app2.SubmitTask(idNew) // should reject
	if errNew != domain.ErrQueueFull {
		t.Errorf("expected ErrQueueFull, got %v", errNew)
	}
}

func TestShutdownRejectsNewTasks(t *testing.T) {
	app := setupApp()
	app.Start(1, 5)
	app.Stop()
	
	id, _ := app.CreateTask("test")
	err := app.SubmitTask(id)
	if err != domain.ErrAppStopped {
		t.Errorf("expected ErrAppStopped, got %v", err)
	}
}
