package service

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/molsbee/alive-common/model/database"
	"github.com/molsbee/alive-common/rabbitmq"
	"github.com/molsbee/alive-common/repository"
)

// HTTPService contains all the required dependencies to get and distribute
// work from http_get_config table
type HTTPService struct {
	repo        *repository.HTTPGetConfig
	queue       rabbitmq.RabbitQueue
	workChannel chan database.HTTPGetConfig
}

// NewHTTPService is a construct for creating a HTTPService struct
func NewHTTPService(db *gorm.DB, queue rabbitmq.RabbitQueue) *HTTPService {
	service := HTTPService{
		repo:        repository.NewHTTPGetConfig(db),
		queue:       queue,
		workChannel: make(chan database.HTTPGetConfig, 20),
	}

	// startup multiple dispatchers in seperate goroutines
	for i := 0; i <= 5; i++ {
		go service.dispatcher()
	}

	return &service
}

// DispatchHTTPGetWork polls database table http_get_config for all rows and issues
// that work to RabbitMQ queue where multiple processes are listening for work.
// Mechanism in main will use this method on a cron intervale to distribute work.
func (s *HTTPService) DispatchHTTPGetWork() {
	configs, err := s.repo.FindAll()
	if err != nil {
		log.Printf("Error: %v", err)
	}

	for _, config := range configs {
		s.workChannel <- config
	}
}

func (s *HTTPService) dispatcher() {
	for w := range s.workChannel {
		err := s.queue.Publish(w)
		if err != nil {
			log.Printf("Failed to publish work to queue %v", err)
		}
	}
}
