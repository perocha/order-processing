package main

import (
	"log"

	"github.com/perocha/order-processing/config"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/repository/eventhub"
	"github.com/perocha/order-processing/pkg/usecase"
)

func main() {
	// Initialize configuration
	cfg := config.InitializeConfig()
	if cfg == nil {
		log.Println("Error: Failed to load configuration")
		panic("Failed to load configuration")
	}

	// Initialize
	createOrder := usecase.CreateOrder{}
	deleteOrder := usecase.DeleteOrder{}
	updateOrder := usecase.UpdateOrder{}
	eventConsumer := usecase.NewEventConsumer(createOrder, deleteOrder, updateOrder)
	log.Printf("Event consumer initialized %v", eventConsumer)

	// Initialize the event hub adapter
	eventHubAdapter, cleanup, err := eventhub.ConsumerInit(cfg.EventHubConnectionString, cfg.EventHubName, cfg.CheckpointStoreContainerName, cfg.CheckpointStoreConnectionString, eventConsumer)
	if err != nil {
		log.Println("Error: Failed to initialize event hub adapter")
		panic("Failed to initialize event hub adapter")
	}
	log.Printf("Event hub adapter initialized %v", eventHubAdapter)

	// Cleanup
	defer cleanup()
}
