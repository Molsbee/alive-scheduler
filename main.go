package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/molsbee/alive-common/rabbitmq"
	"github.com/molsbee/alive-scheduler/resource"
	"github.com/molsbee/alive-scheduler/service"
	"github.com/molsbee/gocron"
)

var databaseURL string
var amqpURL string
var hostName string

func init() {
	flag.StringVar(&databaseURL, "database", "root@tcp(localhost:3306)/alive?parseTime=true", "-database=root@tcp(localhost:3306)/alive?parseTime=true")
	flag.StringVar(&amqpURL, "amqp", "amqp://guest:guest@localhost:5672", "-amqp=amqp://guest:guest@localhost:5672")

	var err error
	if hostName, err = os.Hostname(); err != nil {
		log.Fatalf("Failed to retrieve hostName used as ConsumerTag for RabbitMQ: %v", err)
	}
}

func main() {
	log.Printf("Establishing database connection - %s\n", databaseURL)
	db, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		log.Fatalf("Unable to open database connection - Error: %v\n", err)
	}
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(20)
	db.LogMode(true)

	log.Printf("Establishing connection to rabbitmq - %s\n", amqpURL)
	rabbitmq, err := rabbitmq.NewRabbitMQ(amqpURL, hostName)
	if err != nil {
		log.Fatalf("Unable to setup rabbitmq - Error: %v\n", err)
	}

	queue, err := rabbitmq.NewRabbitQueue("alive.http.get", "alive", "http.get", nil)
	if err != nil {
		log.Printf("Failed to setup exchange and queue %v", err)
	}

	httpService := service.NewHTTPService(db, queue)

	scheduler := gocron.NewScheduler(1)
	scheduler.ScheduleSimpleTask().Minute(1).Run(httpService.DispatchHTTPGetWork)
	scheduler.Start()

	healthResource := resource.NewHealthResource(db, queue)

	router := mux.NewRouter()
	router.HandleFunc("/health", healthResource.Get)

	http.ListenAndServe(":8080", router)
}
