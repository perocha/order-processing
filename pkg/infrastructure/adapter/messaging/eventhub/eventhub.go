package eventhub

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/perocha/goutils/pkg/telemetry"
)

type EventHubAdapterImpl struct {
	ehProcessor      *azeventhubs.Processor
	ehConsumerClient *azeventhubs.ConsumerClient
	checkpointStore  *checkpoints.BlobStore
	checkClient      *container.Client
	ehProducerClient *azeventhubs.ProducerClient
	eventHubName     string
}

// Initializes only the consumer client
func ConsumerInitializer(ctx context.Context, eventHubName, consumerConnectionString, containerName, checkpointStoreConnectionString string) (*EventHubAdapterImpl, error) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// create a container client using a connection string and container name
	checkClient, err := container.NewClientFromConnectionString(checkpointStoreConnectionString, containerName, nil)
	if err != nil {
		telemetryClient.TrackException(ctx, "EventHubAdapter::Error creating container client", err, telemetry.Critical, nil, true)
		return nil, err
	}

	// create a checkpoint store that will be used by the event hub
	checkpointStore, err := checkpoints.NewBlobStore(checkClient, nil)
	if err != nil {
		telemetryClient.TrackException(ctx, "EventHubAdapter::Error creating checkpoint store", err, telemetry.Critical, nil, true)
		return nil, err
	}

	// create a consumer client using a connection string to the namespace and the event hub
	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(consumerConnectionString, eventHubName, azeventhubs.DefaultConsumerGroup, nil)
	if err != nil {
		telemetryClient.TrackException(ctx, "EventHubAdapter::Error creating consumer client", err, telemetry.Critical, nil, true)
		return nil, err
	}

	// Create a processor to receive and process events
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)
	if err != nil {
		telemetryClient.TrackException(ctx, "EventHubAdapter::Error creating processor", err, telemetry.Critical, nil, true)
		return nil, err
	}

	// Obtain the eventHubName from the consumer client
	eventHubProperties, err := consumerClient.GetEventHubProperties(ctx, nil)
	if err != nil {
		telemetryClient.TrackException(ctx, "EventHubAdapter::Error getting event hub properties", err, telemetry.Critical, nil, true)
		return nil, err
	}

	adapter := &EventHubAdapterImpl{
		ehProcessor:      processor,
		ehConsumerClient: consumerClient,
		checkpointStore:  checkpointStore,
		checkClient:      checkClient,
		eventHubName:     eventHubProperties.Name,
	}

	return adapter, nil
}

// Initializes only the producer client
func ProducerInitializer(ctx context.Context, eventHubName, producerConnectionString string) (*EventHubAdapterImpl, error) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create a new producer client
	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(producerConnectionString, eventHubName, nil)
	if err != nil {
		properties := map[string]string{
			"Error": err.Error(),
		}
		telemetryClient.TrackException(ctx, "EventHubAdapter::Failed to initialize producer", err, telemetry.Critical, properties, true)
		return nil, err
	}

	// Obtain the eventHubName from the producer client
	eventHubProperties, err := producerClient.GetEventHubProperties(ctx, nil)
	if err != nil {
		telemetryClient.TrackException(ctx, "EventHubAdapter::Error getting event hub properties", err, telemetry.Critical, nil, true)
		return nil, err
	}

	adapter := &EventHubAdapterImpl{
		ehProducerClient: producerClient,
		eventHubName:     eventHubProperties.Name,
	}

	return adapter, nil
}

// Close the EventHub adapter, both the consumer and producer clients
func (a *EventHubAdapterImpl) Close(ctx context.Context) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	telemetryClient.TrackTrace(ctx, "EventHubAdapter::Close::Stopping event hub consumer and producer clients", telemetry.Information, nil, true)

	// Close the consumer client
	if a.ehConsumerClient != nil {
		err := a.ehConsumerClient.Close(context.TODO())
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "EventHubAdapter::Consumer client close failed", err, telemetry.Critical, properties, true)
			return err
		}
	}

	// Close the producer client
	if a.ehProducerClient != nil {
		err := a.ehProducerClient.Close(ctx)
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "EventHubAdapter::Producer client close failed", err, telemetry.Critical, properties, true)
			return err
		}
	}

	telemetryClient.TrackTrace(ctx, "EventHubAdapter::Closing eventhub", telemetry.Information, nil, true)

	return nil
}
