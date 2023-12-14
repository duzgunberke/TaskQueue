package api

import (
	"github.com/gorilla/mux"
	"github.com/duzgunberke/task-queue/tasks" // Eksik olan import
)

func SetupAPIRoutes(taskQueue *tasks.TaskQueue) *mux.Router { // TaskQueue tipi düzeltildi
	r := mux.NewRouter()

	r.HandleFunc("/enqueue", EnqueueTaskHandler).Methods("POST")
	r.HandleFunc("/tasks", GetTasksHandler).Methods("GET")

	return r
}
