package resource

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/molsbee/alive-common/rabbitmq"
	"github.com/molsbee/alive-scheduler/model"
)

var (
	up   = model.Status{Status: "UP"}
	down = model.Status{Status: "DOWN"}
)

// HealthResource contract for web endpoints
type HealthResource interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type health struct {
	db    *gorm.DB
	queue rabbitmq.RabbitQueue
}

// NewHealthResource constructor for initializing health resource
func NewHealthResource(db *gorm.DB, queue rabbitmq.RabbitQueue) HealthResource {
	return &health{db: db, queue: queue}
}

// Get http HandlerFunc for returning the current health of the application
func (h *health) Get(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	response := model.HealthResponse{
		Status:   up,
		Database: getDatabaseHealth(h.db),
		RabbitMQ: getQueueHealth(h.queue),
	}

	if response.Database == down || response.RabbitMQ.Status == down {
		statusCode = http.StatusServiceUnavailable
		response.Status = down
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func getDatabaseHealth(db *gorm.DB) model.Status {
	err := db.Exec("show tables").Error
	if err != nil {
		return down
	}
	return up
}

func getQueueHealth(queue rabbitmq.RabbitQueue) model.QueueHealth {
	health := model.QueueHealth{
		Status:    up,
		QueueName: queue.QueueName(),
	}

	q, err := queue.Inspect()
	if err != nil {
		health.Status = down
		return health
	}

	health.QueueDepth = q.Messages
	health.ConsumerCount = q.Consumers
	return health
}
