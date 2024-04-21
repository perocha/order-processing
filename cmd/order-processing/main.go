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

	// Initialize App Insights
	telemetryClient, err := telemetry.Initialize(cfg.AppInsightsInstrumentationKey, SERVICE_NAME)
	if err != nil {
		log.Fatalf("Main::Fatal error::Failed to initialize App Insights %s\n", err.Error())
	}
	// Add telemetry object to the context, so that it can be reused across the application
	ctx := context.WithValue(context.Background(), telemetry.TelemetryContextKey, telemetryClient)

	// Initialize CosmosDB repository
	orderRepository, err := cosmosdb.NewCosmosDBOrderRepository(ctx, cfg.CosmosdbEndpoint, cfg.CosmosdbConnectionString, cfg.CosmosdbDatabaseName, cfg.CosmosdbContainerName)
	if err != nil {
		telemetryClient.TrackException(ctx, "Main::Fatal error::Failed to initialize CosmosDB repository", err, telemetry.Critical, nil, true)
		log.Fatalf("Main::Fatal error::Failed to initialize CosmosDB repository %s\n", err.Error())
	}

	// Initialize EventHub subscriber adapter
	//	eventHubAdapter, err := eventhub.EventHubAdapterInit(ctx, cfg.EventHubName, cfg.EventHubConsumerConnectionString, cfg.EventHubProducerConnectionString, cfg.CheckpointStoreContainerName, cfg.CheckpointStoreConnectionString)

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

	telemetryClient.TrackTrace(ctx, "Main::All adapters initialized successfully", telemetry.Information, nil, true)

	// Create a channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Initialize service
	serviceInstance := service.Initialize(ctx, consumerInstance, producerInstance, orderRepository)
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
