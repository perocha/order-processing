package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/perocha/order-processing/config"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/repository/eventhub"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/repository/storage"
	"github.com/perocha/order-processing/pkg/usecase"
)

func main() {
	// Initialize configuration
	cfg := config.InitializeConfig()
	if cfg == nil {
		log.Println("Error: Failed to load configuration")
		panic("Failed to load configuration")
	}

	// Initialize App Insights

	// Initialize CosmosDB repository
	OrderRepository, err := storage.NewCosmosDBOrderRepository(cfg.CosmosDBConnectionString)
	if err != nil {
		log.Println("Error: Failed to initialize CosmosDB repository")
		panic("Failed to initialize CosmosDB repository")
	}
	log.Printf("CosmosDB repository initialized %v\n", OrderRepository)

	// Initialize order processing
	orderProcessing := usecase.NewOrderProcess(OrderRepository)

	// Initialize the event consumer
	eventConsumer := usecase.NewEventConsumer(*orderProcessing)
	log.Printf("Event consumer initialized %v", eventConsumer)

	// Initialize the event hub adapter in a separate goroutine
	var cleanup context.CancelFunc
	go func() {
		var err error
		_, cleanup, err = eventhub.ConsumerInit(cfg.EventHubConnectionString, cfg.EventHubName, cfg.CheckpointStoreContainerName, cfg.CheckpointStoreConnectionString, eventConsumer)
		if err != nil {
			log.Println("Error: Failed to initialize event hub adapter")
			panic("Failed to initialize event hub adapter")
		}
		log.Println("Event hub adapter initialized")
	}()

	// Defer cleanup before entering the infinite loop
	defer func() {
		if cleanup != nil {
			log.Println("Cleaning up resources")
			cleanup()
		}
	}()

	// Create a channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Infinite loop
	for {
		select {
		case <-signals:
			log.Println("Received termination signal")
			return
		case <-time.After(1 * time.Minute):
			log.Println("Waiting for termination signal")
		}
	}
}
