package main

import (
	"log"

	"github.com/perocha/order-processing/config"
	"github.com/perocha/order-processing/pkg/infrastructure/service/eventhub"
)

func main() {
	cfg := config.LoadConfig()
	eventHubService := eventhub.NewService(cfg)
	log.Println("EventHub service started", eventHubService)

	// Initialize other services and start your application

}
