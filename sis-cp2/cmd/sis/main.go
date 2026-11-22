package main

import (
	"fmt"
	"sis-cp2/internal/analyzer"
	"sis-cp2/internal/app"
	"sis-cp2/internal/storage"
)

func main() {
	repo := storage.NewMemoryRepo()
	proc := &analyzer.TextProcessor{}
	application := app.NewApp(repo, proc)
	
	// Start worker pool: 3 workers, queue capacity of 5
	application.Start(3, 5)
	fmt.Println("Worker pool started.")

	// Create and submit multiple tasks
	tasks := []string{"Hello Golang", "Concurrency in Go", "Wait Groups and Channels", "SIS Checkpoint 2"}
	var ids []string
	
	for _, text := range tasks {
		id, err := application.CreateTask(text)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
		
		err = application.SubmitTask(id)
		if err != nil {
			fmt.Printf("Failed to submit task %s: %v\n", id, err)
		} else {
			fmt.Printf("Task %s submitted.\n", id)
		}
	}
	
	// Graceful shutdown: wait for all accepted tasks to finish
	fmt.Println("Initiating graceful shutdown...")
	application.Stop()
	
	// Verify results
	fmt.Println("All workers finished. Results:")
	for _, id := range ids {
		task, _ := application.GetTask(id)
		fmt.Printf("- %s [%s]: Words: %d, Chars: %d\n", task.ID, task.State, task.Result.WordCount, task.Result.CharCount)
	}
}
