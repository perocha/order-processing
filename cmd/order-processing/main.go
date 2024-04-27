package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/perocha/goadapters/database/cosmosdb"
	"github.com/perocha/goadapters/messaging/eventhub"
	"github.com/perocha/goutils/pkg/telemetry"
	"github.com/perocha/order-processing/pkg/config"
	"github.com/perocha/order-processing/pkg/service"
)

const (
	SERVICE_NAME = "OrderProcessing"
	configFile   = "config.yaml"
)

// Initialize configuration
func initialize() (*config.MicroserviceConfig, error) {
	// Load the configuration
	cfg, err := config.InitializeConfig()
	if err != nil {
		return nil, err
	}
	// Refresh the configuration with the latest values
	err = cfg.RefreshConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	// Initialize config
	cfg, err := initialize()
	if err != nil {
		log.Fatalf("Main::Fatal error::Failed to load configuration %s\n", err.Error())
	}

	// Initialize telemetry package
	telemetryConfig := telemetry.NewXTelemetryConfig(cfg.AppInsightsInstrumentationKey, SERVICE_NAME, "debug", 1)
	xTelemetry, err := telemetry.NewXTelemetry(telemetryConfig)
	if err != nil {
		log.Fatalf("Main::Fatal error::Failed to initialize XTelemetry %s\n", err.Error())
	}
	// Add telemetry object to the context, so that it can be reused across the application
	ctx := context.WithValue(context.Background(), telemetry.TelemetryContextKey, xTelemetry)

	// Initialize CosmosDB repository
	orderRepository, err := cosmosdb.NewCosmosdbRepository(ctx, cfg.CosmosdbEndpoint, cfg.CosmosdbConnectionString, cfg.CosmosdbDatabaseName, cfg.CosmosdbContainerName)
	if err != nil {
		xTelemetry.Error(ctx, "Main::Fatal error::Failed to initialize CosmosDB repository", telemetry.String("error", err.Error()))
		panic(err)
	}

	// Initialize EventHub consumer adapter
	consumerInstance, err := eventhub.ConsumerInitializer(ctx, cfg.EventHubNameConsumer, cfg.EventHubConsumerConnectionString, cfg.CheckpointStoreContainerName, cfg.CheckpointStoreConnectionString)
	if err != nil {
		xTelemetry.Error(ctx, "Main::Fatal error::Failed to initialize EventHub", telemetry.String("error", err.Error()))
		panic(err)
	}

	// Initialize EventHub publisher adapter
	producerInstance, err := eventhub.ProducerInitializer(ctx, cfg.EventHubNameProducer, cfg.EventHubProducerConnectionString)
	if err != nil {
		xTelemetry.Error(ctx, "Main::Fatal error::Failed to initialize EventHub", telemetry.String("error", err.Error()))
		panic(err)
	}

	xTelemetry.Info(ctx, "Main::All adapters initialized successfully")

	// Create a channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Initialize service
	serviceInstance := service.Initialize(ctx, consumerInstance, producerInstance, orderRepository)
	go serviceInstance.Start(ctx, signals)

	xTelemetry.Info(ctx, "Main::Service layer initialized successfully")

	// Infinite loop
	for {
		select {
		case <-signals:
			// Termination signal received
			xTelemetry.Info(ctx, "Main::Received termination signal")
			serviceInstance.Stop(ctx)
			return
		case <-time.After(2 * time.Minute):
			// Do nothing
			xTelemetry.Debug(ctx, "Main::Waiting for termination signal")
		}
	}
}
