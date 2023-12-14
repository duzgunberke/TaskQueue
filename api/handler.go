package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func EnqueueTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	taskQueue.EnqueueTask(task)
	w.WriteHeader(http.StatusAccepted)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Bu örnek basit bir şekilde tüm görevleri döndürüyor, gerçek bir uygulama daha karmaşık bir sorgu ve filtreleme yapacaktır.
	var tasks []Task
	for i := 0; i < 5; i++ {
		tasks = append(tasks, Task{
			ID:        i + 1,
			Payload:   fmt.Sprintf("Task %d", i+1),
			Schedule:  time.Now(),
			Interval:  0,
			Timeout:   0,
			MaxRetries: 3,
			Priority:  rand.Intn(5) + 1, // 1 ile 5 arasında rastgele bir öncelik
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
