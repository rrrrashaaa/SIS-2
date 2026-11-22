# SIS Checkpoint 2 - Concurrent Text Analyzer

This is an extension of Checkpoint 1, introducing concurrent task processing with a worker pool.

## Concurrency Architecture
- **Worker Pool:** The application uses a configurable number of worker goroutines (default 3) to process tasks concurrently.
- **Task Queue:** A bounded channel acts as the queue. If the queue is full, new submissions are rejected with an error.
- **Synchronization:** `sync.RWMutex` is used in the memory repository to ensure thread-safe access to task data. `sync.Mutex` protects task ID generation.
- **Shutdown:** Graceful shutdown is supported. The application stops accepting new tasks, closes the queue, and waits for all workers to finish their current tasks using `sync.WaitGroup`.

## Task Lifecycle & Allowed Transitions
1. **Pending**: Created and stored.
2. **Queued**: Submitted for processing and waiting in the bounded queue.
3. **Processing**: Picked up by a worker.
4. **Completed**: Analysis finished.
