package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	taskCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_queue_processed_total",
			Help: "Total number of processed tasks",
		},
		[]string{"status"},
	)

	taskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "task_queue_duration_seconds",
			Help:    "Histogram of task processing durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	priorityCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_queue_priority_total",
			Help: "Total number of tasks processed by priority",
		},
		[]string{"priority"},
	)

	duplicateTaskCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "task_queue_duplicate_total",
			Help: "Total number of duplicate tasks received",
		},
	)

	mutexMap     = make(map[string]*sync.Mutex)
	mutexMapLock sync.Mutex
)

func init() {
	prometheus.MustRegister(taskCounter)
	prometheus.MustRegister(taskDuration)
	prometheus.MustRegister(priorityCounter)
	prometheus.MustRegister(duplicateTaskCounter)
}

