package main

import (
	"context"
	"errors"
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

func main() {
	// Load the configuration
	cfg, err := config.InitializeConfig()
	if err != nil {
		log.Fatalf("Main::Fatal error::Failed to load configuration %s\n", err.Error())
	}
	// Refresh the configuration with the latest values
	err = cfg.RefreshConfig()
	if err != nil {
		log.Fatalf("Main::Fatal error::Failed to refresh configuration %s\n", err.Error())
	}

	telemetryConfig := telemetry.NewXTelemetryConfig(cfg.AppInsightsInstrumentationKey, SERVICE_NAME, "debug", 1)

	xTelemetry, err := telemetry.NewXTelemetry(telemetryConfig)
	if err != nil {
		log.Fatalf("Main::Fatal error::Failed to initialize XTelemetry %s\n", err.Error())
	}

	// Initialize App Insights
	telemetryClient, err := telemetry.Initialize(cfg.AppInsightsInstrumentationKey, SERVICE_NAME, "debug")
	if err != nil {
		log.Fatalf("Main::Fatal error::Failed to initialize App Insights %s\n", err.Error())
	}
	// Add telemetry object to the context, so that it can be reused across the application
	ctx := context.WithValue(context.Background(), telemetry.TelemetryContextKey, telemetryClient)

	// Log the configuration
	xTelemetry.Debug(ctx, "Main::Configuration loaded successfully DEBUG", telemetry.String("AppInsightsInstrumentationKey", cfg.AppInsightsInstrumentationKey))
	xTelemetry.Info(ctx, "Main::Configuration loaded successfully INFO",
		telemetry.String("AppInsightsInstrumentationKey", cfg.AppInsightsInstrumentationKey),
		telemetry.String("CosmosdbEndpoint", cfg.CosmosdbEndpoint),
		telemetry.String("CosmosdbConnectionString", cfg.CosmosdbConnectionString),
	)
	err = errors.New("this is a test error")
	xTelemetry.Error(ctx, "Main::Fatal error::Failed to initialize CosmosDB repository", telemetry.String("error", err.Error()))

	// Initialize CosmosDB repository
	orderRepository, err := cosmosdb.NewCosmosdbRepository(ctx, cfg.CosmosdbEndpoint, cfg.CosmosdbConnectionString, cfg.CosmosdbDatabaseName, cfg.CosmosdbContainerName)
	if err != nil {
		xTelemetry.Error(ctx, "Main::Fatal error::Failed to initialize CosmosDB repository", telemetry.String("error", err.Error()))
		telemetryClient.TrackException(ctx, "Main::Fatal error::Failed to initialize CosmosDB repository", err, telemetry.Critical, nil, true)
		log.Fatalf("Main::Fatal error::Failed to initialize CosmosDB repository %s\n", err.Error())
	}

	// Initialize EventHub consumer adapter
	consumerInstance, err := eventhub.ConsumerInitializer(ctx, cfg.EventHubNameConsumer, cfg.EventHubConsumerConnectionString, cfg.CheckpointStoreContainerName, cfg.CheckpointStoreConnectionString)
	if err != nil {
		telemetryClient.TrackException(ctx, "Main::Fatal error::Failed to initialize EventHub", err, telemetry.Critical, nil, true)
		log.Fatalf("Main::Fatal error::Failed to initialize EventHub %s\n", err.Error())
	}

	// Initialize EventHub publisher adapter
	producerInstance, err := eventhub.ProducerInitializer(ctx, cfg.EventHubNameProducer, cfg.EventHubProducerConnectionString)
	if err != nil {
		telemetryClient.TrackException(ctx, "Main::Fatal error::Failed to initialize EventHub", err, telemetry.Critical, nil, true)
		log.Fatalf("Main::Fatal error::Failed to initialize EventHub %s\n", err.Error())
	}

	xTelemetry.Info(ctx, "Main::All adapters initialized successfully", telemetry.String("ServiceName", SERVICE_NAME))

	// Create a channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Initialize service
	serviceInstance := service.Initialize(ctx, consumerInstance, producerInstance, orderRepository)
	go serviceInstance.Start(ctx, signals)

	xTelemetry.Info(ctx, "Main::Service layer initialized successfully", telemetry.String("ServiceName", SERVICE_NAME))

	// Infinite loop
	for {
		select {
		case <-signals:
			telemetryClient.TrackTrace(ctx, "Main::Received termination signal", telemetry.Information, nil, true)
			serviceInstance.Stop(ctx)
			return
		case <-time.After(2 * time.Minute):
			// Do nothing
			telemetryClient.TrackTrace(ctx, "Main::Waiting for termination signal", telemetry.Verbose, nil, true)
		}
	}
}
