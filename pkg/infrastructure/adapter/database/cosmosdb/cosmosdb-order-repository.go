package cosmosdb

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/perocha/order-processing/pkg/domain/order"
	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

type CosmosDBOrderRepository struct {
	client    *azcosmos.Client
	database  *azcosmos.DatabaseClient
	container *azcosmos.ContainerClient
}

// Initialize CosmosDB repository using the provided connection string
func NewCosmosDBOrderRepository(ctx context.Context, endPoint, connectionString, databaseName, containerName string) (*CosmosDBOrderRepository, error) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create a new default azure credential
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		telemetryClient.TrackException(ctx, "CosmosDBOrderRepository::NewCosmosDBOrderRepository::Error creating default azure credential", err, telemetry.Critical, nil, true)
		return nil, err
	}

	// Create a new CosmosDB client
	clientOptions := azcosmos.ClientOptions{
		EnableContentResponseOnWrite: true,
	}
	client, err := azcosmos.NewClient(endPoint, credential, &clientOptions)
	if err != nil {
		telemetryClient.TrackException(ctx, "CosmosDBOrderRepository::NewCosmosDBOrderRepository::Error creating client", err, telemetry.Critical, nil, true)
		return nil, err
	}

	// Retrieve database
	database, err := client.NewDatabase(databaseName)
	if err != nil {
		return nil, err
	}

	// Create a new container
	container, err := database.NewContainer(containerName)
	if err != nil {
		return nil, err
	}

	return &CosmosDBOrderRepository{
		client:    client,
		database:  database,
		container: container,
	}, nil
}

// Creates a new order in CosmosDB
func (r *CosmosDBOrderRepository) CreateOrder(ctx context.Context, order order.Order) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	properties := order.ToMap()
	telemetryClient.TrackTrace(ctx, "CosmosDBOrderRepository::CreateOrder", telemetry.Information, properties, true)

	// New partition key
	pk := azcosmos.NewPartitionKeyString(order.ProductCategory)

	if order.Id == "" {
		// Generate error code
		err := errors.New("Id is required")
		return err
	}

	// Convert order to json
	orderJson, err := json.Marshal(order)
	if err != nil {
		return err
	}

	properties = map[string]string{
		"orderJson": string(orderJson),
	}
	telemetryClient.TrackTrace(ctx, "CosmosDBOrderRepository::CreateOrder::Order JSON", telemetry.Information, properties, true)

	// Create an item
	_, err = r.container.UpsertItem(ctx, pk, orderJson, nil)
	if err != nil {
		return err
	}

	return nil
}

// Updates an existing order in CosmosDB
func (r *CosmosDBOrderRepository) UpdateOrder(ctx context.Context, order order.Order) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	properties := order.ToMap()
	telemetryClient.TrackTrace(ctx, "CosmosDBOrderRepository::UpdateOrder", telemetry.Information, properties, true)
	return nil
}

// Deletes an order from CosmosDB
func (r *CosmosDBOrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	properties := map[string]string{
		"Id": orderID,
	}
	telemetryClient.TrackTrace(ctx, "CosmosDBOrderRepository::DeleteOrder", telemetry.Information, properties, true)
	return nil
}
