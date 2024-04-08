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
	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/domain/order"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/repository/eventhub"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/repository/storage"
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
	OrderRepository, err := storage.NewCosmosDBOrderRepository(ctx, cfg.CosmosDBConnectionString)
	if err != nil {
		telemetryClient.TrackException(ctx, "Failed to initialize CosmosDB repository", err, telemetry.Critical, nil, true)
		panic("Failed to initialize CosmosDB repository")
	}
	telemetryClient.TrackTrace(ctx, "Main::CosmosDB repository initialized", telemetry.Information, nil, true)

	// Initialize order processing module
	orderProcessing := order.NewOrderProcess(OrderRepository)

	// Initialize the event consumer module
	eventConsumer := event.NewEventConsumer(*orderProcessing)
	telemetryClient.TrackTrace(ctx, "Main::Event consumer initialized", telemetry.Information, nil, true)

	// Initialize the event hub adapter in a separate goroutine
	var cleanup context.CancelFunc
	go func() {
		var err error
		_, cleanup, err = eventhub.ConsumerInit(ctx, cfg.EventHubConnectionString, cfg.EventHubName, cfg.CheckpointStoreContainerName, cfg.CheckpointStoreConnectionString, eventConsumer)
		if err != nil {
			telemetryClient.TrackException(ctx, "Failed to initialize event hub adapter", err, telemetry.Critical, nil, true)
			panic("Failed to initialize event hub adapter")
		}
		telemetryClient.TrackTrace(ctx, "Main::Event hub adapter initialized", telemetry.Information, nil, true)
	}()
	// Defer cleanup before entering the infinite loop
	defer func() {
		if cleanup != nil {
			telemetryClient.TrackTrace(ctx, "Main::Cleaning up resources", telemetry.Information, nil, true)
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
			telemetryClient.TrackTrace(ctx, "Main::Received termination signal", telemetry.Information, nil, true)
			return
		case <-time.After(1 * time.Minute):
			// Do nothing
			log.Println("Main::Waiting for termination signal")
		}
	}
}
