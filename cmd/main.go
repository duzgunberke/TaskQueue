package main

import (
	"net/http"
	"time"

	"github.com/duzgunberke/task-queue/api"
	"github.com/duzgunberke/task-queue/internal/task"
)

func main() {
	taskQueue := task.NewTaskQueue()
	taskQueue.StartWorkers(3)

	apiRouter := api.SetupAPIRoutes(taskQueue)

	// HTTP Sunucusu
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.StripPrefix("/api", apiRouter),
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		server.ListenAndServe()
	}()

	time.Sleep(20 * time.Second) // Simulate the application running for a while
}
