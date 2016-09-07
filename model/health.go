package model

// Status is a shared struct between all Health items composed of a string
// Status ["UP","DOWN"]
type Status struct {
	Status string `json:"status"`
}

// QueueHealth represents the state of a RabbitMQ queue
type QueueHealth struct {
	Status
	QueueName     string
	QueueDepth    int
	ConsumerCount int
}

// HealthResponse represents the overall health of an application
type HealthResponse struct {
	Status
	Database Status      `json:"database"`
	RabbitMQ QueueHealth `json:"rabbitMQ"`
}
