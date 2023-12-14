package main

import (
	"net/http"
	"time"

	"github.com/duzgunberke/task-queue/api"
	"github.com/duzgunberke/task-queue/tasks"
	"github.com/prometheus/client_golang/prometheus/promhttp" // eksik import

)

func main() {
	taskQueue := tasks.NewTaskQueue() // tasks paketini import etmiş olmalısınız
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
