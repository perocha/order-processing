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

	"github.com/perocha/order-processing/pkg/domain"
	"github.com/perocha/order-processing/pkg/usecase"
)

type EventHubAdapter struct {
	eventProcessor   usecase.EventConsumer
	eventhubConsumer *azeventhubs.Processor
}

func ConsumerInit(eventHubConnectionString, eventHubName, containerName, checkpointStoreConnectionString string, eventProcessor usecase.EventConsumer) (*EventHubAdapter, context.CancelFunc, error) {
	// create a container client using a connection string and container name
	checkClient, err := container.NewClientFromConnectionString(checkpointStoreConnectionString, containerName, nil)
	if err != nil {
		log.Println("eventhub-subscriber::Error creating container client", err)
		panic(err)
	}

	// create a checkpoint store that will be used by the event hub
	checkpointStore, err := checkpoints.NewBlobStore(checkClient, nil)
	if err != nil {
		log.Println("eventhub-subscriber::Error creating checkpoint store", err)
		panic(err)
	}

	// create a consumer client using a connection string to the namespace and the event hub
	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(eventHubConnectionString, eventHubName, azeventhubs.DefaultConsumerGroup, nil)
	if err != nil {
		log.Println("eventhub-subscriber::Error creating consumer client", err)
		panic(err)
	}

	// Create a processor to receive and process events
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)
	if err != nil {
		log.Println("eventhub-subscriber::Error creating processor", err)
		panic(err)
	}

	adapter := &EventHubAdapter{
		eventProcessor:   eventProcessor,
		eventhubConsumer: processor,
	}

	// Run all partition clients
	go adapter.dispatchPartitionClients(processor)

	processorCtx, processorCancel := context.WithCancel(context.TODO())

	if err := processor.Run(processorCtx); err != nil {
		log.Println("eventhub-subscriber::Error processor run", err)
		processorCancel()
		consumerClient.Close(context.TODO())
		return nil, nil, err
	}

	cleanup := func() {
		processorCancel()
		consumerClient.Close(context.TODO())
	}

	return adapter, cleanup, nil
}

// For each partition in the event hub, create a partition client with processEvents as the function to process events
func (a *EventHubAdapter) dispatchPartitionClients(processor *azeventhubs.Processor) {
	for {
		// Track time and create a new operation ID, that will be used to track the end to end operation
		//startTime := time.Now()

		// Get the next partition client
		partitionClient := processor.NextPartitionClient(context.TODO())

		if partitionClient == nil {
			// No more partition clients to process
			break
		}

		go func() {
			// Define the operation ID using the defined OperationID type
			//operationID := uuid.New().String()
			//log.Printf("eventhub-subscriber::OperationID::%s::Creating new partition client\n", operationID)

			// Create a new context with the operation ID
			//ctx := context.WithValue(context.Background(), shared.OperationIDKeyContextKey, operationID)

			log.Printf("eventhub-subscriber::PartitionID::%s::Partition client initialized\n", partitionClient.PartitionID())

			// Process events for the partition client
			if err := a.processEvents(context.TODO(), partitionClient); err != nil {
				log.Printf("eventhub-subscriber::Error processing events for partition %s: %v\n", partitionClient.PartitionID(), err)
				panic(err)
			}
		}()
	}
}

// ProcessEvents implements the logic that is executed when events are received from the event hub
// func processEvents(ctx context.Context, partitionClient *azeventhubs.ProcessorPartitionClient) error {
func (a *EventHubAdapter) processEvents(ctx context.Context, partitionClient *azeventhubs.ProcessorPartitionClient) error {
	defer closePartitionResources(partitionClient)

	for {
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 10, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		// Uncomment the following line to verify that the consumer is trying to receive events
		log.Printf("eventhub-subscriber::PartitionID=%s::Processing %d event(s)\n", partitionClient.PartitionID(), len(events))

		for _, event := range events {
			log.Println("eventhub-subscriber::Message received: ", string(event.Body))

			// Events received!! Process the message
			msg := domain.Event{}
			// Unmarshal the event body into the message struct
			err := json.Unmarshal(event.Body, &msg)
			if err != nil {
				// Error unmarshalling the event body, discard the message
				log.Println("eventhub-subscriber::Error unmarshalling event body, discarding message", err)
			} else {
				// Process the message
				err = a.eventProcessor.ConsumeEvent(ctx, msg)
				if err != nil {
					log.Println("eventhub-subscriber::Error processing message, discarding message", err)
					//return err
				}
			}

			log.Printf("eventhub-subscriber::PartitionID::%s::Events received %v\n", partitionClient.PartitionID(), string(event.Body))
			log.Printf("Offset: %d Sequence number: %d MessageID: %s\n", event.Offset, event.SequenceNumber, *event.MessageID)
		}

		if len(events) != 0 {
			if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1], nil); err != nil {
				log.Println("eventhub-subscriber::Error updating checkpoint", err)
				return err
			}
		}
	}
}

// Closes the partition client
func closePartitionResources(partitionClient *azeventhubs.ProcessorPartitionClient) {
	defer partitionClient.Close(context.TODO())
}
