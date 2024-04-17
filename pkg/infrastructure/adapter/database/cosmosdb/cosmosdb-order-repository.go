package cosmosdb

import (
	"context"
	"encoding/json"
	"time"

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
	startTime := time.Now()
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Convert order to json
	orderJson, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Create partition key
	pk := azcosmos.NewPartitionKeyString(order.ProductCategory)

	// Create an item
	_, err = r.container.CreateItem(ctx, pk, orderJson, nil)
	if err != nil {
		return err
	}

	telemetryClient.TrackDependency(ctx, "CosmosDBOrderRepository", "CreateOrder", "CosmosDB", r.client.Endpoint(), true, startTime, time.Now(), nil, true)

	return nil
}

// Updates an existing order in CosmosDB
func (r *CosmosDBOrderRepository) UpdateOrder(ctx context.Context, order order.Order) error {
	startTime := time.Now()
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create partition key
	pk := azcosmos.NewPartitionKeyString(order.ProductCategory)

	// Convert order to json
	orderJson, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Update an item
	_, err = r.container.UpsertItem(ctx, pk, orderJson, nil)
	if err != nil {
		return err
	}

	telemetryClient.TrackDependency(ctx, "CosmosDBOrderRepository", "UpdateOrder", "CosmosDB", r.client.Endpoint(), true, startTime, time.Now(), nil, true)

	return nil
}

// Deletes an order from CosmosDB
func (r *CosmosDBOrderRepository) DeleteOrder(ctx context.Context, id, partitionKey string) error {
	startTime := time.Now()
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create partition key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Delete an item
	_, err := r.container.DeleteItem(ctx, pk, id, nil)
	if err != nil {
		return err
	}

	telemetryClient.TrackDependency(ctx, "CosmosDBOrderRepository", "DeleteOrder", "CosmosDB", r.client.Endpoint(), true, startTime, time.Now(), nil, true)

	return nil
}

// Retrieves an order from CosmosDB
func (r *CosmosDBOrderRepository) GetOrder(ctx context.Context, id, partitionKey string) (*order.Order, error) {
	startTime := time.Now()
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Create partition key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Retrieve an item
	item, err := r.container.ReadItem(ctx, pk, id, nil)
	if err != nil {
		return nil, err
	}

	// Convert item to order
	var order order.Order
	err = json.Unmarshal(item.Value, &order)
	if err != nil {
		return nil, err
	}

	telemetryClient.TrackDependency(ctx, "CosmosDBOrderRepository", "DeleteOrder", "CosmosDB", r.client.Endpoint(), true, startTime, time.Now(), nil, true)

	return &order, nil
}
