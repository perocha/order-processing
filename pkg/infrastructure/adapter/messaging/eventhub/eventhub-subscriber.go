package eventhub

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"

	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

type EventHubAdapterImpl struct {
	ehProcessor      *azeventhubs.Processor
	ehConsumerClient *azeventhubs.ConsumerClient
	checkpointStore  *checkpoints.BlobStore
	checkClient      *container.Client
}

// Initializes a new EventHubAdapter
func EventHubAdapterInit(ctx context.Context, eventHubConnectionString, eventHubName, containerName, checkpointStoreConnectionString string) (*EventHubAdapterImpl, error) {
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
	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(eventHubConnectionString, eventHubName, azeventhubs.DefaultConsumerGroup, nil)
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

	adapter := &EventHubAdapterImpl{
		ehProcessor:      processor,
		ehConsumerClient: consumerClient,
		checkpointStore:  checkpointStore,
		checkClient:      checkClient,
	}

	return adapter, nil
}

func (a *EventHubAdapterImpl) Subscribe(ctx context.Context) (<-chan event.Event, context.CancelFunc, error) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	eventChannel := make(chan event.Event)

	// Run all partition clients
	go a.dispatchPartitionClients(ctx, eventChannel)

	processorCtx, processorCancel := context.WithCancel(context.TODO())

	go func() {
		if err := a.ehProcessor.Run(processorCtx); err != nil {
			telemetryClient.TrackException(ctx, "EventHubAdapter::Subscribe::Error processor run", err, telemetry.Critical, nil, true)
			processorCancel()
			a.ehConsumerClient.Close(context.TODO())
		}
	}()

	return eventChannel, processorCancel, nil
}

func (a *EventHubAdapterImpl) dispatchPartitionClients(ctx context.Context, eventChannel chan event.Event) {
	for {
		telemetryClient := telemetry.GetTelemetryClient(ctx)

		// Get the next partition client
		partitionClient := a.ehProcessor.NextPartitionClient(context.TODO())

		if partitionClient == nil {
			// No more partition clients to process
			break
		}

		go func() {
			telemetryClient.TrackTrace(ctx, "EventHubAdapter::dispatchPartitionClients::Partition ID "+partitionClient.PartitionID()+"::Client initialized", telemetry.Information, nil, true)

			// Process events for the partition client
			if err := a.processEventsForPartition(ctx, partitionClient, eventChannel); err != nil {
				properties := map[string]string{
					"PartitionID": partitionClient.PartitionID(),
				}
				telemetryClient.TrackException(ctx, "EventHubAdapter::dispatchPartitionClients::Error processing events", err, telemetry.Error, properties, true)
				panic(err)
			}
		}()
	}
}

// ProcessEvents implements the logic that is executed when events are received from the event hub
func (a *EventHubAdapterImpl) processEventsForPartition(ctx context.Context, partitionClient *azeventhubs.ProcessorPartitionClient, eventChannel chan event.Event) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	for {
		// Receive events from the partition client with a timeout of 20 seconds
		timeout := time.Second * 20
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), timeout)

		// Limit the wait for a number of events to receive
		limitEvents := 10
		events, err := partitionClient.ReceiveEvents(receiveCtx, limitEvents, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		// Uncomment the following line to verify that the consumer is trying to receive events
		log.Printf("EventHubAdapter::processEventsForPartition::PartitionID=%s::Processing %d event(s)\n", partitionClient.PartitionID(), len(events))

		for _, eventItem := range events {
			log.Println("EventHubAdapter::processEventsForPartition::Message received: ", string(eventItem.Body))

			// Events received!! Process the message
			msg := event.Event{}
			// Unmarshal the event body into the message struct
			err := json.Unmarshal(eventItem.Body, &msg)
			if err != nil {
				// Error unmarshalling the event body, discard the message
				telemetryClient.TrackTrace(ctx, "EventHubAdapter::processEventsForPartition::Error unmarshalling event body", telemetry.Error, nil, true)
				//return err
			} else {
				telemetryClient.TrackTrace(ctx, "EventHubAdapter::processEventsForPartition::PROCESS MESSAGE", telemetry.Information, nil, true)
				// Send the message to the event channel
				eventChannel <- msg
			}

			log.Printf("EventHubAdapter::processEventsForPartition::PartitionID::%s::Events received %v\n", partitionClient.PartitionID(), string(eventItem.Body))
			log.Printf("EventHubAdapter::processEventsForPartition::Offset: %d Sequence number: %d MessageID: %s\n", eventItem.Offset, eventItem.SequenceNumber, *eventItem.MessageID)
		}

		if len(events) != 0 {
			if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1], nil); err != nil {
				telemetryClient.TrackException(ctx, "EventHubAdapter::processEventsForPartition::Error updating checkpoint", err, telemetry.Error, nil, true)
				return err
			}
		}
	}
}

func (a *EventHubAdapterImpl) Init(ctx context.Context) error {
	return nil
}

func (a *EventHubAdapterImpl) Publish(ctx context.Context, message event.Event) error {
	return nil
}

func (a *EventHubAdapterImpl) Unsubscribe(ctx context.Context, topic string) error {
	return nil
}

func (a *EventHubAdapterImpl) Close(ctx context.Context) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	telemetryClient.TrackTrace(ctx, "EventHubAdapter::Close::Stopping event hub listener", telemetry.Information, nil, true)

	return nil
}
