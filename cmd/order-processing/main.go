package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/perocha/goutils/pkg/telemetry"
	"github.com/perocha/order-processing/pkg/config"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/database/cosmosdb"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/messaging/eventhub"
	"github.com/perocha/order-processing/pkg/service"
)

const (
	SERVICE_NAME = "OrderProcessing"
	configFile   = "config.yaml"
)

func main() {
	// Initialize configuration
	cfg, err := config.InitializeConfig(configFile)
	if err != nil {
		// Print error
		log.Printf("Main::Fatal error::Failed to load configuration %s\n", err.Error())
		panic("Main::Failed to load configuration")
	}

	// Initialize App Insights
	telemetryClient, err := telemetry.Initialize(cfg.AppInsightsInstrumentationKey, SERVICE_NAME)
	if err != nil {
		log.Printf("Main::Fatal error::Failed to initialize App Insights %s\n", err.Error())
		panic("Main::Failed to initialize App Insights")
	}
	// Add telemetry object to the context, so that it can be reused across the application
	ctx := context.WithValue(context.Background(), telemetry.TelemetryContextKey, telemetryClient)

	// Initialize CosmosDB repository
	orderRepository, err := cosmosdb.NewCosmosDBOrderRepository(ctx, cfg.CosmosdbEndpoint, cfg.CosmosdbConnectionString, cfg.CosmosdbDatabaseName, cfg.CosmosdbContainerName)
	if err != nil {
		telemetryClient.TrackException(ctx, "Main::Fatal error::Failed to initialize CosmosDB repository", err, telemetry.Critical, nil, true)
		panic("Main::Fatal error::Failed to initialize CosmosDB repository")
	}

	// Initialize EventHub
	eventHubInstance, err := eventhub.EventHubAdapterInit(ctx, cfg.EventHubConnectionString, cfg.EventHubName, cfg.CheckpointStoreContainerName, cfg.CheckpointStoreConnectionString)
	if err != nil {
		telemetryClient.TrackException(ctx, "Main::Fatal error::Failed to initialize EventHub", err, telemetry.Critical, nil, true)
		panic("Main::Fatal error::Failed to initialize EventHub")
	}

	telemetryClient.TrackTrace(ctx, "Main::All adapters initialized successfully", telemetry.Information, nil, true)

	// Create a channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Initialize service
	serviceInstance := service.Initialize(ctx, eventHubInstance, orderRepository)
	go serviceInstance.Start(ctx, signals)

	telemetryClient.TrackTrace(ctx, "Main::Service layer initialized successfully", telemetry.Information, nil, true)

	// Infinite loop
	for {
		select {
		case <-signals:
			telemetryClient.TrackTrace(ctx, "Main::Received termination signal", telemetry.Information, nil, true)
			serviceInstance.Stop(ctx)
			return
		case <-time.After(2 * time.Minute):
			// Do nothing
			log.Println("Main::Waiting for termination signal")
		}
	}
}
