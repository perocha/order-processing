package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/perocha/order-processing/pkg/appcontext"
	"github.com/perocha/order-processing/pkg/config"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/database/cosmosdb"
	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

func main() {
	// Initialize configuration
	cfg := config.InitializeConfig()
	if cfg == nil {
		// Print error
		log.Println("Error: Failed to load configuration")
		panic("Failed to load configuration")
	}

	// Initialize App Insights
	telemetryClient, err := telemetry.Initialize(cfg.AppInsightsInstrumentationKey)
	if err != nil {
		log.Printf("Failed to initialize App Insights %s\n", err.Error())
		panic("Failed to initialize App Insights")
	}
	// Add telemetry object to the context, so that it can be reused across the application
	ctx := context.WithValue(context.Background(), appcontext.TelemetryContextKey, telemetryClient)
	telemetryClient.TrackTrace(ctx, "Main::App Insights initialized", telemetry.Information, nil, true)

	// Initialize CosmosDB repository
	OrderRepository, err := cosmosdb.NewCosmosDBOrderRepository(ctx, cfg.CosmosDBConnectionString)
	if err != nil {
		telemetryClient.TrackException(ctx, "Failed to initialize CosmosDB repository", err, telemetry.Critical, nil, true)
		panic("Failed to initialize CosmosDB repository")
	}
	telemetryClient.TrackTrace(ctx, "Main::CosmosDB repository initialized", telemetry.Information, nil, true)

	// Start the orchestrator

	// Initialize the EventHub adapter
	//	eventHubAdapter, err := eventhub.EventHubAdapterInit(ctx, cfg.EventHubConnectionString, cfg.EventHubName, cfg.CheckpointStoreContainerName, cfg.CheckpointStoreConnectionString)

	// Create a channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Infinite loop
	for {
		select {
		case <-signals:
			telemetryClient.TrackTrace(ctx, "Main::Received termination signal", telemetry.Information, nil, true)
			return
		case <-time.After(1 * time.Minute):
			// Do nothing
			log.Println("Main::Waiting for termination signal")
		}
	}
}
