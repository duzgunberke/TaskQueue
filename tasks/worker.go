package task

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	ID    int
	Tasks chan Task
}

func (w *Worker) ProcessTasks(wg *sync.WaitGroup) {
	defer wg.Done()
	processedTasks := make(map[int]bool)
	for task := range w.Tasks {
		if processedTasks[task.ID] {
			duplicateTaskCounter.Inc()
			continue
		}
		processedTasks[task.ID] = true

		startTime := time.Now()
		success := w.executeTask(task)
		duration := time.Since(startTime).Seconds()

		if success {
			taskCounter.WithLabelValues("success").Inc()
			taskDuration.WithLabelValues("success").Observe(duration)
			priorityCounter.WithLabelValues(fmt.Sprintf("%d", task.Priority)).Inc()
		} else {
			taskCounter.WithLabelValues("failure").Inc()
			taskDuration.WithLabelValues("failure").Observe(duration)

			w.handleTaskFailure(task)
		}
	}
}

func (w *Worker) executeTask(task Task) bool {
	fmt.Printf("Worker %d processing task %d with payload: %s, scheduled at: %s, priority: %d\n", w.ID, task.ID, task.Payload, task.Schedule.Format(time.RFC3339), task.Priority)
	// Simulate task processing
	time.Sleep(2 * time.Second)

	// Simulate a task failure
	if task.ID == 3 || task.ID == 4 {
		return false
	}

	return true
}

func (w *Worker) handleTaskFailure(task Task) {
	if task.MaxRetries > 0 {
		task.MaxRetries--
		fmt.Printf("Task %d failed. Retrying. Remaining retries: %d\n", task.ID, task.MaxRetries)
		w.Tasks <- task
	} else {
		fmt.Printf("Task %d failed. No more retries.\n", task.ID)
	}
}
