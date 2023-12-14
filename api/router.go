package api

import (
	"github.com/gorilla/mux"
)

func setupAPIRoutes(taskQueue *TaskQueue) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/enqueue", EnqueueTaskHandler).Methods("POST")
	r.HandleFunc("/tasks", GetTasksHandler).Methods("GET")

	return r
}
