package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/duzgunberke/task-queue/api"
	"github.com/duzgunberke/task-queue/tasks"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		go api.StartPrometheusMetricsServer()
		http.Handle("/metrics", promhttp.Handler())
		server.ListenAndServe()
	}()

	fmt.Println("Server is running on :8080")
	
	// Uygulamanın bir süre çalışmasını simüle et
	time.Sleep(20 * time.Second)
}
