package cosmosdb

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/perocha/order-processing/pkg/domain/order"
	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

type CosmosDBOrderRepository struct {
	client *azcosmos.Client
}

// Initialize CosmosDB repository using the provided connection string
func NewCosmosDBOrderRepository(ctx context.Context, endPoint string, connectionString string) (*CosmosDBOrderRepository, error) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	credential, err := azcosmos.NewKeyCredential(connectionString)
	if err != nil {
		properties := map[string]string{
			"Error": err.Error(),
		}
		telemetryClient.TrackException(ctx, "CosmosDBOrderRepository::NewCosmosDBOrderRepository::Error creating key credential", err, telemetry.Critical, properties, true)
		return nil, err
	}

	client, err := azcosmos.NewClientWithKey(endPoint, credential, nil)
	if err != nil {
		properties := map[string]string{
			"Error": err.Error(),
		}
		telemetryClient.TrackException(ctx, "CosmosDBOrderRepository::NewCosmosDBOrderRepository::Error creating client", err, telemetry.Critical, properties, true)
		return nil, err
	}

	telemetryClient.TrackTrace(ctx, "CosmosDBOrderRepository::NewCosmosDBOrderRepository::DB client created successfully", telemetry.Information, nil, true)

	return &CosmosDBOrderRepository{
		client: client,
	}, nil
}

// Creates a new order in CosmosDB
func (r *CosmosDBOrderRepository) CreateOrder(ctx context.Context, order order.Order) error {
	// Create a new container
	container, err := r.client.NewContainer("orders", "/id")
	if err != nil {
		return err
	}

	pk := azcosmos.NewPartitionKeyString("1")

	// Convert order to json
	orderJson, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Create an item
	_, err = container.CreateItem(ctx, pk, orderJson, nil)
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
		"OrderID": orderID,
	}
	telemetryClient.TrackTrace(ctx, "CosmosDBOrderRepository::DeleteOrder", telemetry.Information, properties, true)
	return nil
}
