package task

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID        int
	Payload   string
	Schedule  time.Time
	Interval  time.Duration
	Timeout   time.Duration
	Priority  int
	MaxRetries int
}

var (
	mutexMap     = make(map[string]*sync.Mutex)
	mutexMapLock sync.Mutex
)

type TaskQueue struct {
	Tasks   chan Task
	Workers []*Worker
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{Tasks: make(chan Task, 10)}
}


// EnqueueTask adds a task to the task queue.
func (q *TaskQueue) EnqueueTask(task Task) {
	// Eş zamanlı işlemi önlemek için göreve özel bir kilit oluştur
	taskLock := getTaskLock(task.ID)
	taskLock.Lock()
	defer taskLock.Unlock()

	q.Tasks <- task
}

func (q *TaskQueue) StartWorkers(numWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		worker := &Worker{ID: i, Tasks: make(chan Task, 10)}
		q.Workers = append(q.Workers, worker)
		wg.Add(1)
		go worker.ProcessTasks(&wg)
		go q.DispatchTasks(worker.Tasks)
	}
	go func() {
		wg.Wait()
		close(q.Tasks)
		for _, worker := range q.Workers {
			close(worker.Tasks)
		}
	}()
}

func (q *TaskQueue) DispatchTasks(workerTasks chan Task) {
	for task := range q.Tasks {
		if task.Schedule.IsZero() || time.Now().After(task.Schedule) {
			q.handleTaskTimeout(task)
			workerTasks <- task
		} else {
			time.Sleep(task.Schedule.Sub(time.Now()))
			q.handleTaskTimeout(task)
			workerTasks <- task
		}

		if task.Interval > 0 {
			go func(t Task) {
				timer := time.NewTicker(t.Interval)
				defer timer.Stop()
				for {
					select {
					case <-timer.C:
						q.handleTaskTimeout(t)
						workerTasks <- t
					}
				}
			}(task)
		}
	}
	close(workerTasks)
}

func (q *TaskQueue) handleTaskTimeout(task Task) {
	if task.Timeout > 0 {
		go func(t Task) {
			select {
			case <-time.After(t.Timeout):
				fmt.Printf("Task %d timed out. Retrying.\n", t.ID)
				q.EnqueueTask(t)
			}
		}(task)
	}
}

// Görev önceliklerine göre kilit oluştur
func getTaskLock(taskID int) *sync.Mutex {
	key := fmt.Sprintf("task_%d", taskID)

	mutexMapLock.Lock()
	defer mutexMapLock.Unlock()

	mutex, ok := mutexMap[key]
	if !ok {
		mutex = &sync.Mutex{}
		mutexMap[key] = mutex
	}

	return mutex
}
