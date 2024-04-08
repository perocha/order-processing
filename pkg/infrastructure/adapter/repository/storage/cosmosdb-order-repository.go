package storage

import (
	"context"
	"log"

	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/perocha/order-processing/pkg/domain/order"
)

type CosmosDBOrderRepository struct {
	client *azcosmos.Client
}

// Initialize CosmosDB repository using the provided connection string
func NewCosmosDBOrderRepository(ctx context.Context, connectionString string) (*CosmosDBOrderRepository, error) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)

	log.Printf("CosmosDBOrderRepository::NewCosmosDBOrderRepository::Client=%v::Error=%v", client, err)
	telemetryClient.TrackTrace(ctx, "CosmosDBOrderRepository::NewCosmosDBOrderRepository", telemetry.Information, nil, true)

	return nil, nil
	/*
		if err != nil {
			return nil, err
		}

		return &CosmosDBOrderRepository{
			client: client,
		}, nil
	*/
}

// Creates a new order in CosmosDB
func (r *CosmosDBOrderRepository) CreateOrder(ctx context.Context, order order.Order) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	properties := order.ToMap()
	telemetryClient.TrackTrace(ctx, "CosmosDBOrderRepository::CreateOrder", telemetry.Information, properties, true)
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
